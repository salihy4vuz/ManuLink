package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"

	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/config"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/database"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/models"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/middlewares"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/rabbitmq"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/routes"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/services"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/ws"
)

func main() {
	// Initialize environment variables if .env exists
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on system environment variables.")
	}

	// Validate environment variables before proceeding
	config.ValidateEnv()

	// Connect to PostgreSQL
	config.ConnectDB()
	database.ConnectDB()

	// Connect to MinIO
	config.ConnectMinio()
	services.EnsureBucketsExist()

	// Init OpenTelemetry
	shutdownOtel := config.InitOpenTelemetry()
	defer shutdownOtel(context.Background())

	// Auto-Migrate schema
	database.AutoMigrate(
		&models.Tenant{},
		&models.TenantTheme{},
		&models.User{},
		&models.Project{},
		&models.Part{},
		&models.Quote{},
		&models.Order{},
		&models.OrderItem{},
		&models.MESOrder{},
		&models.LMEPriceLog{},
		&models.AuditLog{},
		&models.Payment{},
		&models.ApiKey{},
		&models.WebhookEndpoint{},
		&models.PricingConfig{}, // Faz 20 (E): Dinamik fiyatlandırma parametreleri
	)

	// Seed Database with MVP data
	database.SeedDatabase()

	// Connect to RabbitMQ
	mqUser := os.Getenv("RABBITMQ_USER")
	if mqUser == "" { mqUser = "guest" }
	mqPass := os.Getenv("RABBITMQ_PASSWORD")
	if mqPass == "" { mqPass = "guest" }
	mqHost := os.Getenv("RABBITMQ_HOST")
	if mqHost == "" { mqHost = "localhost" }
	mqPort := os.Getenv("RABBITMQ_PORT")
	if mqPort == "" { mqPort = "5672" }
	mqUrl := fmt.Sprintf("amqp://%s:%s@%s:%s/", mqUser, mqPass, mqHost, mqPort)
	rabbitmq.ConnectRabbitMQ(mqUrl)
	defer rabbitmq.CloseRabbitMQ()

	// Start Background Cron Jobs (Free Currency API Fetcher & LME Prices)
	services.StartCurrencyCronJob()
	services.StartLMECronJob()

	// Initialize Fiber app
	app := fiber.New()

	// Middlewares
	app.Use(recover.New())

	// Faz 1.3: Observability (Request ID & Structured Logging)
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format:     `{"time":"${time}","level":"info","req_id":"${locals:requestid}","method":"${method}","path":"${path}","status":${status},"latency":"${latency}","ip":"${ip}"}` + "\n",
		TimeFormat: time.RFC3339,
		TimeZone:   "UTC",
	}))

	// Faz 1.7: Circuit Breaker
	app.Use(middlewares.CircuitBreakerMiddleware())

	// Faz 1.8: i18n Language parsing
	app.Use(middlewares.I18nMiddleware())

	// Faz 1.4: Security (Helmet & Global Rate Limiting for DDOS)
	app.Use(helmet.New())
	app.Use(limiter.New(limiter.Config{
		Max:        100, // Max 100 requests per IP
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			msg := middlewares.T(c, "error_too_many_requests")
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": msg,
			})
		},
	}))

	// Faz 1.5: Redis / Memory Caching for fast responses (e.g., Currency API)
	app.Use(cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			// Only cache GET requests (like LME rates, parts list)
			return c.Method() != "GET"
		},
		Expiration:   5 * time.Minute,
		CacheControl: true,
	}))

	app.Use(middlewares.CORSMiddleware()) // CORS middleware

	// Swagger B2B Documentation Route
	app.Static("/swagger.json", "./docs/swagger.json")
	app.Get("/swagger", func(c *fiber.Ctx) error {
		return c.SendString(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<title>Swagger UI - B2B API</title>
			<link rel="stylesheet" type="text/css" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.15.5/swagger-ui.css" >
		</head>
		<body>
			<div id="swagger-ui"></div>
			<script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.15.5/swagger-ui-bundle.js"> </script>
			<script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.15.5/swagger-ui-standalone-preset.js"> </script>
			<script>
			window.onload = function() {
				const ui = SwaggerUIBundle({
				url: "/swagger.json",
				dom_id: '#swagger-ui',
				deepLinking: true,
				presets: [
					SwaggerUIBundle.presets.apis,
					SwaggerUIStandalonePreset
				],
				layout: "StandaloneLayout"
				});
				window.ui = ui;
			}
			</script>
		</body>
		</html>`)
	})

	// Initialize WebSockets
	routes.InitWebSocket(app)
	// Initialize WebSocket Hub
	hub := ws.NewHub()
	go hub.Run()

	// Start Phase 6.8 Tenant WebSocket Hub
	services.DefaultHub = services.NewHub()
	go services.DefaultHub.Run()

	// Initialize API Routes
	routes.SetupAPIRoutes(app)

	app.Get("/health", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		// PostgreSQL ping
		dbStatus := "Online"
		dbLatency := "N/A"
		if sqlDB, err := config.DB.DB(); err == nil {
			start := time.Now()
			if err := sqlDB.PingContext(ctx); err != nil {
				dbStatus = "Degraded: " + err.Error()
			} else {
				dbLatency = time.Since(start).String()
			}
		} else {
			dbStatus = "Error: " + err.Error()
		}

		// Redis ping
		redisStatus := "Online"
		redisLatency := "N/A"
		if rc := services.GetRedisClient(); rc != nil {
			start := time.Now()
			if err := rc.Ping(ctx).Err(); err != nil {
				redisStatus = "Degraded: " + err.Error()
			} else {
				redisLatency = time.Since(start).String()
			}
		} else {
			redisStatus = "Not configured"
		}

		// MinIO ping
		minioStatus := "Online"
		if err := config.PingMinio(ctx); err != nil {
			minioStatus = "Degraded: " + err.Error()
		}

		// RabbitMQ connection status
		rabbitStatus := rabbitmq.ConnectionStatus()

		// Rust Engine health: independent check via RabbitMQ queue inspection
		// The geometry-engine is considered healthy if RabbitMQ is online AND
		// the geometry_tasks queue exists and is reachable (consumer count signal).
		rustEngineStatus := "Unknown"
		if rabbitStatus == "Online" {
			ch := rabbitmq.GetChannel()
			if ch != nil {
				queueInfo, inspectErr := ch.QueueInspect("geometry_tasks")
				if inspectErr != nil {
					rustEngineStatus = "Degraded: queue not reachable"
				} else if queueInfo.Consumers == 0 {
					rustEngineStatus = "Degraded: no active consumers (engine may be down)"
				} else {
					rustEngineStatus = "Online"
				}
			} else {
				rustEngineStatus = "Degraded: channel unavailable"
			}
		} else {
			rustEngineStatus = "Offline: RabbitMQ disconnected"
		}

		overallStatus := "Healthy"
		if dbStatus != "Online" || rabbitStatus != "Online" {
			overallStatus = "Degraded"
		}
		if rustEngineStatus != "Online" {
			overallStatus = "Degraded"
		}

		httpStatus := fiber.StatusOK
		if overallStatus == "Degraded" {
			httpStatus = fiber.StatusServiceUnavailable
		}

		return c.Status(httpStatus).JSON(fiber.Map{
			"status":    overallStatus,
			"version":   "1.9.0",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"services": fiber.Map{
				"database":    fiber.Map{"status": dbStatus, "latency": dbLatency},
				"redis":       fiber.Map{"status": redisStatus, "latency": redisLatency},
				"minio":       fiber.Map{"status": minioStatus},
				"rabbitmq":    fiber.Map{"status": rabbitStatus},
				"rust_engine": fiber.Map{"status": rustEngineStatus},
			},
		})
	})

	// WebSocket Endpoint for CAD Plugins (SolidWorks, Fusion 360)
	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		client := ws.NewClient(hub, c)
		client.Hub.Register(client)

		// Allow collection of memory referenced by the caller by doing all work in
		// new goroutines.
		go client.WritePump()
		client.ReadPump()
	}))

	// Graceful Shutdown
	go func() {
		if err := app.Listen(":3000"); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Gracefully shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

// Trigger rebuild

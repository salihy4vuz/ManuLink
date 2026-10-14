package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func SetupSecurityMiddlewares(app *fiber.App) {
	// 1. Recover from Panics
	app.Use(recover.New())

	// 2. Structured Logging
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))

	// 3. CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Configure this to frontend URL in production
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// 4. Rate Limiting (DDoS Protection)
	app.Use(limiter.New(limiter.Config{
		Max:        100,             // Max 100 requests
		Expiration: 1 * time.Minute, // per minute per IP
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests. Please try again later.",
			})
		},
	}))

	// 5. GZIP Compression (Crucial for large GLTF 3D files)
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
}

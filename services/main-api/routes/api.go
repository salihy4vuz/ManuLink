package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/controllers"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/middlewares"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/services"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/telemetry"
)

func SetupAPIRoutes(app *fiber.App) {
	// Global Telemetry Middleware
	app.Use(telemetry.MetricsMiddleware())

	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Multi-Tenant Endpoints
	v1.Post("/auth/login", func(c *fiber.Ctx) error {
		// Mock login endpoint for frontend compatibility
		return c.JSON(fiber.Map{"token": "admin-token"})
	})
	v1.Get("/tenants/config", controllers.GetTenantConfig)
	// Requires standard auth for updating
	// v1.Put("/tenants/config", middlewares.Protected(), middlewares.TenantContext(), controllers.UpdateTenantConfig)
	v1.Put("/tenants/config", middlewares.TenantContext(), controllers.UpdateTenantConfig)

	// ERP Integration (Faz 6.3)
	v1.Post("/orders/:id/sync-erp", middlewares.TenantContext(), controllers.SyncOrderToERP)

	// Dynamic PDF Generation (Faz 6.4)
	v1.Get("/quotes/:id/pdf", middlewares.TenantContext(), controllers.GetQuotePDF)
	v1.Get("/quotes/:id/certificate", middlewares.TenantContext(), controllers.GetMaterialCertificate)

	// WebSocket Production Tracking (Faz 6.8)
	// First middleware checks upgrade header and sets TenantID
	api.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	
	// Register WS route with TenantContext to identify the client
	api.Get("/ws/tracking", middlewares.TenantContext(), websocket.New(controllers.WsTrackingHandler))
	
	// Internal HTTP endpoint to trigger a WS broadcast for testing
	v1.Post("/internal/tracking/update", middlewares.TenantContext(), controllers.InternalUpdateTracking)

	// B2B Quoting Endpoint with Anti-DDoS and JWT Auth
	quoteLimiter := middlewares.NewRateLimiter(middlewares.RateLimiterConfig{
		RedisClient: services.GetRedisClient(),
		MaxRequests: 10,
		Window:      1 * time.Minute,
	})
	api.Post("/quote", quoteLimiter, middlewares.Protected(), controllers.GenerateQuote)
	
	// Frontend SPA Routes (mapped under /v1)
	v1.Post("/upload", controllers.UploadCADFile)
	v1.Post("/quotes", quoteLimiter, middlewares.Protected(), controllers.GenerateQuote)
	v1.Post("/orders", middlewares.Protected(), controllers.PlaceOrder)
	v1.Get("/orders", middlewares.Protected(), controllers.GetOrders)
	v1.Get("/admin/stats", middlewares.Protected(), controllers.GetStats)

	api.Post("/machining/quote", quoteLimiter, middlewares.Protected(), controllers.GenerateMachiningQuote) // Yeni: Talaşlı İmalat
	api.Post("/quotes/:id/approve", middlewares.Protected(), controllers.ApproveQuote) // Sprint 4: CAD-to-MES Otonom Onay
	api.Post("/quote/test-rpc", controllers.TestRPCBridge) // Test RPC Bridge without JWT
	
	generalLimiter := middlewares.NewRateLimiter(middlewares.RateLimiterConfig{
		RedisClient: services.GetRedisClient(),
		MaxRequests: 100,
		Window:      1 * time.Minute,
	})
	api.Get("/lme-prices", generalLimiter, controllers.GetLMEPrices)
	
	// Phase 6.9 Stripe Payments
	api.Post("/payments/create-intent", middlewares.Protected(), controllers.CreatePaymentIntent)
	api.Post("/payments/webhook", controllers.StripeWebhook) // Webhook has its own signature verification, no JWT

	// Phase 6.9 Iyzico Payments
	api.Post("/payments/iyzico/create-intent", middlewares.Protected(), controllers.CreateIyzicoPaymentIntent)
	api.Post("/payments/iyzico/webhook", controllers.IyzicoWebhook)

	// Phase 3.4 AI LME Endpoints
	api.Get("/ai/lme/predict", generalLimiter, controllers.PredictLMEPrice)
	api.Post("/ai/lme/train", middlewares.Protected(), controllers.TriggerLMEModelTraining)
	
	// Faz 10: Generative AI (LLM) DFM Advisor
	api.Post("/ai/advisor/analyze", generalLimiter, controllers.AnalyzeCADWithAI)

	// Faz 11: 2D Nesting Optimizer Endpoint
	v1.Post("/nesting/calculate", generalLimiter, controllers.CalculateNesting)

	// Prometheus Metrics Endpoint
	app.Get("/metrics", telemetry.MetricsHandler)

	// User Subscription endpoints
	legacyFeature := middlewares.RequireFeature(middlewares.FeatureFlagConfig{RedisClient: services.GetRedisClient()}, "legacy_erp")
	api.Post("/subscription/upgrade", legacyFeature, controllers.UpgradeSubscription)

	// B2B Integration Endpoints (Requires standard JWT login to manage)
	b2b := api.Group("/b2b")
	b2b.Use(middlewares.Protected())
	b2b.Post("/api-keys/generate", controllers.GenerateAPIKey)
	b2b.Post("/webhooks/register", controllers.RegisterWebhook)

	// External B2B API (Protected by API Key instead of JWT)
	b2bExternal := api.Group("/external/b2b")
	b2bExternal.Use(middlewares.ApiKeyAuth())
	b2bExternal.Post("/quote", controllers.GenerateQuote)
	b2bExternal.Post("/order/place", controllers.PlaceOrder)
	b2bExternal.Post("/upload", controllers.UploadCADFile)
	b2bExternal.Get("/orders", controllers.GetOrders)
	b2bExternal.Get("/stats", controllers.GetStats)

	// CAD Upload & Export
	cadGroup := api.Group("/cad")
	cadGroup.Get("/formats", controllers.GetSupportedFormats)
	cadGroup.Post("/upload", controllers.UploadCADFile)
	cadGroup.Post("/analyze-dxf", controllers.AnalyzeDXFContent)
	cadGroup.Post("/lantek/export", controllers.LantekExport)
	cadGroup.Get("/:id/viewer", controllers.ViewerGLTF)
	cadGroup.Get("/:id/download-dxf", controllers.DownloadDXF)

	// Forge integration is an enterprise feature, lock behind legacy flag
	// cadGroup.Post("/forge/webhook", legacyFeature, controllers.HandleForgeWebhook)


	// User Support Ticket endpoints
	api.Get("/tickets", legacyFeature, controllers.GetUserTickets)
	api.Post("/tickets/new", legacyFeature, controllers.CreateTicket)

	// Admin Endpoints for Pricing Management
	admin := api.Group("/admin")
	admin.Use(middlewares.Protected())
	admin.Use(middlewares.RequireRole(middlewares.RoleAdmin))

	admin.Get("/materials", controllers.GetMaterials)
	admin.Post("/materials/update", controllers.UpdateMaterialPrice)

	admin.Get("/machines", controllers.GetMachines)
	admin.Post("/machines/update", controllers.UpdateMachineConfig)

	// Admin Stats for Dashboard
	admin.Get("/stats", controllers.GetStats)

	// Admin Order Management
	admin.Put("/orders/:id/status", controllers.UpdateOrderStatus)

	// Faz 20 (E): Dinamik Fiyatlandırma Konfigürasyonu
	admin.Get("/pricing-config", controllers.GetPricingConfigHandler)
	admin.Put("/pricing-config", controllers.UpdatePricingConfigHandler)
}

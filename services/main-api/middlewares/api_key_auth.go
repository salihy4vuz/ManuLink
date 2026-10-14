package middlewares

import (
	"strings"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/services"
	"github.com/gofiber/fiber/v2"
)

// ApiKeyAuth middleware validates B2B API keys via X-API-Key header
func ApiKeyAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get("X-API-Key")
		
		// Optional: Allow Bearer token in Authorization header for API Key
		if apiKey == "" {
			authHeader := c.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer sk_") {
				apiKey = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if apiKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing X-API-Key header",
			})
		}

		tenantID, err := services.ValidateAPIKey(apiKey)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or revoked API Key",
			})
		}

		// Store tenantID in locals for the next handlers
		c.Locals("tenant_id", tenantID)

		return c.Next()
	}
}

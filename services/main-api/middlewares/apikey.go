package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/database"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/models"
)

// APIKeyAuth validates B2B API keys
func APIKeyAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get("X-API-Key")
		if apiKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "X-API-Key header required",
			})
		}

		// Admin key bypass for internal systems
		if apiKey == "admin-key" {
			return c.Next()
		}

		var keyRecord models.ApiKey
		if err := database.DB.Where("key = ? AND is_active = ?", apiKey, true).First(&keyRecord).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or revoked API key",
			})
		}

		// Store tenant ID for downstream handlers
		c.Locals("tenant_id", keyRecord.TenantID)
		return c.Next()
	}
}

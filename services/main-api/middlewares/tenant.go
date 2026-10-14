package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/database"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/models"
)

// TenantContext extracts the tenant domain from the request and injects the Tenant ID
func TenantContext() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Try to get tenant from custom header (useful for API/mobile clients)
		domain := c.Get("X-Tenant-Domain")

		// 2. Fallback to Hostname if header is empty (useful for Web clients)
		if domain == "" {
			domain = c.Hostname()
			// In a real scenario, you might parse "acme.sacmetal.com" to "acme"
		}

		var tenant models.Tenant
		if err := database.DB.Where("domain = ?", domain).First(&tenant).Error; err != nil {
			// If no tenant is matched, we can either reject or assign a default fallback
			// For strict multi-tenancy, we return 404
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Tenant not found for domain: " + domain,
			})
		}

		// Inject TenantID into Fiber Locals for subsequent handlers
		c.Locals("tenant_id", tenant.ID)

		return c.Next()
	}
}

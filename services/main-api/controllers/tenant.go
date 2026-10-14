package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/database"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/models"
)

// GetTenantConfig returns the theme configuration for the requested domain
func GetTenantConfig(c *fiber.Ctx) error {
	domain := c.Query("hostname")
	if domain == "" {
		domain = c.Hostname()
	}

	var tenant models.Tenant
	if err := database.DB.Preload("Theme").Where("domain = ?", domain).First(&tenant).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Tenant not found",
		})
	}

	// Build the response mapped to the Vue.js store expectations
	response := fiber.Map{
		"id":             tenant.ID,
		"name":           tenant.Name,
		"primaryColor":   tenant.Theme.PrimaryColor,
		"secondaryColor": tenant.Theme.SecondaryColor,
		"logoUrl":        tenant.Theme.LogoURL,
		"faviconUrl":     tenant.Theme.FaviconURL,
	}

	return c.JSON(response)
}

// UpdateTenantConfig allows the admin to update their white-label settings
func UpdateTenantConfig(c *fiber.Ctx) error {
	// In real life, TenantID would come from auth token + locals
	tenantID := c.Locals("tenant_id")
	if tenantID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	type UpdateRequest struct {
		PrimaryColor   string `json:"primaryColor"`
		SecondaryColor string `json:"secondaryColor"`
		LogoURL        string `json:"logoUrl"`
	}

	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var theme models.TenantTheme
	if err := database.DB.Where("tenant_id = ?", tenantID).First(&theme).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Theme not found"})
	}

	theme.PrimaryColor = req.PrimaryColor
	theme.SecondaryColor = req.SecondaryColor
	theme.LogoURL = req.LogoURL

	database.DB.Save(&theme)

	return c.JSON(fiber.Map{"message": "Theme updated successfully"})
}

package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/models"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/services"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/database"
)

func GenerateAPIKey(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(uint)
	
	type request struct {
		Name string `json:"name"`
	}
	var req request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	key, err := services.GenerateAPIKey(tenantID, req.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate key"})
	}

	return c.JSON(fiber.Map{
		"message": "Key generated successfully. Store this key, it will not be shown again.",
		"api_key": key,
	})
}

func RegisterWebhook(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(uint)
	
	type request struct {
		URL    string `json:"url"`
		Secret string `json:"secret"`
		Events string `json:"events"`
	}
	var req request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	endpoint := models.WebhookEndpoint{
		TenantID: tenantID,
		URL:      req.URL,
		Secret:   req.Secret,
		Events:   req.Events,
		IsActive: true,
	}

	if err := database.DB.Create(&endpoint).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not register webhook"})
	}

	return c.JSON(fiber.Map{
		"message": "Webhook registered successfully",
		"id":      endpoint.ID,
	})
}

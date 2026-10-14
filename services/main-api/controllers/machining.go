package controllers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/config"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/models"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/services"
	"gorm.io/datatypes"
)

// GenerateMachiningQuote takes a STEP file and calculates CNC Milling cost
func GenerateMachiningQuote(c *fiber.Ctx) error {
	var req services.MachiningQuoteRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	resp, err := services.CalculateMachiningQuote(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Machining pricing engine failed",
			"details": err.Error(),
		})
	}

	tenantID, _ := c.Locals("TenantID").(uint)
	if tenantID == 0 {
		tenantID = 1
	}

	detailsJSON, _ := json.Marshal(resp)

	quoteRecord := models.Quote{
		TenantID:       tenantID,
		TotalCost:      resp.TotalCostUSD,
		FinalPrice:     resp.FinalSellingPrice,
		Currency:       resp.Currency,
		Status:         "draft",
		QuoteType:      "machining", // Key differentiator
		PricingDetails: datatypes.JSON(detailsJSON),
	}

	if err := config.DB.Create(&quoteRecord).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"quote_id": quoteRecord.ID,
		"quote": resp,
	})
}

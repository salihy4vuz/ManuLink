package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/services"
)

// GetLMEPrices returns the latest LME prices
func GetLMEPrices(c *fiber.Ctx) error {
	// Fetch real dynamic prices from the services package
	allMaterials := services.GetAllMaterials()
	
	prices := fiber.Map{}
	for key, data := range allMaterials {
		prices[key] = data.PricePerKg
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    prices,
	})
}


// PredictLMEPrice predicts future LME prices using AI
func PredictLMEPrice(c *fiber.Ctx) error {
	// Stub implementation
	return c.JSON(fiber.Map{
		"success": true,
		"prediction": fiber.Map{
			"steel":    860.0,
			"aluminum": 2150.0,
			"copper":   8600.0,
		},
	})
}

// TriggerLMEModelTraining triggers the AI model training for LME prices
func TriggerLMEModelTraining(c *fiber.Ctx) error {
	// Stub implementation
	return c.JSON(fiber.Map{
		"success": true,
		"message": "AI model training started successfully",
	})
}

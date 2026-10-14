package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/services"
)

// AnalyzeCADWithAI endpoint
func AnalyzeCADWithAI(c *fiber.Ctx) error {
	var req services.DFMRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Geçersiz AI parametreleri.",
		})
	}

	// Call the AI Advisor (DFM local implementation)
	aiResult := services.AnalyzeDFM(req)

	return c.JSON(fiber.Map{
		"success": true,
		"data":    aiResult,
	})
}

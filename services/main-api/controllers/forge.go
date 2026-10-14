//go:build ignore
// Rapor Kararı: Autodesk Forge, aylık ücretli kurumsal bir platform.
// KOBİ stratejisiyle çeliştiği için derleme dışı bırakıldı. Dosya silinmedi.

package controllers

import (
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/rabbitmq"
)

// ForgeWebhookPayload represents the expected JSON from Autodesk Forge Model Derivative API
type ForgeWebhookPayload struct {
	Version string `json:"version"`
	ResourceURN string `json:"resourceUrn"`
	Hook struct {
		System string `json:"system"`
		Event  string `json:"event"`
	} `json:"hook"`
	Payload struct {
		ProjectID    string  `json:"projectId"`
		Material     string  `json:"material"`
		ThicknessMm  float64 `json:"thickness_mm"`
		BoundingBox  struct {
			X float64 `json:"x"`
			Y float64 `json:"y"`
		} `json:"boundingBox"`
	} `json:"payload"`
}

// HandleForgeWebhook receives CAD data from Autodesk Forge (Cloud CAD)
func HandleForgeWebhook(c *fiber.Ctx) error {
	var body ForgeWebhookPayload

	if err := c.BodyParser(&body); err != nil {
		log.Printf("Forge Webhook Parsing Error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid Forge Payload",
		})
	}

	log.Printf("Received Autodesk Forge Webhook for URN: %s", body.ResourceURN)

	// In a real enterprise system, we route this to RabbitMQ for Pricing and CRM updates
	taskPayload := map[string]interface{}{
		"source":         "Autodesk_Forge",
		"material":       body.Payload.Material,
		"thickness_mm":   body.Payload.ThicknessMm,
		"bounding_box_x": body.Payload.BoundingBox.X,
		"bounding_box_y": body.Payload.BoundingBox.Y,
		"bend_count":     4, // Usually extracted from Forge metadata
	}

	// Trigger Rust Engine or Go Pricing Engine asynchronously
	_, err := rabbitmq.CallRustEngine(c.Context(), "forge_sync", taskPayload)
	if err != nil {
		log.Printf("Failed to sync Forge data to RabbitMQ: %v", err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Webhook processed successfully",
	})
}

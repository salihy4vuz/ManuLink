package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/rabbitmq"
)

type PartToNest struct {
	ID       string  `json:"id"`
	Width    float64 `json:"width"`
	Height   float64 `json:"height"`
	Quantity int     `json:"quantity"`
}

type NestingRequest struct {
	SheetWidth    float64      `json:"sheet_width"`
	SheetHeight   float64      `json:"sheet_height"`
	AllowRotation bool         `json:"allow_rotation"`
	Parts         []PartToNest `json:"parts"`
}

// CalculateNesting sends a nesting layout request to the Rust CAM Engine
func CalculateNesting(c *fiber.Ctx) error {
	var req NestingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Geçersiz Nesting parametreleri.",
		})
	}

	// Payload struct matching Rust expected schema
	taskPayload := map[string]interface{}{
		"sheet_width":    req.SheetWidth,
		"sheet_height":   req.SheetHeight,
		"allow_rotation": req.AllowRotation,
		"parts":          req.Parts,
	}

	// Make RPC call to Rust geometry engine (blocking for result)
	// We use "nesting_request" as the queue/action name
	result, err := rabbitmq.CallRustEngine(c.Context(), "nesting", taskPayload)
	if err != nil {
		log.Printf("Failed to calculate nesting via Rust engine: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Nesting motoru cevap vermedi.",
		})
	}

	// Result should already be the JSON struct output from NestingResponse in Rust
	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

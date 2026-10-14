package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/rabbitmq"
)

func TestRPCBridge(c *fiber.Ctx) error {
	var payload interface{}
	if err := c.BodyParser(&payload); err != nil {
		// If no body, use a default dummy payload
		payload = map[string]interface{}{
			"material":  "Steel",
			"thickness": 2.0,
			"radius":    2.0,
			"angle":     90.0,
		}
	}

	taskType := c.Query("task", "k_factor")

	// Call the Rust Engine
	res, err := rabbitmq.CallRustEngine(c.Context(), taskType, payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Type("json").SendString(res) // Pass through the raw JSON from Rust
}

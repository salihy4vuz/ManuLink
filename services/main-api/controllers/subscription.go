package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/services"
)

type UpgradeRequest struct {
	Tier string `json:"tier"` // "Pro" or "Enterprise"
}

// UpgradeSubscription handles the mocked Stripe upgrade
func UpgradeSubscription(c *fiber.Ctx) error {
	username := c.Get("X-User", "admin")
	
	var req UpgradeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	tier := services.SubscriptionTier(req.Tier)
	
	if err := services.UpgradeUserTier(username, tier); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Subscription upgraded to " + req.Tier,
		"tier":    tier,
	})
}

package controllers

import (
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/rabbitmq"
)

// StartOnboarding - POST /api/v1/onboarding/start
func StartOnboarding(c *fiber.Ctx) error { 
	// RabbitMQ'ya "tenant.provision" eventi at
	ch := rabbitmq.GetChannel()
	if ch != nil {
		payload, _ := json.Marshal(map[string]string{"event": "tenant_created", "company": "Test"})
		ch.PublishWithContext(context.Background(), "", "tenant.provision", false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
		})
	}
	return c.SendStatus(fiber.StatusOK) 
}

// SetupMachines - POST /api/v1/onboarding/machines
func SetupMachines(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) }

// SetupMaterials - POST /api/v1/onboarding/materials
func SetupMaterials(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) }

// SetupERP - POST /api/v1/onboarding/erp
func SetupERP(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) }

// CompleteOnboarding - POST /api/v1/onboarding/complete
func CompleteOnboarding(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) }

// GetOnboardingStatus - GET /api/v1/onboarding/status
func GetOnboardingStatus(c *fiber.Ctx) error { return c.JSON(fiber.Map{"progress": 100}) }

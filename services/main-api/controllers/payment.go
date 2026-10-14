package controllers

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"

	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/database"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/models"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/services"
)

// CheckoutRequest represents the body sent from frontend to initiate payment
type CheckoutRequest struct {
	OrderID uint `json:"order_id"`
}

// CreatePaymentIntent handles the generation of Stripe Client Secret
func CreatePaymentIntent(c *fiber.Ctx) error {
	tenantIDRaw := c.Locals("tenant_id")
	if tenantIDRaw == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	tenantID := tenantIDRaw.(uint)

	var req CheckoutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Fetch Order from DB to ensure amount cannot be manipulated by frontend
	var order models.Order
	if err := database.DB.Where("id = ? AND status = ?", req.OrderID, models.StatusPendingPayment).First(&order).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found or not in pending payment state"})
	}

	// Create Intent
	clientSecret, err := services.CreatePaymentIntent(order.TotalAmountUSD, order.ID, tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create payment intent", "details": err.Error()})
	}

	return c.JSON(fiber.Map{
		"client_secret": clientSecret,
		"order_id":      order.ID,
		"amount":        order.TotalAmountUSD,
	})
}

// StripeWebhook handles asynchronous payment confirmations from Stripe
func StripeWebhook(c *fiber.Ctx) error {
	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if endpointSecret == "" {
		endpointSecret = "whsec_mockSecretForDevOnly"
	}

	payload := c.Body()
	sigHeader := c.Get("Stripe-Signature")

	event, err := webhook.ConstructEvent(payload, sigHeader, endpointSecret)
	if err != nil {
		log.Printf("⚠️  Webhook signature verification failed. %v\n", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// Handle the event
	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			log.Printf("Error parsing webhook JSON: %v\n", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		log.Printf("💰 PaymentIntent for %d cents was successful!", paymentIntent.Amount)
		
		// In a real scenario, we would parse metadata to get the OrderID
		// and update the order status in the database to models.StatusInProduction
		// e.g.:
		// database.DB.Model(&models.Order{}).Where("id = ?", orderID).Update("status", models.StatusInProduction)
		
		// And we would also Trigger the ERP sync and WebSocket broadcast!

	default:
		log.Printf("Unhandled event type: %s\n", event.Type)
	}

	return c.SendStatus(fiber.StatusOK)
}

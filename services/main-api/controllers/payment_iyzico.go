package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/database"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/models"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/services"
)

// CreateIyzicoPaymentIntent handles the generation of Iyzico Checkout Form Token
func CreateIyzicoPaymentIntent(c *fiber.Ctx) error {
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

	// Create Intent via Iyzico
	token, err := services.CreateIyzicoPaymentIntent(order.TotalAmountUSD, order.ID, tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create Iyzico payment intent", "details": err.Error()})
	}

	return c.JSON(fiber.Map{
		"token":         token,
		"order_id":      order.ID,
		"amount":        order.TotalAmountUSD,
		"payment_url":   "https://sandbox-api.iyzipay.com/payment/iyzipay/checkoutform/initialize/auth/ecom", // Usually handled via Iyzico SDK
	})
}

// IyzicoWebhook handles asynchronous payment confirmations from Iyzico
func IyzicoWebhook(c *fiber.Ctx) error {
	// Payload will contain token and status from Iyzico redirect/callback
	// In a real integration, we query Iyzico's /payment/iyzipay/checkoutform/auth/ecom/detail endpoint using the token
	// to verify if the payment was actually successful.
	
	type IyzicoCallbackPayload struct {
		Token  string `json:"token" form:"token"`
		Status string `json:"status" form:"status"`
	}

	var payload IyzicoCallbackPayload
	if err := c.BodyParser(&payload); err != nil {
		log.Printf("Error parsing Iyzico webhook: %v\n", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	log.Printf("💰 Iyzico Webhook received! Token: %s, Status: %s", payload.Token, payload.Status)

	if payload.Status == "success" {
		// Verify token and update order
		// database.DB.Model(&models.Order{}).Where("payment_id = ?", payload.Token).Update("status", models.StatusInProduction)
	}

	// Iyzico usually requires a 200 OK or redirect to success page
	return c.SendStatus(fiber.StatusOK)
}

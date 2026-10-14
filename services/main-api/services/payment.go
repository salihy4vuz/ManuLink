package services

import (
	"errors"
	"fmt"
	"log"

	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/models"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/database"
	"os"
	
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
)

// ProcessPayment creates a real payment intent via Stripe
func ProcessPayment(tenantID uint, orderID uint, amount float64, currency string) (*models.Payment, error) {
	if amount <= 0 {
		return nil, errors.New("invalid payment amount")
	}

	apiKey := os.Getenv("STRIPE_SECRET_KEY")
	if apiKey == "" {
		return nil, errors.New("Stripe API key not configured")
	}
	stripe.Key = apiKey

	// Create a PaymentIntent
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(amount * 100)), // Stripe expects smallest currency unit (cents)
		Currency: stripe.String(currency),
		Metadata: map[string]string{
			"order_id":  fmt.Sprintf("%d", orderID),
			"tenant_id": fmt.Sprintf("%d", tenantID),
		},
	}
	
	pi, err := paymentintent.New(params)
	if err != nil {
		log.Printf("[PAYMENT] Failed to create payment intent: %v\n", err)
		return nil, err
	}

	payment := models.Payment{
		TenantID:      tenantID,
		OrderID:       orderID,
		TransactionID: pi.ID,
		Provider:      "stripe",
		Amount:        amount,
		Currency:      currency,
		Status:        "pending", // Wait for webhook to mark as completed
	}

	if err := database.DB.Create(&payment).Error; err != nil {
		return nil, err
	}

	log.Printf("[PAYMENT] Created Stripe PaymentIntent %s for Order %d\n", pi.ID, orderID)

	return &payment, nil
}

package services

import (
	"errors"
	"os"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
)

// CreatePaymentIntent calls Stripe API to create an intent to pay for an order
// It returns the client_secret which the frontend needs to complete the payment.
func CreatePaymentIntent(amountUSD float64, orderID uint, tenantID uint) (string, error) {
	apiKey := os.Getenv("STRIPE_SECRET_KEY")
	if apiKey == "" {
		// Fallback for development if not provided in .env
		apiKey = "sk_test_mockKeyForDevOnly"
	}
	stripe.Key = apiKey

	// Stripe accepts amounts in the smallest currency unit (e.g., cents)
	amountCents := int64(amountUSD * 100)

	if amountCents <= 0 {
		return "", errors.New("invalid payment amount")
	}

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amountCents),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		Metadata: map[string]string{
			"order_id":  string(rune(orderID)), // store order ID for webhook matching
			"tenant_id": string(rune(tenantID)),
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return "", err
	}

	return pi.ClientSecret, nil
}

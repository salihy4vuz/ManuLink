package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

// IyzicoCheckoutRequest represents the payload to create a checkout form
type IyzicoCheckoutRequest struct {
	Locale          string  `json:"locale"`
	ConversationID  string  `json:"conversationId"`
	Price           string  `json:"price"`
	PaidPrice       string  `json:"paidPrice"`
	Currency        string  `json:"currency"`
	BasketID        string  `json:"basketId"`
	PaymentGroup    string  `json:"paymentGroup"`
	CallbackURL     string  `json:"callbackUrl"`
	// In a real app, you'd add Buyer, ShippingAddress, BillingAddress, and BasketItems objects here
}

// IyzicoCheckoutResponse represents the response from Iyzico API
type IyzicoCheckoutResponse struct {
	Status          string `json:"status"`
	ErrorCode       string `json:"errorCode"`
	ErrorMessage    string `json:"errorMessage"`
	CheckoutFormContent string `json:"checkoutFormContent"`
	Token           string `json:"token"`
	PaymentPageURL  string `json:"paymentPageUrl"`
}

// CreateIyzicoPaymentIntent connects to Iyzico REST API to generate a Checkout Form Token
func CreateIyzicoPaymentIntent(amountUSD float64, orderID uint, tenantID uint) (string, error) {
	apiKey := os.Getenv("IYZICO_API_KEY")
	secretKey := os.Getenv("IYZICO_SECRET_KEY")
	baseURL := os.Getenv("IYZICO_BASE_URL")

	if baseURL == "" {
		baseURL = "https://sandbox-api.iyzipay.com" // Default to sandbox
	}
	if apiKey == "" || secretKey == "" {
		return "", errors.New("Iyzico API keys are not configured in environment variables")
	}

	if amountUSD <= 0 {
		return "", errors.New("invalid payment amount")
	}

	// Format price to string (Iyzico expects string like "15.0")
	priceStr := fmt.Sprintf("%.2f", amountUSD)

	payload := IyzicoCheckoutRequest{
		Locale:         "tr",
		ConversationID: fmt.Sprintf("%d-%d", tenantID, orderID),
		Price:          priceStr,
		PaidPrice:      priceStr,
		Currency:       "USD", // Or TRY based on order.Currency
		BasketID:       fmt.Sprintf("ORD-%d", orderID),
		PaymentGroup:   "B2B",
		CallbackURL:    "https://api.sacmetal.com/api/v1/payments/iyzico/webhook", // Example Callback
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", baseURL+"/payment/iyzipay/checkoutform/initialize/auth/ecom", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	
	// Authentication headers for Iyzico PKI String generation goes here in a real implementation
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var iyzicoResp IyzicoCheckoutResponse
	if err := json.NewDecoder(resp.Body).Decode(&iyzicoResp); err != nil {
		return "", err
	}

	if iyzicoResp.Status != "success" {
		return "", fmt.Errorf("iyzico error: %s - %s", iyzicoResp.ErrorCode, iyzicoResp.ErrorMessage)
	}

	// Returning the Token which the frontend needs to render the Iyzico iframe
	return iyzicoResp.Token, nil
}

package services

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/models"
	"log"
)

// ERPWebhookPayload represents the Industry 4.0 standard JSON pushed to ERP/MES
type ERPWebhookPayload struct {
	OrderID         uint        `json:"order_id"`
	QuoteID         uint        `json:"quote_id"`
	Status          string      `json:"status"`
	ProductionSteps interface{} `json:"routing"`
	Timestamp       time.Time   `json:"timestamp"`
}

// DispatchMESOrderToERP sends the MES Order data to the configured ERP Webhook URL with Exponential Backoff Resilience
func DispatchMESOrderToERP(order models.MESOrder, webhookURL string) {
	if webhookURL == "" {
		log.Println("[Webhook Dispatcher] No webhook URL configured. Skipping.")
		return
	}

	payload := ERPWebhookPayload{
		OrderID:         order.ID,
		QuoteID:         order.QuoteID,
		Status:          order.Status,
		ProductionSteps: order.ProductionSteps, // This contains the detailed routing
		Timestamp:       time.Now(),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[Webhook Dispatcher] Failed to marshal payload: %v\n", err)
		return
	}

	// Exponential Backoff implementation
	maxRetries := 3
	backoff := 2 * time.Second

	go func() {
		for attempt := 1; attempt <= maxRetries; attempt++ {
			log.Printf("[Webhook Dispatcher] Pushing MESOrder %d to %s (Attempt %d/%d)...\n", order.ID, webhookURL, attempt, maxRetries)

			resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
			
			if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
				log.Printf("[Webhook Dispatcher] Successfully pushed to %s\n", webhookURL)
				resp.Body.Close()
				return
			}

			if resp != nil {
				resp.Body.Close()
				log.Printf("[Webhook Dispatcher] Failed with status %d\n", resp.StatusCode)
			} else {
				log.Printf("[Webhook Dispatcher] Request failed: %v\n", err)
			}

			if attempt < maxRetries {
				log.Printf("[Webhook Dispatcher] Retrying in %v...\n", backoff)
				time.Sleep(backoff)
				backoff *= 2 // Exponentially increase wait time
			}
		}

		log.Printf("[Webhook Dispatcher] CRITICAL: Failed to push MESOrder %d to ERP after %d attempts. Dead Letter logic should handle this.\n", order.ID, maxRetries)
	}()
}

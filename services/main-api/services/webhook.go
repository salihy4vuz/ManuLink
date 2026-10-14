package services

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/database"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/models"
)

type WebhookPayload struct {
	Event     string      `json:"event"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// DispatchWebhookToTenants sends an event to all active endpoints for a tenant
func DispatchWebhookToTenants(tenantID uint, event string, data interface{}) {
	var endpoints []models.WebhookEndpoint
	database.DB.Where("tenant_id = ? AND is_active = ?", tenantID, true).Find(&endpoints)

	for _, ep := range endpoints {
		if strings.Contains(ep.Events, event) || ep.Events == "*" {
			go sendWebhook(ep, event, data)
		}
	}
}

func sendWebhook(ep models.WebhookEndpoint, event string, data interface{}) {
	payload := WebhookPayload{
		Event:     event,
		Timestamp: time.Now().Unix(),
		Data:      data,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Errorf("Failed to marshal webhook payload: %v", err)
		return
	}

	// Sign payload
	mac := hmac.New(sha256.New, []byte(ep.Secret))
	mac.Write(body)
	signature := hex.EncodeToString(mac.Sum(nil))

	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		req, _ := http.NewRequest("POST", ep.URL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Signature", signature)

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)

		if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			log.Infof("Successfully sent webhook %s to %s", event, ep.URL)
			resp.Body.Close()
			return
		}

		if resp != nil {
			log.Warnf("Webhook failed with status %d: %s", resp.StatusCode, ep.URL)
			resp.Body.Close()
		} else {
			log.Warnf("Webhook failed: %v", err)
		}

		time.Sleep(time.Duration(i*2) * time.Second)
	}

	log.Errorf("Webhook delivery failed after %d retries for %s", maxRetries, ep.URL)
}

// DispatchWebhookEvent legacy adapter if still needed
func DispatchWebhookEvent(event string, data interface{}, targetURL string) {
	go func() {
		payload := WebhookPayload{
			Event:     event,
			Timestamp: time.Now().Unix(),
			Data:      data,
		}
		body, _ := json.Marshal(payload)
		http.Post(targetURL, "application/json", bytes.NewBuffer(body))
	}()
}

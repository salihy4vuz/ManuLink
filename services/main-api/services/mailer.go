package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"
)

// SendQuotationEmail sends the generated PDF via SMTP
func SendQuotationEmail(customerEmail, customerName string, pdfBytes []byte) error {
	// In production, these should be loaded from .env
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT") // e.g., "587"
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	// If SMTP is not configured, we just log and skip
	if smtpHost == "" || smtpPort == "" {
		return errors.New("SMTP credentials missing, cannot send email")
	}

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	// Build the email with attachment (MIME standard)
	boundary := "Apex-Sheet-Metal-Boundary-123456789"
	
	var body bytes.Buffer
	body.WriteString(fmt.Sprintf("To: %s\r\n", customerEmail))
	body.WriteString("Subject: Your Sheet Metal Quotation from Apex Inc.\r\n")
	body.WriteString("MIME-Version: 1.0\r\n")
	body.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n\r\n", boundary))

	// Email text body
	body.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	body.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n\r\n")
	body.WriteString(fmt.Sprintf("Dear %s,\n\nPlease find attached the official Proforma Invoice for your sheet metal manufacturing request.\n\nDue to LME metal volatility, prices are valid for 24 hours.\n\nRegards,\nApex Autonomous Factory\n\n", customerName))

	// Email PDF Attachment
	body.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	body.WriteString("Content-Type: application/pdf; name=\"Quotation.pdf\"\r\n")
	body.WriteString("Content-Disposition: attachment; filename=\"Quotation.pdf\"\r\n")
	body.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")
	
	// Encode bytes to Base64
	encoded := base64.StdEncoding.EncodeToString(pdfBytes)
	
	// Break base64 string into lines of 76 characters per MIME standard
	for i := 0; i < len(encoded); i += 76 {
		end := i + 76
		if end > len(encoded) {
			end = len(encoded)
		}
		body.WriteString(encoded[i:end] + "\r\n")
	}
	body.WriteString("\r\n")

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{customerEmail}, body.Bytes())
	if err != nil {
		log.Printf("Failed to send email to %s: %v\n", customerEmail, err)
		return err
	}

	log.Printf("Successfully sent Proforma PDF to %s\n", customerEmail)
	return nil
}

// DispatchWebhook sends the quotation data to an external ERP (SAP, Logo, etc.)
func DispatchWebhook(customerName string, quote PricingResponse) {
	log.Printf("[WEBHOOK DISPATCHER] Sending Quote Data for %s to ERP Webhook...\n", customerName)
	
	webhookURL := os.Getenv("ERP_WEBHOOK_URL")
	if webhookURL == "" {
		log.Println("[WEBHOOK DISPATCHER] ERP_WEBHOOK_URL is not configured. Webhook dispatch aborted.")
		return
	}

	payload := map[string]interface{}{
		"customer_name": customerName,
		"quote_details": quote,
		"timestamp":     time.Now().Format(time.RFC3339),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[WEBHOOK DISPATCHER] Failed to marshal payload: %v\n", err)
		return
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("[WEBHOOK DISPATCHER] HTTP POST Failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		log.Printf("[WEBHOOK DISPATCHER] ERP rejected webhook with status: %d\n", resp.StatusCode)
		return
	}
	
	log.Printf("[WEBHOOK DISPATCHER] Success! Quote total: %.2f %s logged to ERP.\n", quote.FinalSellingPrice, quote.Currency)
}

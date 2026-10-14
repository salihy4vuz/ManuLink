package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// SendEmailNotification sends an email to a user via SendGrid API
func SendEmailNotification(to string, subject string, body string) error {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	if apiKey == "" {
		return errors.New("SENDGRID_API_KEY is not configured")
	}

	payload := map[string]interface{}{
		"personalizations": []map[string]interface{}{
			{
				"to": []map[string]string{{"email": to}},
			},
		},
		"from": map[string]string{"email": "noreply@sacmetal.com"},
		"subject": subject,
		"content": []map[string]string{
			{
				"type": "text/plain",
				"value": body,
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://api.sendgrid.com/v3/mail/send", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("SendGrid returned status: %d", resp.StatusCode)
	}

	log.Printf("[Email Sent] To: %s | Subject: %s\n", to, subject)
	return nil
}

// SendSMSNotification sends an SMS to a user via Twilio
func SendSMSNotification(phone string, message string) error {
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	twilioNumber := os.Getenv("TWILIO_PHONE_NUMBER")

	if accountSid == "" || authToken == "" || twilioNumber == "" {
		return errors.New("Twilio credentials are not configured")
	}

	urlStr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSid)
	
	msgData := url.Values{}
	msgData.Set("To", phone)
	msgData.Set("From", twilioNumber)
	msgData.Set("Body", message)
	
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(msgData.Encode()))
	if err != nil {
		return err
	}
	
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Twilio returned status: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	log.Printf("[SMS Sent] To: %s", phone)
	return nil
}

// NotifyOrderCreated orchestrates notifications for a new order
func NotifyOrderCreated(email, phone string, orderID uint) {
	msg := fmt.Sprintf("Your order #%d has been successfully created and is waiting for production.", orderID)
	
	if email != "" {
		go SendEmailNotification(email, "Order Created", msg)
	}
	
	if phone != "" {
		go SendSMSNotification(phone, msg)
	}
}

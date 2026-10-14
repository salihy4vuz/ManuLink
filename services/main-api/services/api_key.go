package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/models"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/database"
)

// GenerateAPIKey generates a new plain API key and saves its hash to DB
func GenerateAPIKey(tenantID uint, name string) (string, error) {
	// Generate random 32 bytes
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	
	plainKey := "sk_" + hex.EncodeToString(bytes)
	
	// Hash the key for storage
	hash := sha256.Sum256([]byte(plainKey))
	hashedKey := hex.EncodeToString(hash[:])

	apiKey := models.ApiKey{
		TenantID: tenantID,
		KeyHash:  hashedKey,
		Name:     name,
		Revoked:  false,
	}

	if err := database.DB.Create(&apiKey).Error; err != nil {
		return "", err
	}

	// We only return the plain key once!
	return plainKey, nil
}

// ValidateAPIKey checks if the provided plain key is valid and returns its TenantID
func ValidateAPIKey(plainKey string) (uint, error) {
	if plainKey == "" {
		return 0, errors.New("empty API key")
	}

	hash := sha256.Sum256([]byte(plainKey))
	hashedKey := hex.EncodeToString(hash[:])

	var apiKey models.ApiKey
	if err := database.DB.Where("key_hash = ? AND revoked = ?", hashedKey, false).First(&apiKey).Error; err != nil {
		return 0, errors.New("invalid or revoked API key")
	}

	return apiKey.TenantID, nil
}

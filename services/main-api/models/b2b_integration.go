package models

import (
	"time"

	"gorm.io/gorm"
)

// ApiKey model for B2B API integrations
type ApiKey struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TenantID  uint           `gorm:"not null;index" json:"tenant_id"`
	KeyHash   string         `gorm:"not null" json:"-"`
	Name      string         `gorm:"not null" json:"name"` // E.g., "ERP Integration Key"
	Revoked   bool           `gorm:"default:false" json:"revoked"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// WebhookEndpoint model for sending events to customers
type WebhookEndpoint struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	TenantID    uint           `gorm:"not null;index" json:"tenant_id"`
	URL         string         `gorm:"not null" json:"url"`
	Secret      string         `gorm:"not null" json:"secret"` // Used to sign the webhook payload
	Events      string         `gorm:"not null" json:"events"` // Comma separated, e.g., "order.created,production.completed"
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

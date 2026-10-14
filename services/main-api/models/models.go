package models

import (
	"time"

	"gorm.io/gorm"
)



// User model with Soft Delete and Multi-Tenant Support
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TenantID  uint           `gorm:"not null;index" json:"tenant_id"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	Name      string         `json:"name"`
	Role      string         `gorm:"default:'user'" json:"role"` // admin, operator, customer
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Project model
type Project struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	TenantID    uint           `gorm:"not null;index" json:"tenant_id"`
	UserID      uint           `gorm:"not null" json:"user_id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description"`
	Status      string         `gorm:"default:'draft'" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}


// LMEPriceLog model for LME Fiyat Logları
type LMEPriceLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Material  string    `gorm:"not null" json:"material"`
	Price     float64   `gorm:"not null" json:"price"`
	Currency  string    `gorm:"default:'USD'" json:"currency"`
	Date      time.Time `gorm:"not null" json:"date"`
	CreatedAt time.Time `json:"created_at"`
}

// AuditLog model for tracking changes (Anti-Espionage & Security)
type AuditLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TenantID   uint      `gorm:"not null;index" json:"tenant_id"`
	UserID     uint      `gorm:"not null;index" json:"user_id"`
	Action     string    `gorm:"not null" json:"action"`     // e.g., "DOWNLOAD_CAD", "UPDATE_PRICE"
	EntityType string    `gorm:"not null" json:"entity_type"` // e.g., "Order", "Project"
	EntityID   uint      `json:"entity_id"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	CreatedAt  time.Time `gorm:"index" json:"created_at"`
}

// Payment model for handling Stripe/Iyzico payments
type Payment struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	TenantID      uint           `gorm:"not null;index" json:"tenant_id"`
	OrderID       uint           `gorm:"not null" json:"order_id"`
	TransactionID string         `gorm:"not null" json:"transaction_id"`
	Provider      string         `gorm:"not null" json:"provider"` // stripe, iyzico, bank_transfer
	Amount        float64        `gorm:"not null" json:"amount"`
	Currency      string         `gorm:"default:'USD'" json:"currency"`
	Status        string         `gorm:"default:'pending'" json:"status"` // pending, completed, failed
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

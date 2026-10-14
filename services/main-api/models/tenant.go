package models

import (
	"time"

	"gorm.io/gorm"
)

// Tenant represents a B2B customer (White-Label Client)
type Tenant struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	Domain    string         `gorm:"uniqueIndex;not null" json:"domain"` // e.g., acme.sacmetal.com or acme-factory
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Theme TenantTheme `gorm:"foreignKey:TenantID" json:"theme"`

	// Multi-Region ERP Configuration
	ERPSystem   string `gorm:"default:'none'" json:"erp_system"` // e.g., "sap", "odoo", "logo", "netsuite"
	ERPEndpoint string    `json:"erp_endpoint"` // Endpoint of the external system
	ERPToken    string    `json:"erp_token"`    // Access token for the external system
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	FeatureFlag string    `json:"feature_flag"` // Special flags
	Users       []User    `json:"users"`        // A tenant can have many users
}

// TenantTheme holds the white-label branding configurations
type TenantTheme struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	TenantID       uint   `gorm:"uniqueIndex;not null" json:"tenant_id"`
	PrimaryColor   string `gorm:"default:'#3498db'" json:"primaryColor"`
	SecondaryColor string `gorm:"default:'#2c3e50'" json:"secondaryColor"`
	LogoURL        string `gorm:"default:'/default-logo.svg'" json:"logoUrl"`
	FaviconURL     string `gorm:"default:'/favicon.ico'" json:"faviconUrl"`
}

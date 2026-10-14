package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Quote struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	TenantID       uint           `gorm:"not null;index" json:"tenant_id"`
	CustomerID     uint           `gorm:"index" json:"customer_id"`
	TotalCost      float64        `json:"total_cost"`
	FinalPrice     float64        `json:"final_price"`
	Currency       string         `gorm:"size:3;default:'USD'" json:"currency"`
	Status         string         `gorm:"size:20;default:'draft'" json:"status"` // draft, pending_approval, approved, rejected
	QuoteType      string         `gorm:"size:20;default:'sheet_metal'" json:"quote_type"` // sheet_metal, machining
	PricingDetails datatypes.JSON `gorm:"type:jsonb" json:"pricing_details"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate is a GORM hook to log the creation of a Quote
func (q *Quote) BeforeCreate(tx *gorm.DB) (err error) {
	auditLog := AuditLog{
		TenantID:   q.TenantID,
		UserID:     q.CustomerID, // Assuming Customer is the user creating it
		Action:     "CREATE_QUOTE",
		EntityType: "Quote",
		EntityID:   0, // Not available yet
	}
	return tx.Create(&auditLog).Error
}

// BeforeUpdate is a GORM hook to log updates to a Quote (e.g. status changes)
func (q *Quote) BeforeUpdate(tx *gorm.DB) (err error) {
	auditLog := AuditLog{
		TenantID:   q.TenantID,
		UserID:     q.CustomerID,
		Action:     "UPDATE_QUOTE",
		EntityType: "Quote",
		EntityID:   q.ID,
	}
	return tx.Create(&auditLog).Error
}

package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type MESOrder struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	TenantID        uint           `gorm:"not null;index" json:"tenant_id"`
	QuoteID         uint           `gorm:"index" json:"quote_id"`
	Status          string         `gorm:"default:'queued'" json:"status"` // queued, cutting, bending, done
	ProductionSteps datatypes.JSON `gorm:"type:jsonb;index:,type:gin" json:"production_steps"` // Machines to be used, estimated times, operator instructions
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

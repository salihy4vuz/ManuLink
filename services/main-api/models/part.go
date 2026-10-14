package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Part struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	TenantID     uint           `gorm:"not null;index" json:"tenant_id"`
	FileHash     string         `gorm:"uniqueIndex;not null" json:"file_hash"`
	FileName      string         `json:"file_name"`
	FileURL       string         `json:"file_url"`        // MinIO link to raw CAD file
	ViewerGLTFUrl string         `json:"viewer_gltf_url"` // Endpoint returning 3D GLTF
	DXFUrl        string         `json:"dxf_url"`         // Endpoint returning 2D DXF
	Material      string         `json:"material"`
	ThicknessMm   float64        `json:"thickness_mm"`
	GeometryData  datatypes.JSON `gorm:"type:jsonb;index:,type:gin" json:"geometry_data"` // Holds outer_perimeter, bend_count, volumes, etc.
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate is a GORM hook to log CAD file uploads
func (p *Part) BeforeCreate(tx *gorm.DB) (err error) {
	auditLog := AuditLog{
		TenantID:   p.TenantID,
		UserID:     0, // Would be fetched from context in a real flow
		Action:     "UPLOAD_CAD_PART",
		EntityType: "Part",
		EntityID:   0,
	}
	return tx.Create(&auditLog).Error
}

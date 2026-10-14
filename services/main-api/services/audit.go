package services

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/database"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/models"
)

// LogAudit helper function to record user actions securely
func LogAudit(c *fiber.Ctx, action, entityType string, entityID uint) {
	tenantIDRaw := c.Locals("tenant_id")
	userIDRaw := c.Locals("user_id") // Assuming auth middleware sets this

	if tenantIDRaw == nil || userIDRaw == nil {
		log.Println("[AUDIT_ERROR] Missing tenant_id or user_id in context")
		return
	}

	tenantID := tenantIDRaw.(uint)
	userID := userIDRaw.(uint)

	auditLog := models.AuditLog{
		TenantID:   tenantID,
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		IPAddress:  c.IP(),
		UserAgent:  c.Get("User-Agent"),
	}

	if err := database.DB.Create(&auditLog).Error; err != nil {
		log.Printf("[AUDIT_ERROR] Failed to save audit log: %v\n", err)
	}
}

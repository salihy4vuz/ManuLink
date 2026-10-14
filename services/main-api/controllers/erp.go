package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/database"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/models"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/services"
)

// SyncOrderToERP manually triggers pushing an order to the Tenant's ERP
func SyncOrderToERP(c *fiber.Ctx) error {
	tenantIDRaw := c.Locals("tenant_id")
	if tenantIDRaw == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	tenantID := tenantIDRaw.(uint)

	orderIDParam := c.Params("id")
	orderID, err := strconv.ParseUint(orderIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	var order models.Order
	if err := database.DB.Where("id = ? AND user_id IN (SELECT id FROM users WHERE tenant_id = ?)", orderID, tenantID).First(&order).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
	}

	var tenant models.Tenant
	if err := database.DB.Where("id = ?", tenantID).First(&tenant).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Tenant not found"})
	}

	var items []models.OrderItem
	database.DB.Where("order_id = ?", order.ID).Find(&items)

	// Trigger ERP Sync
	if err := services.SyncOrderToERP(tenant, order, items); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to sync with ERP",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"message": "Order successfully pushed to ERP (MRP)"})
}

// QueueERPSync - POST /api/v1/erp/sync
func QueueERPSync(c *fiber.Ctx) error {
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "ERP sync işi kuyruğa eklendi"})
}

// OdooWebhook - POST /api/v1/erp/odoo/webhook
func OdooWebhook(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "received", "erp": "odoo"})
}

// ERPNextWebhook - POST /api/v1/erp/erpnext/webhook
func ERPNextWebhook(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "received", "erp": "erpnext"})
}

// GetERPSyncStatus - GET /api/v1/erp/status/{jobID}
func GetERPSyncStatus(c *fiber.Ctx) error {
	jobID := c.Params("jobID")
	return c.JSON(fiber.Map{"status": "completed", "jobID": jobID})
}

// LogoSQLSync - POST /api/v1/erp/logo/sql/sync
func LogoSQLSync(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "queued", "erp": "logo_sql"})
}

// MikroSQLSync - POST /api/v1/erp/mikro/sql/sync
func MikroSQLSync(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "queued", "erp": "mikro_sql"})
}

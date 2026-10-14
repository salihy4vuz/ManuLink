package controllers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/database"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/models"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/services"
)

// GetQuotePDF generates and returns a PDF file for a specific order/quote
func GetQuotePDF(c *fiber.Ctx) error {
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
	if err := database.DB.Preload("Theme").Where("id = ?", tenantID).First(&tenant).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Tenant error"})
	}

	// Generate PDF
	pdfBytes, err := services.GenerateQuotePDF(tenant, order)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to generate PDF",
			"details": err.Error(),
		})
	}

	// Serve the PDF to the browser
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=\"Quote_"+orderIDParam+".pdf\"")
	
	// Add Audit Log
	services.LogAudit(c, "DOWNLOAD_PDF", "Order", uint(orderID))
	
	return c.Send(pdfBytes)
}

// GetMaterialCertificate generates and returns an EN 10204 3.1 PDF certificate
func GetMaterialCertificate(c *fiber.Ctx) error {
	orderIDParam := c.Params("id")

	// Dummy data for simulation
	cert := services.MaterialCertificate{
		CertificateNo: "CERT-" + orderIDParam,
		IssuedDate:    time.Now(),
		PartNo:        "PART-" + orderIDParam,
		Material:      "S235JR",
		HeatNo:        "H-123456",
		Thickness:     3.0,
	}

	pdfBytes, err := services.GenerateMaterialCertificatePDF(cert)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to generate certificate PDF",
			"details": err.Error(),
		})
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=\"Certificate_"+orderIDParam+".pdf\"")
	
	return c.Send(pdfBytes)
}

package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/models"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/services"
)

// GetPricingConfig godoc
// @Summary      Tenant fiyatlandırma konfigürasyonunu getir
// @Description  Admin kullanıcısı için tenant bazında tüm fiyatlandırma parametrelerini döndürür.
// @Tags         Admin
// @Produce      json
// @Param        tenant_id  query  int  false  "Tenant ID (default: JWT'den alınır)"
// @Success      200  {object}  models.PricingConfig
// @Router       /api/admin/pricing-config [get]
func GetPricingConfigHandler(c *fiber.Ctx) error {
	tenantID, err := getTenantIDFromCtx(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	cfg := services.GetPricingConfig(tenantID)
	return c.JSON(fiber.Map{
		"status": "ok",
		"config": cfg,
	})
}

// UpdatePricingConfig godoc
// @Summary      Tenant fiyatlandırma konfigürasyonunu güncelle
// @Description  Admin kullanıcısı için fiyatlandırma parametrelerini (punch maliyeti, gaz tipi, vb.) günceller.
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        config  body  models.PricingConfig  true  "Güncellenecek parametreler"
// @Success      200  {object}  fiber.Map
// @Failure      400  {object}  fiber.Map
// @Router       /api/admin/pricing-config [put]
func UpdatePricingConfigHandler(c *fiber.Ctx) error {
	tenantID, err := getTenantIDFromCtx(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var updated models.PricingConfig
	if err := c.BodyParser(&updated); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body: " + err.Error(),
		})
	}

	if err := services.UpdatePricingConfig(tenantID, updated); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update pricing config: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":    "updated",
		"tenant_id": tenantID,
		"message":   "Pricing configuration updated successfully. Changes will take effect immediately.",
	})
}

// getTenantIDFromCtx JWT locals veya query param'dan tenant ID alır.
func getTenantIDFromCtx(c *fiber.Ctx) (uint, error) {
	// Önce JWT locals'tan dene
	if tenantID, ok := c.Locals("tenant_id").(uint); ok && tenantID > 0 {
		return tenantID, nil
	}
	// Query param'dan dene (superadmin use case)
	if qID := c.Query("tenant_id"); qID != "" {
		id, err := strconv.ParseUint(qID, 10, 32)
		if err == nil && id > 0 {
			return uint(id), nil
		}
	}
	return 1, nil // default tenant
}

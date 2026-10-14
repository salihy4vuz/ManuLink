package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/config"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/models"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/repositories"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/services"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// GenerateQuote handles incoming requests for quoting sheet metal parts
func GenerateQuote(c *fiber.Ctx) error {
	var req services.PricingRequest

	// Parse JSON body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid JSON payload",
			"details": err.Error(),
		})
	}

	// For B2B requests, API key might be associated with a specific customer
	// we could extract customer ID from locals, but for now we proceed.

	// Calculate Quote using our massive pricing engine
	resp, err := services.CalculateQuote(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to calculate quote",
			"details": err.Error(),
		})
	}

	// Fetch TenantID
	tenantID, _ := c.Locals("TenantID").(uint)
	if tenantID == 0 {
		tenantID = 1 // default for now
	}

	// Prepare PricingDetails JSONB
	detailsJSON, _ := json.Marshal(resp)

	// Save to DB
	quoteRecord := models.Quote{
		TenantID:       tenantID,
		TotalCost:      resp.TotalCostUSD,
		FinalPrice:     resp.FinalSellingPrice,
		Currency:       resp.Currency,
		Status:         "draft",
		PricingDetails: datatypes.JSON(detailsJSON),
	}

	// Transactional Boundary (UoW)
	uow := repositories.NewUnitOfWork(config.DB)
	quoteRepo := repositories.NewQuoteRepository()

	err = uow.Do(c.Context(), func(tx *gorm.DB) error {
		return quoteRepo.Create(c.Context(), tx, &quoteRecord)
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to save quote to database via UoW",
		})
	}

	// Cache in Redis
	redisKey := fmt.Sprintf("%d", quoteRecord.ID)
	services.CachePriceCalculation(context.Background(), redisKey, quoteRecord)

	// Return the detailed pricing response
	return c.JSON(fiber.Map{
		"success": true,
		"quote_id": quoteRecord.ID,
		"quote":   resp,
	})
}

// ApproveQuote handles transitioning a Quote to Approved status and spinning off an MESOrder
func ApproveQuote(c *fiber.Ctx) error {
	quoteIDParam := c.Params("id")
	quoteID, err := strconv.ParseUint(quoteIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid quote ID"})
	}

	var quote models.Quote
	if err := config.DB.First(&quote, quoteID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Quote not found"})
	}

	if quote.Status == "approved" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Quote is already approved"})
	}

	// Transactional Boundary
	uow := repositories.NewUnitOfWork(config.DB)
	
	var mesOrder models.MESOrder

	err = uow.Do(c.Context(), func(tx *gorm.DB) error {
		// 1. Update Quote Status
		quote.Status = "approved"
		if err := tx.Save(&quote).Error; err != nil {
			return err
		}

		// 2. Generate Industrial Routing / Production Steps
		var routing []map[string]interface{}
		
		if quote.QuoteType == "machining" {
			// Extract MachineType from PricingDetails JSON
			var pricing map[string]interface{}
			json.Unmarshal(quote.PricingDetails, &pricing)
			
			machineType := "milling_3axis"
			if mt, ok := pricing["machine_type"].(string); ok {
				machineType = mt
			}

			if machineType == "lathe" {
				routing = []map[string]interface{}{
					{"operation": "CNC Torna (Turning)", "machine": "Mazak", "estimated_time_min": 30.0, "status": "pending"},
					{"operation": "Kalite Kontrol", "machine": "Kumpas/Mikrometre", "estimated_time_min": 5.0, "status": "pending"},
				}
			} else if machineType == "swiss_lathe" {
				routing = []map[string]interface{}{
					{"operation": "Kayar Otomat (Swiss Lathe)", "machine": "Citizen", "estimated_time_min": 15.0, "status": "pending"},
					{"operation": "Yıkama & Paketleme", "machine": "Ultrasonik", "estimated_time_min": 5.0, "status": "pending"},
				}
			} else if machineType == "welding" {
				// Metal Çatımı ve Kaynağı (Welding & Framework)
				routing = []map[string]interface{}{
					{"operation": "Metal Çatımı & Fikstür", "machine": "Montaj Masası", "estimated_time_min": 60.0, "status": "pending"},
					{"operation": "Gazaltı/TIG Kaynak", "machine": "Kaynak Kabini", "estimated_time_min": 120.0, "status": "pending"},
					{"operation": "Tahribatsız Muayene (NDT)", "machine": "Kalite Kontrol", "estimated_time_min": 20.0, "status": "pending"},
				}
			} else if machineType == "drilling" {
				// Delik Delme (Drilling)
				routing = []map[string]interface{}{
					{"operation": "Kılavuz Çekme & Delik Delme", "machine": "Sütunlu Matkap", "estimated_time_min": 10.0, "status": "pending"},
					{"operation": "Kontrol", "machine": "Mastarlama", "estimated_time_min": 2.0, "status": "pending"},
				}
			} else if machineType == "press" {
				// Eksantrik / Hidrolik Pres (Stamping)
				routing = []map[string]interface{}{
					{"operation": "Kalıp Bağlama", "machine": "Eksantrik Pres 100T", "estimated_time_min": 45.0, "status": "pending"},
					{"operation": "Pres Baskı", "machine": "Eksantrik Pres 100T", "estimated_time_min": 0.5, "status": "pending"}, // Seri üretim
				}
			} else {
				// milling_3axis
				routing = []map[string]interface{}{
					{"operation": "CNC Dik İşlem (Milling)", "machine": "Haas VF-2", "estimated_time_min": 45.0, "status": "pending"},
					{"operation": "Çapak Alma", "machine": "Manuel", "estimated_time_min": 5.0, "status": "pending"},
					{"operation": "Kalite Kontrol", "machine": "CMM", "estimated_time_min": 10.0, "status": "pending"},
				}
			}
		} else {
			routing = []map[string]interface{}{
				{"operation": "Lazer Kesim", "machine": "Bystronic 10kW", "estimated_time_min": 12.5, "status": "pending"},
				{"operation": "Abkant Pres (Bending)", "machine": "Trumpf TruBend", "estimated_time_min": 4.0, "status": "pending"},
				{"operation": "Kalite Kontrol", "machine": "Manuel", "estimated_time_min": 2.0, "status": "pending"},
			}
		}
		
		routingJSON, _ := json.Marshal(routing)

		// 3. Create MESOrder
		mesOrder = models.MESOrder{
			TenantID:        quote.TenantID,
			QuoteID:         quote.ID,
			Status:          "queued",
			ProductionSteps: datatypes.JSON(routingJSON),
		}

		if err := tx.Create(&mesOrder).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to approve quote and generate MES order",
			"details": err.Error(),
		})
	}

	// 4. Dispatch Webhook to ERP (Asynchronous)
	var tenant models.Tenant
	if err := config.DB.First(&tenant, quote.TenantID).Error; err != nil {
		tenant.ERPEndpoint = "http://localhost:3000/api/mock-erp/webhook" // Fallback if missing
	}
	services.DispatchMESOrderToERP(mesOrder, tenant.ERPEndpoint)

	return c.JSON(fiber.Map{
		"message": "Quote approved successfully. MES Order generated and dispatched to ERP.",
		"mes_order_id": mesOrder.ID,
	})
}

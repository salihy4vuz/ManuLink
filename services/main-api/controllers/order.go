package controllers

import (
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/database"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/models"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/services"
)

type PlaceOrderRequest struct {
	CustomerID        uint                   `json:"customer_id"`
	QuoteData         services.PricingRequest `json:"quote_data"`
	PartName          string                 `json:"part_name"`
	ShippingAddress   string                 `json:"shipping_address"`
}

// PlaceOrder converts a pricing request directly into a pending order
func PlaceOrder(c *fiber.Ctx) error {
	var req PlaceOrderRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	// Recalculate quote to ensure validity
	resp, err := services.CalculateQuote(req.QuoteData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to calculate quote for order"})
	}

	// Create Order
	order := models.Order{
		CustomerID:        req.CustomerID,
		TotalAmountUSD:    resp.TotalCostUSD,
		FinalSellingPrice: resp.FinalSellingPrice,
		Currency:          resp.Currency,
		Status:            models.StatusPendingPayment,
		ShippingAddress:   req.ShippingAddress,
		EstimatedDelivery: time.Now().AddDate(0, 0, 7), // 7 days default
	}

	// Use Transaction
	tx := database.DB.Begin()
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create order"})
	}

	// Create OrderItem
	orderItem := models.OrderItem{
		OrderID:         order.ID,
		PartName:        req.PartName,
		Material:        req.QuoteData.Material,
		ThicknessMm:     req.QuoteData.ThicknessMm,
		Quantity:        req.QuoteData.Quantity,
		UnitCostUSD:     resp.UnitCostUSD,
		TotalCostUSD:    resp.TotalCostUSD,
		LaserCutting:    true,
		BendingCount:    req.QuoteData.BendCount,
		ComplexBending:  req.QuoteData.ComplexBending,
		TappingCount:    req.QuoteData.TappingCount,
	}

	if err := tx.Create(&orderItem).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create order item"})
	}
	
	tx.Commit()

	// Optionally simulate payment via stripe mock
	payment, err := services.ProcessPayment(req.CustomerID, order.ID, resp.FinalSellingPrice, resp.Currency)
	if err == nil && payment != nil {
		database.DB.Model(&order).Update("status", models.StatusInProduction)
		
		// Send notification
		services.NotifyOrderCreated("customer@example.com", "+1234567890", order.ID)
		
		// Trigger webhook
		services.DispatchWebhookToTenants(req.CustomerID, "order.created", order)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"order_id": order.ID,
		"status":   order.Status,
	})
}

// UpdateOrderStatus allows admin to update the order status
func UpdateOrderStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	
	type UpdateRequest struct {
		Status string `json:"status"`
	}
	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "order not found"})
	}

	oldStatus := order.Status
	order.Status = models.OrderStatus(req.Status)
	database.DB.Save(&order)

	// Trigger webhook on status change
	if oldStatus != order.Status {
		services.DispatchWebhookToTenants(order.CustomerID, "order.status_updated", order)
	}

	return c.JSON(fiber.Map{"success": true, "order": order})
}



// GetOrders retrieves the list of orders for the authenticated B2B tenant
func GetOrders(c *fiber.Ctx) error {
	// In production, extract TenantID from API Key
	var orders []models.Order
	
	// Preload nested items
	if err := database.DB.Preload("Items").Find(&orders).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch orders"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"orders": orders,
	})
}

// GetStats returns summary statistics for the admin dashboard
func GetStats(c *fiber.Ctx) error {
	var activeOrders int64
	var completedOrders int64

	database.DB.Model(&models.Order{}).Where("status != ?", "Completed").Count(&activeOrders)
	database.DB.Model(&models.Order{}).Where("status = ?", "Completed").Count(&completedOrders)

	var totalRevenue float64
	database.DB.Model(&models.Order{}).Where("status = ?", "Completed").Select("COALESCE(SUM(final_selling_price), 0)").Scan(&totalRevenue)

	return c.JSON(fiber.Map{
		"success": true,
		"active_orders": activeOrders,
		"completed_orders": completedOrders,
		"total_revenue": totalRevenue, 
	})
}

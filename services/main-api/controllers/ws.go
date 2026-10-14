package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/services"
)

// WsTrackingHandler manages the WebSocket upgrades and client loop
func WsTrackingHandler(c *websocket.Conn) {
	// Normally we would extract the TenantID from a JWT token sent in query or headers during upgrade.
	// For this phase, we use locals populated by middleware before the upgrade.
	tenantIDRaw := c.Locals("tenant_id")
	if tenantIDRaw == nil {
		log.Println("[WebSocket] Missing tenant_id in locals")
		return
	}
	tenantID := tenantIDRaw.(uint)

	client := &services.Client{
		ID:       uuid.New().String(),
		TenantID: tenantID,
		Conn:     c,
		Send:     make(chan []byte, 256),
	}

	services.DefaultHub.Register <- client

	// Pump messages to the WebSocket client
	go func() {
		defer func() {
			services.DefaultHub.Unregister <- client
			c.Close()
		}()
		for {
			message, ok := <-client.Send
			if !ok {
				c.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		}
	}()

	// Read loop (to keep connection alive and catch disconnects)
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
	}

	services.DefaultHub.Unregister <- client
}

// InternalUpdateTracking is a mock endpoint to broadcast a status update to a Tenant's WS clients
func InternalUpdateTracking(c *fiber.Ctx) error {
	tenantIDRaw := c.Locals("tenant_id")
	if tenantIDRaw == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	tenantID := tenantIDRaw.(uint)

	type UpdatePayload struct {
		OrderID int    `json:"order_id"`
		Status  string `json:"status"` // e.g. "cutting", "bending", "done"
	}

	var payload UpdatePayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	// Broadcast via Hub
	services.DefaultHub.BroadcastToTenant(tenantID, services.WsMessage{
		Type:    "production_update",
		Payload: payload,
	})

	return c.JSON(fiber.Map{"message": "Broadcast sent to tenant clients"})
}

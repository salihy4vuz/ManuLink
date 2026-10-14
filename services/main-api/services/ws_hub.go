package services

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
)

// WsMessage represents a standard JSON message sent to clients
type WsMessage struct {
	Type    string      `json:"type"`     // e.g., "status_update", "nesting_complete"
	Payload interface{} `json:"payload"`
}

// Client represents a connected WebSocket client linked to a specific Tenant
type Client struct {
	ID       string
	TenantID uint
	Conn     *websocket.Conn
	Send     chan []byte
}

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients mapped by TenantID to ensure data isolation
	// TenantID -> map of client pointers
	Clients map[uint]map[*Client]bool

	// Inbound messages from the clients (if we need to handle incoming WS messages).
	Broadcast chan []byte

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client

	mu sync.RWMutex
}

// Global Hub instance
var DefaultHub *Hub

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[uint]map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	log.Println("[WebSocket] Hub is running...")
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if h.Clients[client.TenantID] == nil {
				h.Clients[client.TenantID] = make(map[*Client]bool)
			}
			h.Clients[client.TenantID][client] = true
			h.mu.Unlock()
			log.Printf("[WebSocket] Client %s registered to Tenant %d\n", client.ID, client.TenantID)

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.TenantID][client]; ok {
				delete(h.Clients[client.TenantID], client)
				close(client.Send)
				if len(h.Clients[client.TenantID]) == 0 {
					delete(h.Clients, client.TenantID)
				}
				log.Printf("[WebSocket] Client %s unregistered from Tenant %d\n", client.ID, client.TenantID)
			}
			h.mu.Unlock()

		case _ = <-h.Broadcast:
			// For generic broadcasts to ALL tenants (rarely used in B2B).
			// Use BroadcastToTenant instead for data isolation.
		}
	}
}

// BroadcastToTenant sends a message only to clients belonging to the specified Tenant
func (h *Hub) BroadcastToTenant(tenantID uint, msg WsMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("[WebSocket] JSON marshal error:", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	if tenantClients, ok := h.Clients[tenantID]; ok {
		for client := range tenantClients {
			select {
			case client.Send <- data:
			default:
				close(client.Send)
				delete(h.Clients[tenantID], client)
			}
		}
	}
}

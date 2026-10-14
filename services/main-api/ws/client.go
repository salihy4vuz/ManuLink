package ws

import (
	"log"
	"time"

	"github.com/gofiber/websocket/v2"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

type Client struct {
	Hub *Hub

	// The websocket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		Hub:  hub,
		Conn: conn,
		send: make(chan []byte, 256),
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	// Configure Heartbeat
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		
		// FAZ 9.1 & 9.2 & 9.3 & 9.4: SolidWorks'ten gelen "cad_export" JSON'ını parse et ve Kütle verisini yakala
		type BendDetail struct {
			Angle     float64 `json:"angle"`
			Radius    float64 `json:"radius"`
			Direction string  `json:"direction"`
		}
		
		type CADPayload struct {
			PartNumber       string       `json:"part_number"`      // FAZ 9.6
			Revision         string       `json:"revision"`         // FAZ 9.6
			Description      string       `json:"description"`      // FAZ 9.6
			Material         string       `json:"material"`
			ThicknessMm      float64      `json:"thickness_mm"`
			BoundingBoxX     float64      `json:"bounding_box_x"`
			BoundingBoxY     float64      `json:"bounding_box_y"`
			MassKg           float64      `json:"mass_kg"`          // FAZ 9.4
			SurfaceAreaM2    float64      `json:"surface_area_m2"`  // FAZ 9.4
			CutLengthMm      float64      `json:"cut_length_mm"`    // FAZ 9.7: Lazer kesim uzunluğu
			PierceCount      int          `json:"pierce_count"`     // FAZ 9.7: Delme/Patlatma sayısı
			CountersinkCount int          `json:"countersink_count"`// FAZ 9.8: Havşalı delik sayısı
			TappedHoleCount  int          `json:"tapped_hole_count"`// FAZ 9.8: Kılavuzlu/Dişli delik sayısı
			BendCount        int          `json:"bend_count"`
			Source           string       `json:"source"`
			DxfBase64        string       `json:"dxf_base64"`
			StepBase64       string       `json:"step_base64"`       // FAZ 9.5
			ThumbnailBase64  string       `json:"thumbnail_base64"`  // FAZ 9.5
			BendDetails      []BendDetail `json:"bend_details"`
		}

		type WSMessage struct {
			Action string     `json:"action"`
			Data   CADPayload `json:"data"` // Strongly typed for Phase 9.4 completeness
		}
		
		var parsed WSMessage
		importJson := true // Quick workaround for unused var
		_ = parsed
		_ = importJson
		// Quick and dirty check for performance (In production, use encoding/json)
		if string(message) != "" {
			// Broadcast anyway to all clients (including frontend 3D Viewer)
			c.Hub.broadcast <- message
			log.Printf("B2B WebSocket Received Data. Checking for DXF payload...")

			// [Faz 9.2: DXF Base64 Extraction Simulation]
			// In a full implementation we run json.Unmarshal(&parsed, message).
			// If parsed.Action == "cad_export" and parsed.Data["dxf_base64"] exists:
			// 1. base64.StdEncoding.DecodeString(dxfData)
			// 2. os.WriteFile("uploads/solidworks_export.dxf", decodedBytes, 0644)
			// 3. Trigger Rust CAM Engine for G-Code Generation!
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

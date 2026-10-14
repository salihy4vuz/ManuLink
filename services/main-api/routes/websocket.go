package routes

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// InitWebSocket sets up the CAD bridge websocket endpoints
func InitWebSocket(app *fiber.App) {
	// Middleware to check if it's a websocket upgrade request
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/cad-bridge", websocket.New(func(c *websocket.Conn) {
		log.Println("CAD Bridge WebSocket Connected:", c.RemoteAddr())

		var (
			msgType int
			msg     []byte
			err     error
		)

		for {
			if msgType, msg, err = c.ReadMessage(); err != nil {
				log.Println("CAD Bridge WebSocket Disconnected:", err)
				break
			}
			log.Printf("Received message: %s\n", msg)

			// Echo back or process
			if err = c.WriteMessage(msgType, msg); err != nil {
				log.Println("write:", err)
				break
			}
		}
	}))
}

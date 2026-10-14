package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/services"
)

type NewTicketRequest struct {
	Subject  string `json:"subject"`
	Message  string `json:"message"`
	Priority string `json:"priority"`
}

// GetUserTickets returns all tickets for the logged in user
func GetUserTickets(c *fiber.Ctx) error {
	username := c.Get("X-User", "admin") // Mock JWT
	
	tickets := services.GetTicketsByUser(username)
	
	return c.JSON(fiber.Map{
		"status":  "success",
		"tickets": tickets,
	})
}

// CreateTicket creates a new support ticket
func CreateTicket(c *fiber.Ctx) error {
	username := c.Get("X-User", "admin") // Mock JWT
	
	var req NewTicketRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Subject == "" || req.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Subject and message are required"})
	}

	priority := services.TicketPriority(req.Priority)
	if priority != services.PriorityLow && priority != services.PriorityMedium && priority != services.PriorityHigh {
		priority = services.PriorityMedium
	}

	ticket := services.CreateTicket(username, req.Subject, req.Message, priority)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Ticket created successfully",
		"ticket":  ticket,
	})
}

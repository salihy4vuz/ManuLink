package services

import (
	"sync"
	"time"

	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/database"
)

type TicketStatus string

const (
	StatusOpen       TicketStatus = "Open"
	StatusInProgress TicketStatus = "In Progress"
	StatusResolved   TicketStatus = "Resolved"
)

type TicketPriority string

const (
	PriorityLow    TicketPriority = "Low"
	PriorityMedium TicketPriority = "Medium"
	PriorityHigh   TicketPriority = "High"
)

type Ticket struct {
	ID        int            `json:"id"`
	Username  string         `json:"username"`
	Subject   string         `json:"subject"`
	Message   string         `json:"message"`
	Status    TicketStatus   `json:"status"`
	Priority  TicketPriority `json:"priority"`
	CreatedAt time.Time      `json:"created_at"`
}



var (
	tickets     = make([]Ticket, 0)
	ticketIDSeq = 1
	ticketMutex sync.RWMutex
)

func init() {
	// Try loading from JSON DB
	_ = database.Load(nil, &tickets)
	
	if len(tickets) == 0 {
		tickets = append(tickets, Ticket{
			ID:        ticketIDSeq,
			Username:  "admin",
			Subject:   "G-Code generated from DXF is misaligned",
			Message:   "The laser cutting origin seems to be shifted by 5mm.",
			Status:    StatusResolved,
			Priority:  PriorityHigh,
			CreatedAt: time.Now().Add(-48 * time.Hour),
		})
	}
	
	for _, t := range tickets {
		if t.ID >= ticketIDSeq {
			ticketIDSeq = t.ID + 1
		}
	}
}

func CreateTicket(username, subject, message string, priority TicketPriority) Ticket {
	ticketMutex.Lock()
	defer ticketMutex.Unlock()

	t := Ticket{
		ID:        ticketIDSeq,
		Username:  username,
		Subject:   subject,
		Message:   message,
		Status:    StatusOpen,
		Priority:  priority,
		CreatedAt: time.Now(),
	}
	tickets = append(tickets, t)
	ticketIDSeq++

	// Persist
	go saveDB()

	// Faz 1.6: Webhook Notification (e.g. to Discord/Slack)
	DispatchWebhookEvent("TICKET_CREATED", t, "http://localhost:3000/mock/discord/webhook")

	return t
}

func saveDB() {
	ticketMutex.RLock()
	defer ticketMutex.RUnlock()
	// Create a copy of slice for interface casting
	var tIf []interface{}
	for _, t := range tickets {
		tIf = append(tIf, t)
	}
	_ = database.Save(nil, tIf) // In real code, we'd pass both
}

func GetTicketsByUser(username string) []Ticket {
	ticketMutex.RLock()
	defer ticketMutex.RUnlock()

	var userTickets []Ticket
	for _, t := range tickets {
		if t.Username == username {
			userTickets = append(userTickets, t)
		}
	}

	// Reverse to show newest first
	for i, j := 0, len(userTickets)-1; i < j; i, j = i+1, j-1 {
		userTickets[i], userTickets[j] = userTickets[j], userTickets[i]
	}

	return userTickets
}

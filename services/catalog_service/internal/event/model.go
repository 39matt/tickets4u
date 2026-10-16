package event

import "time"

type Event struct {
	ID               string    `json:"id"`
	Title            string    `json:"title"`
	TotalTickets     int       `json:"total_tickets"`
	AvailableTickets int       `json:"available_tickets"`
	PriceCents       int       `json:"price_cents"`
	CreatedAt        time.Time `json:"created_at"`
}

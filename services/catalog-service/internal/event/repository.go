package event

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrNotFound = fmt.Errorf("event not found")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetEvents(ctx context.Context) ([]Event, error) {
	query := `
		SELECT id, title, total_tickets, available_tickets, price_cents, created_at
		FROM events
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	events := make([]Event, 0)
	for rows.Next() {
		var e Event
		err = rows.Scan(&e.ID, &e.Title, &e.TotalTickets, &e.AvailableTickets, &e.PriceCents, &e.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, Event{})

	}

	return events, nil
}

func (r *Repository) GetEvent(ctx context.Context, id string) (*Event, error) {
	query := `SELECT * FROM events WHERE id = $1`

	var event Event
	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&event.ID, &event.Title, &event.TotalTickets, &event.AvailableTickets, &event.PriceCents, &event.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to query event: %w", err)
	}

	return &event, nil
}

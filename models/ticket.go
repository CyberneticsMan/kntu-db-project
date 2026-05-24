package models

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Ticket struct {
	ID        int     `json:"id"`
	EventName string  `json:"event_name"`
	Price     float64 `json:"price"`
}

func GetTicketByID(db *pgxpool.Pool, id int) (*Ticket, error) {
	var t Ticket
	query := `SELECT id, event_name, price FROM tickets WHERE id = $1`

	err := db.QueryRow(context.Background(), query, id).Scan(&t.ID, &t.EventName, &t.Price)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

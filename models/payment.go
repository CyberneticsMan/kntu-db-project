package models

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Payment struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	ReservationID int       `json:"reservation_id"`
	Amount        float64   `json:"amount"`
	Method        string    `json:"method"`
	Status        string    `json:"status"`
	ProcessedAt   time.Time `json:"processed_at"`
}

func (p *Payment) Validate() error {
	if p.UserID <= 0 {
		return fmt.Errorf("user id is required")
	}
	if p.ReservationID <= 0 {
		return fmt.Errorf("reservation id is required")
	}
	if p.Amount < 0 {
		return fmt.Errorf("amount cannot be negative")
	}
	if p.Method == "" {
		return fmt.Errorf("payment method is required")
	}
	return nil
}

func CreatePayment(db *pgxpool.Pool, payment *Payment) (*Payment, error) {
	if err := payment.Validate(); err != nil {
		return nil, err
	}
	query := `
		INSERT INTO payments (user_id, reservation_id, amount, method, status, processed_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, user_id, reservation_id, amount, method, status, processed_at`
	if err := db.QueryRow(context.Background(), query, payment.UserID, payment.ReservationID, payment.Amount, payment.Method, payment.Status).Scan(
		&payment.ID, &payment.UserID, &payment.ReservationID, &payment.Amount, &payment.Method, &payment.Status, &payment.ProcessedAt,
	); err != nil {
		return nil, err
	}
	return payment, nil
}

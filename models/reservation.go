package models

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Reservation struct {
	ID          int        `json:"id"`
	UserID      int        `json:"user_id"`
	TicketID    int        `json:"ticket_id"`
	Status      string     `json:"status"`
	ReservedAt  time.Time  `json:"reserved_at"`
	ExpiresAt   time.Time  `json:"expires_at"`
	PaidAt      *time.Time `json:"paid_at,omitempty"`
	CancelledAt *time.Time `json:"cancelled_at,omitempty"`
}

func (r *Reservation) Validate() error {
	if r.UserID <= 0 {
		return fmt.Errorf("user id is required")
	}
	if r.TicketID <= 0 {
		return fmt.Errorf("ticket id is required")
	}
	if r.ExpiresAt.Before(r.ReservedAt) {
		return fmt.Errorf("expires_at cannot be before reserved_at")
	}
	return nil
}

func scanReservation(row interface{ Scan(dest ...any) error }) (*Reservation, error) {
	var r Reservation
	if err := row.Scan(&r.ID, &r.UserID, &r.TicketID, &r.Status, &r.ReservedAt, &r.ExpiresAt, &r.PaidAt, &r.CancelledAt); err != nil {
		return nil, err
	}
	return &r, nil
}

func CreateReservation(db *pgxpool.Pool, reservation *Reservation) (*Reservation, error) {
	if err := reservation.Validate(); err != nil {
		return nil, err
	}

	tx, err := db.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	var remaining int
	if err := tx.QueryRow(context.Background(), `SELECT remaining_capacity FROM tickets WHERE id = $1 FOR UPDATE`, reservation.TicketID).Scan(&remaining); err != nil {
		return nil, err
	}
	if remaining <= 0 {
		return nil, fmt.Errorf("ticket is sold out")
	}

	if _, err := tx.Exec(context.Background(), `UPDATE tickets SET remaining_capacity = remaining_capacity - 1 WHERE id = $1`, reservation.TicketID); err != nil {
		return nil, err
	}

	query := `
		INSERT INTO reservations (user_id, ticket_id, status, reserved_at, expires_at)
		VALUES ($1, $2, 'reserved', NOW(), $3)
		RETURNING id, user_id, ticket_id, status, reserved_at, expires_at, paid_at, cancelled_at`
	if err := tx.QueryRow(context.Background(), query, reservation.UserID, reservation.TicketID, reservation.ExpiresAt).Scan(
		&reservation.ID, &reservation.UserID, &reservation.TicketID, &reservation.Status, &reservation.ReservedAt, &reservation.ExpiresAt, &reservation.PaidAt, &reservation.CancelledAt,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return nil, err
	}
	return reservation, nil
}

func GetReservationByID(db *pgxpool.Pool, id int) (*Reservation, error) {
	query := `SELECT id, user_id, ticket_id, status, reserved_at, expires_at, paid_at, cancelled_at FROM reservations WHERE id = $1`
	return scanReservation(db.QueryRow(context.Background(), query, id))
}

func ListReservationsByUser(db *pgxpool.Pool, userID int) ([]*Reservation, error) {
	query := `SELECT id, user_id, ticket_id, status, reserved_at, expires_at, paid_at, cancelled_at FROM reservations WHERE user_id = $1 ORDER BY reserved_at DESC`
	rows, err := db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reservations := make([]*Reservation, 0)
	for rows.Next() {
		r, err := scanReservation(rows)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return reservations, nil
}

func CancelReservation(db *pgxpool.Pool, id int) (*Reservation, error) {
	tx, err := db.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	reservation, err := scanReservation(tx.QueryRow(context.Background(), `SELECT id, user_id, ticket_id, status, reserved_at, expires_at, paid_at, cancelled_at FROM reservations WHERE id = $1 FOR UPDATE`, id))
	if err != nil {
		return nil, err
	}
	if reservation.Status == "cancelled" {
		return reservation, nil
	}

	if _, err := tx.Exec(context.Background(), `UPDATE reservations SET status = 'cancelled', cancelled_at = NOW() WHERE id = $1`, id); err != nil {
		return nil, err
	}
	if _, err := tx.Exec(context.Background(), `UPDATE tickets SET remaining_capacity = remaining_capacity + 1 WHERE id = $1`, reservation.TicketID); err != nil {
		return nil, err
	}
	reservation.Status = "cancelled"
	now := time.Now()
	reservation.CancelledAt = &now

	if err := tx.Commit(context.Background()); err != nil {
		return nil, err
	}
	return reservation, nil
}

func MarkReservationPaid(db *pgxpool.Pool, id int) (*Reservation, error) {
	query := `UPDATE reservations SET status = 'paid', paid_at = NOW() WHERE id = $1 AND status = 'reserved' RETURNING id, user_id, ticket_id, status, reserved_at, expires_at, paid_at, cancelled_at`
	return scanReservation(db.QueryRow(context.Background(), query, id))
}

func ExpireReservation(db *pgxpool.Pool, id int) error {
	_, err := db.Exec(context.Background(), `UPDATE reservations SET status = 'expired' WHERE id = $1`, id)
	if err != nil {
		return err
	}
	_, err = db.Exec(context.Background(), `UPDATE tickets SET remaining_capacity = remaining_capacity + 1 WHERE id = (SELECT ticket_id FROM reservations WHERE id = $1)`, id)
	return err
}

func FindExpiredReservations(db *pgxpool.Pool) ([]*Reservation, error) {
	rows, err := db.Query(context.Background(), `SELECT id, user_id, ticket_id, status, reserved_at, expires_at, paid_at, cancelled_at FROM reservations WHERE status = 'reserved' AND expires_at <= NOW()`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reservations := make([]*Reservation, 0)
	for rows.Next() {
		r, err := scanReservation(rows)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, r)
	}
	return reservations, rows.Err()
}

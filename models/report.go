package models

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Report struct {
	ID            int        `json:"id"`
	UserID        int        `json:"user_id"`
	TicketID      *int       `json:"ticket_id,omitempty"`
	ReservationID *int       `json:"reservation_id,omitempty"`
	Category      string     `json:"category"`
	Subject       string     `json:"subject"`
	Body          string     `json:"body"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
	ResolvedAt    *time.Time `json:"resolved_at,omitempty"`
}

func (r *Report) Validate() error {
	if r.UserID <= 0 {
		return fmt.Errorf("user id is required")
	}
	if r.Category == "" {
		return fmt.Errorf("category is required")
	}
	if r.Subject == "" {
		return fmt.Errorf("subject is required")
	}
	if r.Body == "" {
		return fmt.Errorf("body is required")
	}
	return nil
}

func CreateReport(db *pgxpool.Pool, report *Report) (*Report, error) {
	if err := report.Validate(); err != nil {
		return nil, err
	}

	query := `
		INSERT INTO reports (user_id, ticket_id, reservation_id, category, subject, body, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		RETURNING id, user_id, ticket_id, reservation_id, category, subject, body, status, created_at, resolved_at`
	if err := db.QueryRow(context.Background(), query, report.UserID, report.TicketID, report.ReservationID, report.Category, report.Subject, report.Body, report.Status).Scan(
		&report.ID, &report.UserID, &report.TicketID, &report.ReservationID, &report.Category, &report.Subject, &report.Body, &report.Status, &report.CreatedAt, &report.ResolvedAt,
	); err != nil {
		return nil, err
	}
	return report, nil
}

func scanReport(row interface{ Scan(dest ...any) error }) (*Report, error) {
	var report Report
	if err := row.Scan(&report.ID, &report.UserID, &report.TicketID, &report.ReservationID, &report.Category, &report.Subject, &report.Body, &report.Status, &report.CreatedAt, &report.ResolvedAt); err != nil {
		return nil, err
	}
	return &report, nil
}

func GetReportByID(db *pgxpool.Pool, id int) (*Report, error) {
	query := `SELECT id, user_id, ticket_id, reservation_id, category, subject, body, status, created_at, resolved_at FROM reports WHERE id = $1`
	return scanReport(db.QueryRow(context.Background(), query, id))
}

func ListReportsByUser(db *pgxpool.Pool, userID int) ([]*Report, error) {
	query := `SELECT id, user_id, ticket_id, reservation_id, category, subject, body, status, created_at, resolved_at FROM reports WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reports := make([]*Report, 0)
	for rows.Next() {
		report, err := scanReport(rows)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return reports, nil
}

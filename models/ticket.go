package models

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Ticket struct {
	ID                int       `json:"id"`
	EventName         string    `json:"event_name"`
	SportType         string    `json:"sport_type"`
	HomeTeam          string    `json:"home_team"`
	AwayTeam          string    `json:"away_team"`
	Venue             string    `json:"venue"`
	City              string    `json:"city"`
	EventTime         time.Time `json:"event_time"`
	Price             float64   `json:"price"`
	RemainingCapacity int       `json:"remaining_capacity"`
	Category          string    `json:"category"`
}

type TicketFilters struct {
	SportType string
	HomeTeam  string
	AwayTeam  string
	City      string
	Venue     string
	Category  string
	Search    string
	MinPrice  *float64
	MaxPrice  *float64
}

func (t *Ticket) Validate() error {
	if t.EventName == "" {
		return fmt.Errorf("event name is required")
	}
	if t.SportType == "" {
		return fmt.Errorf("sport type is required")
	}
	if t.HomeTeam == "" {
		return fmt.Errorf("home team is required")
	}
	if t.AwayTeam == "" {
		return fmt.Errorf("away team is required")
	}
	if t.Venue == "" {
		return fmt.Errorf("venue is required")
	}
	if t.City == "" {
		return fmt.Errorf("city is required")
	}
	if t.Price < 0 {
		return fmt.Errorf("price cannot be negative")
	}
	if t.RemainingCapacity < 0 {
		return fmt.Errorf("remaining capacity cannot be negative")
	}
	return nil
}

func scanTicket(row Scanner) (*Ticket, error) {
	var ticket Ticket
	if err := row.Scan(
		&ticket.ID,
		&ticket.EventName,
		&ticket.SportType,
		&ticket.HomeTeam,
		&ticket.AwayTeam,
		&ticket.Venue,
		&ticket.City,
		&ticket.EventTime,
		&ticket.Price,
		&ticket.RemainingCapacity,
		&ticket.Category,
	); err != nil {
		return nil, err
	}
	return &ticket, nil
}

type Scanner interface {
	Scan(dest ...any) error
}

func GetTicketByID(db *pgxpool.Pool, id int) (*Ticket, error) {
	query := `SELECT id, event_name, sport_type, home_team, away_team, venue, city, event_time, price, remaining_capacity, category FROM tickets WHERE id = $1`
	return scanTicket(db.QueryRow(context.Background(), query, id))
}

func ListTickets(db *pgxpool.Pool, filters TicketFilters) ([]*Ticket, error) {
	query := `SELECT id, event_name, sport_type, home_team, away_team, venue, city, event_time, price, remaining_capacity, category FROM tickets WHERE 1=1`
	args := make([]any, 0)
	appendFilter := func(condition string, value any) {
		args = append(args, value)
		query += fmt.Sprintf(condition, len(args))
	}

	if filters.SportType != "" {
		appendFilter(" AND sport_type = $%d", filters.SportType)
	}
	if filters.HomeTeam != "" {
		appendFilter(" AND home_team ILIKE $%d", "%"+filters.HomeTeam+"%")
	}
	if filters.AwayTeam != "" {
		appendFilter(" AND away_team ILIKE $%d", "%"+filters.AwayTeam+"%")
	}
	if filters.City != "" {
		appendFilter(" AND city ILIKE $%d", "%"+filters.City+"%")
	}
	if filters.Venue != "" {
		appendFilter(" AND venue ILIKE $%d", "%"+filters.Venue+"%")
	}
	if filters.Category != "" {
		appendFilter(" AND category ILIKE $%d", "%"+filters.Category+"%")
	}
	if filters.Search != "" {
		search := "%" + filters.Search + "%"
		args = append(args, search, search, search, search, search)
		query += fmt.Sprintf(` AND (event_name ILIKE $%d OR home_team ILIKE $%d OR away_team ILIKE $%d OR venue ILIKE $%d OR category ILIKE $%d)`, len(args)-4, len(args)-3, len(args)-2, len(args)-1, len(args))
	}
	if filters.MinPrice != nil {
		appendFilter(" AND price >= $%d", *filters.MinPrice)
	}
	if filters.MaxPrice != nil {
		appendFilter(" AND price <= $%d", *filters.MaxPrice)
	}

	query += ` ORDER BY event_time ASC, price ASC`

	rows, err := db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tickets := make([]*Ticket, 0)
	for rows.Next() {
		ticket, err := scanTicket(rows)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tickets, nil
}

func CreateTicket(db *pgxpool.Pool, ticket *Ticket) (*Ticket, error) {
	if err := ticket.Validate(); err != nil {
		return nil, err
	}

	query := `
		INSERT INTO tickets (event_name, sport_type, home_team, away_team, venue, city, event_time, price, remaining_capacity, category)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING id`
	if err := db.QueryRow(context.Background(), query,
		ticket.EventName, ticket.SportType, ticket.HomeTeam, ticket.AwayTeam, ticket.Venue, ticket.City, ticket.EventTime, ticket.Price, ticket.RemainingCapacity, ticket.Category,
	).Scan(&ticket.ID); err != nil {
		return nil, err
	}
	return ticket, nil
}

func UpdateTicket(db *pgxpool.Pool, ticket *Ticket) error {
	if err := ticket.Validate(); err != nil {
		return err
	}

	query := `
		UPDATE tickets
		SET event_name = $1, sport_type = $2, home_team = $3, away_team = $4, venue = $5, city = $6, event_time = $7, price = $8, remaining_capacity = $9, category = $10
		WHERE id = $11`
	_, err := db.Exec(context.Background(), query,
		ticket.EventName, ticket.SportType, ticket.HomeTeam, ticket.AwayTeam, ticket.Venue, ticket.City, ticket.EventTime, ticket.Price, ticket.RemainingCapacity, ticket.Category, ticket.ID,
	)
	return err
}

func DeleteTicket(db *pgxpool.Pool, id int) error {
	_, err := db.Exec(context.Background(), `DELETE FROM tickets WHERE id = $1`, id)
	return err
}

func AdjustTicketCapacity(db *pgxpool.Pool, id int, delta int) error {
	_, err := db.Exec(context.Background(), `UPDATE tickets SET remaining_capacity = remaining_capacity + $1 WHERE id = $2`, delta, id)
	return err
}

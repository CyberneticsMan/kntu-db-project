package services

import (
	"github.com/CyberneticsMan/kntu-db-project/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TicketService struct {
	DB *pgxpool.Pool
}

func (ts *TicketService) GetTicketByID(id int) (*models.Ticket, error) {
	return models.GetTicketByID(ts.DB, id)
}

func (ts *TicketService) ListTickets(filters models.TicketFilters) ([]*models.Ticket, error) {
	return models.ListTickets(ts.DB, filters)
}

func (ts *TicketService) CreateTicket(ticket *models.Ticket) (*models.Ticket, error) {
	return models.CreateTicket(ts.DB, ticket)
}

func (ts *TicketService) UpdateTicket(ticket *models.Ticket) error {
	return models.UpdateTicket(ts.DB, ticket)
}

func (ts *TicketService) DeleteTicket(id int) error {
	return models.DeleteTicket(ts.DB, id)
}

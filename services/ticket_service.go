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

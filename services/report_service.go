package services

import (
	"github.com/CyberneticsMan/kntu-db-project/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReportService struct {
	DB *pgxpool.Pool
}

func (rs *ReportService) CreateReport(req *models.CreateReportRequest) (*models.Report, error) {
	report := &models.Report{
		UserID:        req.UserID,
		TicketID:      req.TicketID,
		ReservationID: req.ReservationID,
		Category:      req.Category,
		Subject:       req.Subject,
		Body:          req.Body,
		Status:        "pending",
	}
	return models.CreateReport(rs.DB, report)
}

func (rs *ReportService) GetReportByID(id int) (*models.Report, error) {
	return models.GetReportByID(rs.DB, id)
}

func (rs *ReportService) ListReportsByUser(userID int) ([]*models.Report, error) {
	return models.ListReportsByUser(rs.DB, userID)
}

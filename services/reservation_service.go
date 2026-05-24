package services

import (
	"time"

	"github.com/CyberneticsMan/kntu-db-project/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReservationService struct {
	DB *pgxpool.Pool
}

func (rs *ReservationService) CreateReservation(req *models.CreateReservationRequest) (*models.Reservation, error) {
	minutes := req.Minutes
	if minutes <= 0 {
		minutes = 10
	}
	reservation := &models.Reservation{
		UserID:     req.UserID,
		TicketID:   req.TicketID,
		ReservedAt: time.Now(),
		ExpiresAt:  time.Now().Add(time.Duration(minutes) * time.Minute),
	}
	return models.CreateReservation(rs.DB, reservation)
}

func (rs *ReservationService) GetReservationByID(id int) (*models.Reservation, error) {
	return models.GetReservationByID(rs.DB, id)
}

func (rs *ReservationService) ListReservationsByUser(userID int) ([]*models.Reservation, error) {
	return models.ListReservationsByUser(rs.DB, userID)
}

func (rs *ReservationService) CancelReservation(id int) (*models.Reservation, error) {
	return models.CancelReservation(rs.DB, id)
}

func (rs *ReservationService) PayReservation(req *models.PayReservationRequest, reservationID int) (*models.Payment, *models.Reservation, error) {
	reservation, err := models.GetReservationByID(rs.DB, reservationID)
	if err != nil {
		return nil, nil, err
	}
	ticket, err := models.GetTicketByID(rs.DB, reservation.TicketID)
	if err != nil {
		return nil, nil, err
	}
	updatedReservation, err := models.MarkReservationPaid(rs.DB, reservationID)
	if err != nil {
		return nil, nil, err
	}
	payment := &models.Payment{
		UserID:        reservation.UserID,
		ReservationID: reservationID,
		Amount:        ticket.Price,
		Method:        req.Method,
		Status:        "successful",
	}
	createdPayment, err := models.CreatePayment(rs.DB, payment)
	if err != nil {
		return nil, nil, err
	}
	return createdPayment, updatedReservation, nil
}

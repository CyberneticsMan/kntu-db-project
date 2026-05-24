package models

type CreateReportRequest struct {
	UserID        int    `json:"user_id" binding:"required"`
	TicketID      *int   `json:"ticket_id"`
	ReservationID *int   `json:"reservation_id"`
	Category      string `json:"category" binding:"required"`
	Subject       string `json:"subject" binding:"required"`
	Body          string `json:"body" binding:"required"`
}

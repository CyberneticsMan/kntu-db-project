package models

type CreateReservationRequest struct {
	UserID   int `json:"user_id" binding:"required"`
	TicketID int `json:"ticket_id" binding:"required"`
	Minutes  int `json:"minutes" binding:"omitempty,min=1,max=1440"`
}

type PayReservationRequest struct {
	Method string `json:"method" binding:"required"`
}

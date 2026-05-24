package routes

import (
	"github.com/CyberneticsMan/kntu-db-project/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, tc *controllers.TicketController, uc *controllers.UserController, ac *controllers.AuthController, rc *controllers.ReservationController, rcpt *controllers.ReportController) {
	api := router.Group("/api/v1")
	{
		api.POST("/auth/register", ac.Register)
		api.POST("/auth/login", ac.Login)

		api.GET("/tickets", tc.ListTickets)
		api.POST("/tickets", tc.CreateTicket)
		api.GET("/tickets/:id", tc.GetTicket)
		api.PUT("/tickets/:id", tc.UpdateTicket)
		api.DELETE("/tickets/:id", tc.DeleteTicket)

		api.GET("/users", uc.ListUsers)
		api.GET("/users/:id", uc.GetUser)
		api.POST("/users", uc.CreateUser)
		api.PUT("/users/:id", uc.UpdateUser)
		api.DELETE("/users/:id", uc.DeleteUser)

		api.POST("/reservations", rc.CreateReservation)
		api.GET("/reservations/:id", rc.GetReservation)
		api.GET("/users/:id/reservations", rc.ListUserReservations)
		api.POST("/reservations/:id/cancel", rc.CancelReservation)
		api.POST("/reservations/:id/pay", rc.PayReservation)

		api.POST("/reports", rcpt.CreateReport)
		api.GET("/reports/:id", rcpt.GetReport)
		api.GET("/users/:id/reports", rcpt.ListUserReports)
	}
}

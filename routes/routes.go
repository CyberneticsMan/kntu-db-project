package routes

import (
	"github.com/CyberneticsMan/kntu-db-project/controllers"
	"github.com/CyberneticsMan/kntu-db-project/middleware"
	"github.com/CyberneticsMan/kntu-db-project/models"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, tc *controllers.TicketController, uc *controllers.UserController, ac *controllers.AuthController, rc *controllers.ReservationController, rcpt *controllers.ReportController) {
	api := router.Group("/api/v1")
	{
		api.POST("/auth/register", ac.Register)
		api.POST("/auth/login", ac.Login)

		api.GET("/tickets", tc.ListTickets)
		api.GET("/tickets/:id", tc.GetTicket)

		authenticated := api.Group("")
		authenticated.Use(middleware.Auth())
		{
			authenticated.POST("/users", uc.CreateUser)
			authenticated.GET("/users/:id", uc.GetUser)
			authenticated.PUT("/users/:id", uc.UpdateUser)
			authenticated.DELETE("/users/:id", uc.DeleteUser)

			authenticated.POST("/reservations", rc.CreateReservation)
			authenticated.GET("/reservations/:id", rc.GetReservation)
			authenticated.GET("/users/:id/reservations", rc.ListUserReservations)
			authenticated.POST("/reservations/:id/cancel", rc.CancelReservation)
			authenticated.POST("/reservations/:id/pay", rc.PayReservation)

			authenticated.POST("/reports", rcpt.CreateReport)
			authenticated.GET("/reports/:id", rcpt.GetReport)
			authenticated.GET("/users/:id/reports", rcpt.ListUserReports)
		}

		staffAdmin := authenticated.Group("")
		staffAdmin.Use(middleware.RequireRoles(models.RoleStaff, models.RoleAdmin))
		{
			staffAdmin.GET("/users", uc.ListUsers)
			staffAdmin.POST("/tickets", tc.CreateTicket)
			staffAdmin.PUT("/tickets/:id", tc.UpdateTicket)
			staffAdmin.DELETE("/tickets/:id", tc.DeleteTicket)
		}

		adminOnly := authenticated.Group("")
		adminOnly.Use(middleware.RequireRoles(models.RoleAdmin))
		{
			adminOnly.POST("/admin/users", uc.CreatePrivilegedUser)
		}
	}
}

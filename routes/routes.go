package routes

import (
	"github.com/CyberneticsMan/kntu-db-project/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, tc *controllers.TicketController, uc *controllers.UserController, ac *controllers.AuthController) {
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
	}
}

package routes

import (
	"github.com/CyberneticsMan/kntu-db-project/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, tc *controllers.TicketController, uc *controllers.UserController) {
	api := router.Group("/api/v1")
	{
		api.GET("/tickets/:id", tc.GetTicket)
		// Additional routes can be added here

		// User routes
		api.GET("/users/:id", uc.GetUser)
		api.POST("/users", uc.CreateUser)
		api.PUT("/users", uc.UpdateUser)
		api.DELETE("/users/:id", uc.DeleteUser)

	}
}

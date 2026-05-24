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

		api.GET("/tickets/:id", tc.GetTicket)
		api.GET("/users/:id", uc.GetUser)
		api.POST("/users", uc.CreateUser)
		api.PUT("/users", uc.UpdateUser)
		api.DELETE("/users/:id", uc.DeleteUser)
	}
}

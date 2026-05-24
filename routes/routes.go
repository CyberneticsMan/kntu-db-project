package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/CyberneticsMan/kntu-db-project/controllers"
)

func SetupRoutes(router *gin.Engine, tc *controllers.TicketController) {
	api := router.Group("/api/v1")
	{
		api.GET("/tickets/:id", tc.GetTicket)
	}
}

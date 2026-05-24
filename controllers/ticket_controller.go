package controllers

import (
	"net/http"
	"strconv"

	"github.com/CyberneticsMan/kntu-db-project/services"
	"github.com/gin-gonic/gin"
)

// TicketController holds the service dependency
type TicketController struct {
	Service *services.TicketService
}

// GetTicket handles the GET request for a single ticket
// @Summary Get ticket by ID
// @Description Get a ticket's details by its ID
// @Tags tickets
// @Param id path int true "Ticket ID"
// @Success 200 {object} models.Ticket
// @Failure 400 "Invalid ticket ID"
// @Failure 404 "Ticket not found"
// @Router /tickets/{id} [get]
func (tc *TicketController) GetTicket(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	ticket, err := tc.Service.GetTicketByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

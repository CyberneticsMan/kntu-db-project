package controllers

import (
	"net/http"
	"strconv"

	"github.com/CyberneticsMan/kntu-db-project/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TicketController holds the database dependency
type TicketController struct {
	DB *pgxpool.Pool
}

// GetTicket handles the GET request for a single ticket
func (tc *TicketController) GetTicket(c *gin.Context) {
	// 1. Extract the ID from the URL parameter
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	// 2. Call the Model layer to fetch the data
	ticket, err := models.GetTicketByID(tc.DB, id)
	if err != nil {
		// Log the error internally and return a generic message
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	// 3. Return the data as JSON
	c.JSON(http.StatusOK, ticket)
}

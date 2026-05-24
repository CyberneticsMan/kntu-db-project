package controllers

import (
	"net/http"
	"strconv"

	"github.com/CyberneticsMan/kntu-db-project/models"
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

// ListTickets handles the GET request to list/search tickets
// @Summary List tickets
// @Description List tickets with optional filters
// @Tags tickets
// @Produce json
// @Success 200 {array} models.Ticket
// @Router /tickets [get]
func (tc *TicketController) ListTickets(c *gin.Context) {
	filters := models.TicketFilters{
		SportType: c.Query("sport_type"),
		HomeTeam:  c.Query("home_team"),
		AwayTeam:  c.Query("away_team"),
		City:      c.Query("city"),
		Venue:     c.Query("venue"),
		Category:  c.Query("category"),
		Search:    c.Query("search"),
	}

	if minPrice := c.Query("min_price"); minPrice != "" {
		if value, err := strconv.ParseFloat(minPrice, 64); err == nil {
			filters.MinPrice = &value
		}
	}
	if maxPrice := c.Query("max_price"); maxPrice != "" {
		if value, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			filters.MaxPrice = &value
		}
	}

	tickets, err := tc.Service.ListTickets(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

// CreateTicket handles the POST request to create a ticket
// @Summary Create ticket
// @Description Create a new ticket record
// @Tags tickets
// @Accept json
// @Produce json
// @Param ticket body models.Ticket true "Ticket payload"
// @Success 201 {object} models.Ticket
// @Router /tickets [post]
func (tc *TicketController) CreateTicket(c *gin.Context) {
	var ticket models.Ticket
	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdTicket, err := tc.Service.CreateTicket(&ticket)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdTicket)
}

// UpdateTicket handles the PUT request to update a ticket
// @Summary Update ticket
// @Description Update an existing ticket record
// @Tags tickets
// @Accept json
// @Produce json
// @Param id path int true "Ticket ID"
// @Param ticket body models.Ticket true "Ticket payload"
// @Success 200 {object} models.Ticket
// @Router /tickets/{id} [put]
func (tc *TicketController) UpdateTicket(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	var ticket models.Ticket
	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ticket.ID = id

	if err := tc.Service.UpdateTicket(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// DeleteTicket handles the DELETE request to remove a ticket
// @Summary Delete ticket
// @Description Delete a ticket by its ID
// @Tags tickets
// @Param id path int true "Ticket ID"
// @Success 200 {object} map[string]string
// @Router /tickets/{id} [delete]
func (tc *TicketController) DeleteTicket(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	if err := tc.Service.DeleteTicket(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket deleted successfully"})
}

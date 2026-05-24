package controllers

import (
	"net/http"
	"strconv"

	"github.com/CyberneticsMan/kntu-db-project/models"
	"github.com/CyberneticsMan/kntu-db-project/services"
	"github.com/gin-gonic/gin"
)

type ReservationController struct {
	Service *services.ReservationService
}

// @Security Bearer
func (rc *ReservationController) CreateReservation(c *gin.Context) {
	var req models.CreateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reservation, err := rc.Service.CreateReservation(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, reservation)
}

// @Security Bearer
func (rc *ReservationController) GetReservation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reservation ID"})
		return
	}

	reservation, err := rc.Service.GetReservationByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found"})
		return
	}

	c.JSON(http.StatusOK, reservation)
}

// @Security Bearer
func (rc *ReservationController) ListUserReservations(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	reservations, err := rc.Service.ListReservationsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reservations)
}

// @Security Bearer
func (rc *ReservationController) CancelReservation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reservation ID"})
		return
	}

	reservation, err := rc.Service.CancelReservation(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reservation)
}

// @Security Bearer
func (rc *ReservationController) PayReservation(c *gin.Context) {
	reservationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reservation ID"})
		return
	}

	var req models.PayReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, updatedReservation, err := rc.Service.PayReservation(&req, reservationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payment": payment, "reservation": updatedReservation})
}

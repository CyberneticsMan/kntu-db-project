package controllers

import (
	"net/http"
	"strconv"

	"github.com/CyberneticsMan/kntu-db-project/models"
	"github.com/CyberneticsMan/kntu-db-project/services"
	"github.com/gin-gonic/gin"
)

type ReportController struct {
	Service *services.ReportService
}

// @Security Bearer
func (rc *ReportController) CreateReport(c *gin.Context) {
	var req models.CreateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	report, err := rc.Service.CreateReport(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, report)
}

// @Security Bearer
func (rc *ReportController) GetReport(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}

	report, err := rc.Service.GetReportByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
		return
	}

	c.JSON(http.StatusOK, report)
}

// @Security Bearer
func (rc *ReportController) ListUserReports(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	reports, err := rc.Service.ListReportsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reports)
}

package controllers

import (
	"net/http"
	"strconv"

	"github.com/CyberneticsMan/kntu-db-project/models"
	"github.com/CyberneticsMan/kntu-db-project/services"
	"github.com/gin-gonic/gin"
)

// UserController holds the service dependency
type UserController struct {
	Service *services.UserService
}

// GetUser handles the GET request for a single user
// @Summary Get user by ID
// @Description Get a user's details by their ID
// @Tags users
// @Param id path int true "User ID"
// @Success 200 {object} models.UserResponse
// @Failure 400 "Invalid user ID"
// @Failure 404 "User not found"
// @Router /users/{id} [get]
func (uc *UserController) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := uc.Service.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, models.NewUserResponse(user))
}

// CreateUser handles the POST request to create a new user
// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param userRequest body models.CreateUserRequest true "User creation request"
// @Success 201 {object} models.UserResponse
// @Failure 400 "Invalid request body"
// @Router /users [post]
func (uc *UserController) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default role to customer if not provided
	if req.Role == "" {
		req.Role = models.RoleCustomer
	}

	createdUser, err := uc.Service.CreateUser(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.NewUserResponse(createdUser))
}

// UpdateUser handles the PUT request to update a user
// @Summary Update an existing user
// @Description Update a user's details
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param userRequest body models.UpdateUserRequest true "User update request"
// @Success 200 {object} models.UserResponse
// @Failure 400 "Invalid request or user ID"
// @Router /users/{id} [put]
func (uc *UserController) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ID = id

	updatedUser, err := uc.Service.UpdateUser(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.NewUserResponse(updatedUser))
}

// DeleteUser handles the DELETE request to delete a user
// @Summary Delete a user
// @Description Delete a user by their ID
// @Tags users
// @Param id path int true "User ID"
// @Success 200 "User deleted successfully"
// @Failure 400 "Invalid user ID"
// @Failure 500 "Server error"
// @Router /users/{id} [delete]
func (uc *UserController) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := uc.Service.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

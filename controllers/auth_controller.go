package controllers

import (
	"net/http"

	"github.com/CyberneticsMan/kntu-db-project/models"
	"github.com/CyberneticsMan/kntu-db-project/services"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	AuthService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{AuthService: authService}
}

// Register handles user registration
// @Summary Register a new user
// @Description Create a new user account with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param registerRequest body models.RegisterRequest true "Register request"
// @Success 201 {object} models.AuthResponse
// @Failure 400 "Invalid request or user already exists"
// @Router /auth/register [post]
func (ac *AuthController) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ac.AuthService.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.AuthResponse{
		Token: "",
		User:  models.NewUserResponse(user),
	})
}

// Login handles user login
// @Summary User login
// @Description Authenticate user with email and password, returns JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param loginRequest body models.LoginRequest true "Login request"
// @Success 200 {object} models.AuthResponse "Login successful with JWT token"
// @Failure 400 "Invalid request"
// @Failure 401 "Invalid credentials"
// @Router /auth/login [post]
func (ac *AuthController) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := ac.AuthService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{
		Token: token,
		User:  models.NewUserResponse(user),
	})
}

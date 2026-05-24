package services

import (
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/CyberneticsMan/kntu-db-project/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserService *UserService
	jwtSecret   string
}

func NewAuthService(userService *UserService) *AuthService {
	return &AuthService{
		UserService: userService,
		jwtSecret:   os.Getenv("JWT_SECRET"),
	}
}

func (s *AuthService) Register(req *models.RegisterRequest) (string, *models.User, error) {
	if req == nil {
		return "", nil, errors.New("request is nil")
	}

	if existing, err := s.UserService.GetByEmail(req.Email); err == nil && existing != nil {
		return "", nil, errors.New("email already registered")
	} else if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", nil, err
	}

	user := &models.CreateUserRequest{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  req.Password,
	}

	log.Printf("Registering user: %+v", user) // Debug log

	createdUser, err := s.UserService.CreateUser(user)
	if err != nil {
		return "", nil, err
	}

	token, err := s.generateJWT(createdUser)
	if err != nil {
		return "", nil, err
	}

	return token, createdUser, nil
}

func (s *AuthService) Login(req *models.LoginRequest) (string, *models.User, error) {
	user, err := s.UserService.GetByEmail(req.Email)
	if err != nil {
		return "", nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	token, err := s.generateJWT(user)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *AuthService) generateJWT(user *models.User) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   strconv.Itoa(user.ID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   claims.ExpiresAt.Unix(),
	})

	return token.SignedString([]byte(s.jwtSecret))
}

package services

import (
	"errors"
	"os"
	"time"

	"github.com/CyberneticsMan/kntu-db-project/models"
	"github.com/golang-jwt/jwt/v5"
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

func (s *AuthService) Register(req *models.RegisterRequest) (*models.User, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}

	if existing, _ := s.UserService.GetByEmail(req.Email); existing != nil {
		return nil, errors.New("email already registered")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.CreateUserRequest{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  string(hashed),
	}

	return s.UserService.Create(user)
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
		Subject:   string(rune(user.ID)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   claims.ExpiresAt.Unix(),
	})

	return token.SignedString([]byte(s.jwtSecret))
}

package services

import (
	"log"

	"github.com/CyberneticsMan/kntu-db-project/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	DB *pgxpool.Pool
}

func (us *UserService) GetUserByID(id int) (*models.User, error) {
	return models.GetUserByID(us.DB, id)
}

func (us *UserService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	user := req.ToModel()
	log.Printf("Received user creation request: %+v\n", user) // Debug log
	user.Role = models.RoleCustomer                           // Default role for new users
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)
	log.Printf("Creating user: %+v\n", user) // Debug log
	if err := user.Validate(); err != nil {
		return nil, err
	}
	return models.CreateUser(us.DB, user)
}

func (us *UserService) UpdateUser(req *models.UpdateUserRequest) (*models.User, error) {
	user := req.ToModel()
	if err := user.Validate(); err != nil {
		return nil, err
	}
	if err := models.UpdateUser(us.DB, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (us *UserService) DeleteUser(id int) error {
	return models.DeleteUser(us.DB, id)
}

func (us *UserService) GetByEmail(email string) (*models.User, error) {
	return models.GetUserByEmail(us.DB, email)
}

func (us *UserService) Update(req *models.UpdateUserRequest) (*models.User, error) {
	user := req.ToModel()
	if err := user.Validate(); err != nil {
		return nil, err
	}
	if err := models.UpdateUser(us.DB, user); err != nil {
		return nil, err
	}
	return user, nil
}

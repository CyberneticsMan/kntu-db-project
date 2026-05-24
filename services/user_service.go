package services

import (
	"errors"
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
	return us.createUserWithRole(req.ToModel(), models.RoleCustomer)
}

func (us *UserService) CreatePrivilegedUser(req *models.CreatePrivilegedUserRequest) (*models.User, error) {
	if req.Role != models.RoleAdmin && req.Role != models.RoleStaff {
		return nil, errors.New("invalid role for privileged user creation")
	}
	return us.createUserWithRole(req.ToModel(), req.Role)
}

func (us *UserService) createUserWithRole(user *models.User, role models.UserRole) (*models.User, error) {
	log.Printf("Received user creation request: %+v\n", user) // Debug log
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Role = role
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

func (us *UserService) ListUsers() ([]*models.User, error) {
	return models.ListUsers(us.DB)
}

func (us *UserService) PromoteToAdmin(email, phone string) (*models.User, error) {
	return models.SetUserRoleByIdentifier(us.DB, models.RoleAdmin, email, phone)
}

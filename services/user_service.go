package services

import (
	"github.com/CyberneticsMan/kntu-db-project/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService struct {
	DB *pgxpool.Pool
}

func (us *UserService) GetUserByID(id int) (*models.User, error) {
	return models.GetUserByID(us.DB, id)
}

func (us *UserService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	user := req.ToModel()
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

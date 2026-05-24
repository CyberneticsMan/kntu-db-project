package models

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Define a custom type for UserRole
type UserRole string

// Use constants for allowed roles to ensure type safety
const (
	RoleAdmin    UserRole = "admin"
	RoleCustomer UserRole = "customer"
	RoleStaff    UserRole = "staff"
)

type User struct {
	ID       int      `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
	Role     UserRole `json:"role"` // Use the custom type here
}

// IsValid checks if the provided role is one of the allowed constants
func (r UserRole) IsValid() bool {
	switch r {
	case RoleAdmin, RoleCustomer, RoleStaff:
		return true
	}
	return false
}

func (u *User) Validate() error {
	if u.Username == "" {
		return fmt.Errorf("username is required")
	}
	if u.Email == "" {
		return fmt.Errorf("email is required")
	}
	if u.Phone == "" {
		return fmt.Errorf("phone is required")
	}
	if !u.Role.IsValid() {
		return fmt.Errorf("invalid role: %s", u.Role)
	}
	return nil
}

func GetUserByID(db *pgxpool.Pool, id int) (*User, error) {
	var u User
	query := `SELECT id, username, email, phone, role FROM users WHERE id = $1`

	err := db.QueryRow(context.Background(), query, id).Scan(&u.ID, &u.Username, &u.Email, &u.Phone, &u.Role)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func CreateUser(db *pgxpool.Pool, user *User) (*User, error) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	query := `INSERT INTO users (username, email, phone, role) VALUES ($1, $2, $3, $4) RETURNING id`
	err := db.QueryRow(context.Background(), query, user.Username, user.Email, user.Phone, user.Role).Scan(&user.ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func UpdateUser(db *pgxpool.Pool, user *User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	query := `UPDATE users SET username = $1, email = $2, phone = $3, role = $4 WHERE id = $5`
	_, err := db.Exec(context.Background(), query, user.Username, user.Email, user.Phone, user.Role, user.ID)
	return err
}

func DeleteUser(db *pgxpool.Pool, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := db.Exec(context.Background(), query, id)
	return err
}

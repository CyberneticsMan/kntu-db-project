package models_test

import (
	"context"
	"testing"

	"github.com/CyberneticsMan/kntu-db-project/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// Example using a real test database URL from env (recommended)
// or a local test DB like "postgres://user:pass@localhost:5432/kntu_test"
func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	dsn := "postgres://postgres:postgres@localhost:5432/kntu_test?sslmode=disable"

	pool, err := pgxpool.New(context.Background(), dsn)
	require.NoError(t, err)

	// Optionally truncate tables here

	return pool
}

func TestGetUserByID_Found(t *testing.T) {
	db := setupTestDB(t)	
	defer db.Close()

	// Arrange: insert a user row for test
	_, err := db.Exec(context.Background(),
		`INSERT INTO users (id, username, email, phone, role)
		 VALUES ($1, $2, $3, $4, $5)`,
		1, "testuser", "test@example.com", "09120000000", models.RoleCustomer,
	)
	require.NoError(t, err)

	// Act
	u, err := models.GetUserByID(db, 1)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, u)
	require.Equal(t, 1, u.ID)
	require.Equal(t, "testuser", u.Username)
	require.Equal(t, models.RoleCustomer, u.Role)
}

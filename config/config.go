package config

import (
	"fmt"
	"os"
)

// Config holds application configuration values loaded from the environment.
type Config struct {
	DatabaseURL string
}

// LoadConfig reads configuration from environment variables and returns a Config.
// It currently requires DATABASE_URL to be set.
func LoadConfig() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}
	return &Config{DatabaseURL: dbURL}, nil
}

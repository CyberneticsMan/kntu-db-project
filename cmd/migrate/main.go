package migrate

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/CyberneticsMan/kntu-db-project/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := runMigrations(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("Migrations completed successfully!")
}

func runMigrations(db *pgxpool.Pool) error {
	// Run the SQL files in the correct order from the migrations directory
	// Read Migration files from the "migrations" directory using os.ReadDir and execute them in order

	migrationFiles, err := os.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, file := range migrationFiles {
		if file.IsDir() {
			continue
		}
		if err := executeSQLFile(db, "migrations/"+file.Name()); err != nil {
			return fmt.Errorf("failed to execute %s: %w", file.Name(), err)
		}
	}

	return nil
}

func executeSQLFile(db *pgxpool.Pool, filePath string) error {
	sqlBytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("could not read SQL file: %w", err)
	}

	sql := string(sqlBytes)
	_, err = db.Exec(context.Background(), sql)
	if err != nil {
		return fmt.Errorf("failed to execute SQL: %w", err)
	}

	return nil
}

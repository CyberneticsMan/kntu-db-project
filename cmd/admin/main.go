package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/CyberneticsMan/kntu-db-project/config"
	"github.com/CyberneticsMan/kntu-db-project/services"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	email := flag.String("email", "", "User email")
	phone := flag.String("phone", "", "User phone number")
	flag.Parse()

	if *email == "" && *phone == "" {
		log.Fatal("either --email or --phone must be provided")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	userService := &services.UserService{DB: db}
	updatedUser, err := userService.PromoteToAdmin(*email, *phone)
	if err != nil {
		log.Fatalf("Failed to promote user: %v", err)
	}

	log.Printf("User promoted to admin: id=%d email=%s phone=%s role=%s", updatedUser.ID, updatedUser.Email, updatedUser.Phone, updatedUser.Role)
	fmt.Fprintln(os.Stdout, "done")
}

package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/CyberneticsMan/kntu-db-project/controllers"
	"github.com/CyberneticsMan/kntu-db-project/routes"
)

func main() {
	// 1. Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// 2. Connect to database using pgxpool
	dbURL := os.Getenv("DATABASE_URL")
	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	// 3. Initialize Gin router
	router := gin.Default()

	// 4. Dependency injection for controllers
	ticketController := &controllers.TicketController{DB: dbPool}

	// 5. Define routes
	routes.SetupRoutes(router, ticketController)

	// 6. Run server
	router.Run(":8080")
}

package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/CyberneticsMan/kntu-db-project/controllers"
	_ "github.com/CyberneticsMan/kntu-db-project/docs" // This line is important for swagger initialization
	"github.com/CyberneticsMan/kntu-db-project/routes"
	"github.com/CyberneticsMan/kntu-db-project/services"
)

// @title           KNTU DB Project API
// @version         1.0
// @description     A Go Gin API for managing tickets and users
// @basePath         /api/v1
// @schemes http https
// @securityDefinitions.apikey  Bearer
// @in                          header
// @name                        Authorization
// @description                 Enter the token with the `Bearer: ` prefix, e.g. "Bearer abcde12345"

func main() {
	// 1. Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	log.SetOutput(os.Stdout)

	// 2. Connect to database using pgxpool
	dbURL := os.Getenv("DATABASE_URL")
	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	// 3. Initialize Gin router
	router := gin.Default()

	// 4. Dependency injection for services and controllers
	ticketService := &services.TicketService{DB: dbPool}
	userService := &services.UserService{DB: dbPool}
	reservationService := &services.ReservationService{DB: dbPool}
	reportService := &services.ReportService{DB: dbPool}

	ticketController := &controllers.TicketController{Service: ticketService}
	userController := &controllers.UserController{Service: userService}
	reservationController := &controllers.ReservationController{Service: reservationService}
	reportController := &controllers.ReportController{Service: reportService}

	authService := services.NewAuthService(userService)
	authController := controllers.NewAuthController(authService)

	// 5. Setup Swagger docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 6. Define routes
	routes.SetupRoutes(router, ticketController, userController, authController, reservationController, reportController)

	// 7. Run server
	router.Run(":8080")
}

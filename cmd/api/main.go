package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/CyberneticsMan/kntu-db-project/controllers" // این مسیر را با نام ماژول خود جایگزین کنید
	"github.com/CyberneticsMan/kntu-db-project/routes"
)

func main() {
	// ۱. بارگذاری متغیرهای محیطی از فایل .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// ۲. اتصال به دیتابیس (pgxpool)
	dbURL := os.Getenv("DATABASE_URL")
	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	// ۳. مقداردهی اولیه Gin
	router := gin.Default()

	// ۴. تزریق وابستگی (Dependency Injection) به کنترلرها
	ticketController := &controllers.TicketController{DB: dbPool}

	// ۵. تعریف روت‌ها
	routes.SetupRoutes(router, ticketController)

	// ۶. اجرای سرور
	router.Run(":8080")
}

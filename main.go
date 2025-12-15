package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"prestasi-mahasiswa-api/database"
	"prestasi-mahasiswa-api/routes"
    

    _ "prestasi-mahasiswa-api/docs" 

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

// @title           Sistem Pelaporan Prestasi Mahasiswa API
// @version         1.0
// @description     API Server untuk manajemen prestasi mahasiswa, dosen wali, dan admin.
// @termsOfService  http://swagger.io/terms/

// @contact.name    Tim Pengembang
// @contact.email   support@unair.ac.id

// @host            localhost:3000
// @BasePath        /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// 1. Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables.")
	}

	// 2. Inisialisasi Database
	pgPool, mongoClient, err := database.ConnectDatabases()
	if err != nil {
		log.Fatalf("Failed to initialize databases: %v", err)
	}
	defer pgPool.Close()
	
	// 3. Init Fiber App
	app := fiber.New(fiber.Config{
		AppName: "Sistem Pelaporan Prestasi Mahasiswa API",
	})

	app.Use(cors.New())
	app.Use(logger.New())

	// 4. Setup Routes
	routes.SetupRoutes(app, pgPool, mongoClient)

	// 5. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Graceful Shutdown Setup
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		log.Println("Gracefully shutting down...")

		if err := app.Shutdown(); err != nil {
			log.Printf("Error during Fiber shutdown: %v", err)
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting MongoDB: %v", err)
		}
	}()

	log.Printf("🚀 Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}
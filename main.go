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

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

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

	// Middleware Global
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

		// Block sampai sinyal diterima (Perbaikan S1005)
		<-c
		log.Println("Gracefully shutting down...")

		// Shutdown Fiber
		if err := app.Shutdown(); err != nil {
			log.Printf("Error during Fiber shutdown: %v", err)
		}
		
		// Disconnect MongoDB properly
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
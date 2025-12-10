package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"prestasi-mahasiswa-api/database"
	"prestasi-mahasiswa-api/routes"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables (Opsional, tapi sangat disarankan)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables.")
	}

	// Inisialisasi Database
	pgConn, mongoClient, err := database.ConnectDatabases()
	if err != nil {
		log.Fatalf("Failed to initialize databases: %v", err)
	}
	defer pgConn.Close()
	defer func() {
		if err = mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting MongoDB: %v", err)
		}
	}()

	// Inisialisasi Router
	r := mux.NewRouter()

	// Daftarkan Semua Routes
	routes.SetupRoutes(r, pgConn, mongoClient)

	// Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	addr := fmt.Sprintf(":%s", port)

	srv := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("🚀 Server starting on port %s", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}
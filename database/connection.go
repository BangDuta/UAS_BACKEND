package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectDatabases menginisialisasi koneksi ke PostgreSQL dan MongoDB.
func ConnectDatabases() (*pgxpool.Pool, *mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// --- 1. PostgreSQL (Relational Data & RBAC) ---
	pgURL := os.Getenv("POSTGRES_URL")
	if pgURL == "" {
		log.Println("POSTGRES_URL not set. Using default.")
		pgURL = "postgres://user:password@localhost:5432/prestasi_db"
	}

	pgConn, err := pgxpool.New(ctx, pgURL)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to connect to PostgreSQL: %w", err)
	}

	if err = pgConn.Ping(ctx); err != nil {
		return nil, nil, fmt.Errorf("PostgreSQL ping failed: %w", err)
	}
	log.Println("✅ Connected to PostgreSQL successfully.")

	// --- 2. MongoDB (Dynamic Achievement Data) ---
	mongoURL := os.Getenv("MONGO_URL")
	if mongoURL == "" {
		log.Println("MONGO_URL not set. Using default.")
		mongoURL = "mongodb://localhost:27017"
	}
	
	clientOptions := options.Client().ApplyURI(mongoURL)
	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return pgConn, nil, fmt.Errorf("unable to connect to MongoDB: %w", err)
	}

	if err = mongoClient.Ping(ctx, nil); err != nil {
		return pgConn, nil, fmt.Errorf("MongoDB ping failed: %w", err)
	}
	log.Println("✅ Connected to MongoDB successfully.")

	return pgConn, mongoClient, nil
}
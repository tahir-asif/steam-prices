// Package database provides PostgreSQL connection handling and data access
// functions for managing games and their price history.
package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// ConnectDB establishes a connection to PostgreSQL using the DATABASE_URL
// environment variable. It falls back to a local Docker connection string
// if the variable is not set.
func ConnectDB() *sql.DB {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://steam_user:steam_pass@localhost:5432/steam_prices?sslmode=disable"
		log.Println("DATABASE_URL not set, using fallback local connection string")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	// Verify the connection is alive
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to database successfully")
	return db
}


package main

import (
    "database/sql"
    "log"
    "os"
    _ "github.com/lib/pq"
)

// ConnectDB establishes a connection to PostgreSQL and returns the database handle.
func ConnectDB() *sql.DB {
    // Get connection string from environment variable, or use default local string
    connStr := os.Getenv("DATABASE_URL")
    if connStr == "" {
        // Fallback for local development
        connStr = "postgres://steam_user:steam_pass@localhost:5432/steam_prices?sslmode=disable"
        log.Println("Using default local database connection")
    }

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatalf("Failed to open database connection: %v", err)
    }

    // Verify the connection is actually alive
    if err := db.Ping(); err != nil {
        log.Fatalf("Failed to ping database: %v", err)
    }

    log.Println("Connected to database successfully")
    return db
}

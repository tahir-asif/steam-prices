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

// InsertGame adds a new game or updates the name if the steam_app_id already exists.
// It returns the internal game ID.
func InsertGame(db *sql.DB, steamAppID int, name string) (int, error) {
	var id int
	err := db.QueryRow(`
        INSERT INTO games (steam_app_id, name)
        VALUES ($1, $2)
        ON CONFLICT (steam_app_id) DO UPDATE SET name = EXCLUDED.name
        RETURNING id
    `, steamAppID, name).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetGameIDBySteamAppID retrieves the internal ID for a given Steam App ID.
// Returns 0 and sql.ErrNoRows if not found.
func GetGameIDBySteamAppID(db *sql.DB, steamAppID int) (int, error) {
	var id int
	err := db.QueryRow(`SELECT id FROM games WHERE steam_app_id = $1`, steamAppID).Scan(&id)
	return id, err
}

// GetLastRecordedPrice returns the most recent price and currency for a game.
// If no price history exists, it returns 0, "", and sql.ErrNoRows.
func GetLastRecordedPrice(db *sql.DB, gameID int) (price int, currency string, err error) {
	err = db.QueryRow(`
        SELECT price, currency
        FROM price_history
        WHERE game_id = $1
        ORDER BY recorded_at DESC
        LIMIT 1
    `, gameID).Scan(&price, &currency)
	return
}

// InsertPriceRecord adds a new price history entry.
func InsertPriceRecord(db *sql.DB, gameID, price int, currency string) error {
	_, err := db.Exec(`
        INSERT INTO price_history (game_id, price, currency)
        VALUES ($1, $2, $3)
    `, gameID, price, currency)
	return err
}

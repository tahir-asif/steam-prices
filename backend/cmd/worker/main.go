// Package main is the entry point for the background worker that periodically
// fetches price data from Steam and updates the database.
package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/tahir-asif/steam-prices/internal/config"
	"github.com/tahir-asif/steam-prices/internal/database"
	"github.com/tahir-asif/steam-prices/internal/steam"
)

func main() {
	// Load environment variables (DATABASE_URL, etc.)
	config.LoadEnv()

	// Connect to the database
	db := database.ConnectDB()
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	// Create a Steam API client
	steamClient := steam.NewClient()

	// List of Steam App IDs to track (popular games)
	trackedAppIDs := []int{
		730,    // Counter-Strike 2
		570,    // Dota 2
		440,    // Team Fortress 2
		292030, // The Witcher 3: Wild Hunt
		271590, // Grand Theft Auto V
	}

	log.Printf("Worker started. Tracking %d games.", len(trackedAppIDs))

	// Run once immediately on startup
	runPriceCheck(db, steamClient, trackedAppIDs)

	// Then run every hour
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		runPriceCheck(db, steamClient, trackedAppIDs)
	}
}

// runPriceCheck fetches current prices for all tracked games and records
// changes in the database.
func runPriceCheck(db *sql.DB, client *steam.Client, appIDs []int) {
	log.Println("Starting price check cycle...")

	for _, appID := range appIDs {
		log.Printf("Fetching game %d...", appID)

		// Fetch game details from Steam
		game, err := client.FetchGameDetails(appID)
		if err != nil {
			log.Printf("  Error fetching game %d: %v", appID, err)
			continue
		}

		// Insert or update game in database, get internal ID
		gameID, err := database.InsertGame(db, appID, game.Name)
		if err != nil {
			log.Printf("  Error inserting game %d: %v", appID, err)
			continue
		}

		// Get the last recorded price for this game
		lastPrice, lastCurrency, err := database.GetLastRecordedPrice(db, gameID)
		priceChanged := false
		if err != nil {
			// No previous price (sql.ErrNoRows), so this is a new record
			priceChanged = true
		} else {
			// Compare with last known price
			priceChanged = (lastPrice != game.Price) || (lastCurrency != game.Currency)
		}

		// Insert a new price record only if it changed or is new
		if priceChanged {
			err = database.InsertPriceRecord(db, gameID, game.Price, game.Currency)
			if err != nil {
				log.Printf("  Error inserting price record for %d: %v", appID, err)
				continue
			}
			if lastPrice == 0 && lastCurrency == "" {
				log.Printf("  Initial price recorded for %s: %d %s",
					game.Name, game.Price, game.Currency)
			} else {
				log.Printf("  Price changed for %s: %d %s (was %d %s)",
					game.Name, game.Price, game.Currency, lastPrice, lastCurrency)
			}
		} else {
			log.Printf("  No price change for %s (%d %s)",
				game.Name, game.Price, game.Currency)
		}
	}

	log.Println("Price check cycle complete.")
}

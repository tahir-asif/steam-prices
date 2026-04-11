package main

import (
	"database/sql"
	"log"

	"github.com/tahir-asif/steam-prices/internal/config"
	"github.com/tahir-asif/steam-prices/internal/database"
	"github.com/tahir-asif/steam-prices/internal/steam"
)

func main() {
	config.LoadEnv()

	db := database.ConnectDB()
	defer db.Close()

	steamClient := steam.NewClient()

	log.Println("Starting scheduled price check...")
	runPriceCheck(db, steamClient)
	log.Println("Scheduled price check complete.")
}

// getTrackedAppIDs retrieves all Steam App IDs currently stored in the games table.
// If the query fails, it falls back to a default list of popular games.
func getTrackedAppIDs(db *sql.DB) []int {
	rows, err := db.Query("SELECT steam_app_id FROM games")
	if err != nil {
		log.Printf("Error fetching tracked games: %v", err)
		// Fallback to a few popular games if the database query fails
		return []int{730, 570, 440, 292030, 271590}
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}

// runPriceCheck orchestrates the price checking process for all tracked games.
func runPriceCheck(db *sql.DB, client *steam.Client) {
	appIDs := getTrackedAppIDs(db)
	if len(appIDs) == 0 {
		log.Println("No tracked games found.")
		return
	}

	log.Printf("Checking prices for %d games...", len(appIDs))

	for _, appID := range appIDs {
		processGame(db, client, appID)
	}

	log.Println("Price check cycle finished.")
}

// processGame fetches the current price for a single Steam App ID and updates the database.
func processGame(db *sql.DB, client *steam.Client, appID int) {
	game, err := client.FetchGameDetails(appID)
	if err != nil {
		log.Printf("Error fetching game %d: %v", appID, err)
		return
	}

	gameID, err := database.InsertGame(db, appID, game.Name)
	if err != nil {
		log.Printf("Error inserting/updating game %d: %v", appID, err)
		return
	}

	if hasPriceChanged(db, gameID, game.Price, game.Currency) {
		if err := database.InsertPriceRecord(db, gameID, game.Price, game.Currency); err != nil {
			log.Printf("Error inserting price for game %d: %v", appID, err)
		} else {
			log.Printf("Price updated for %s: %d %s", game.Name, game.Price, game.Currency)
		}
	}
}

// hasPriceChanged checks whether the current price differs from the last recorded price.
func hasPriceChanged(db *sql.DB, gameID, newPrice int, newCurrency string) bool {
	lastPrice, lastCurrency, err := database.GetLastRecordedPrice(db, gameID)
	if err != nil {
		// No previous price found, so treat as a change
		return true
	}
	return lastPrice != newPrice || lastCurrency != newCurrency
}

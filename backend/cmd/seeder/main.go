package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/tahir-asif/steam-prices/internal/config"
	"github.com/tahir-asif/steam-prices/internal/database"
	"github.com/tahir-asif/steam-prices/internal/steam"
)

// SteamSpyTop100Response is the structure returned by the Steam Spy API.
type SteamSpyTop100Response map[string]struct {
	AppID int    `json:"appid"`
	Name  string `json:"name"`
}

func main() {
	config.LoadEnv()

	db := database.ConnectDB()
	defer db.Close()

	steamClient := steam.NewClient()

	log.Println("Fetching top 100 games from Steam Spy...")
	appIDs, err := fetchTop100AppIDs()
	if err != nil {
		log.Fatalf("Failed to fetch top games: %v", err)
	}

	log.Printf("Found %d games. Starting to seed...", len(appIDs))

	successCount := 0
	for i, appID := range appIDs {
		log.Printf("[%d/%d] Processing App ID %d...", i+1, len(appIDs), appID)

		game, err := steamClient.FetchGameDetails(appID)
		if err != nil {
			log.Printf("  Error fetching details for %d: %v", appID, err)
			continue
		}

		gameID, err := database.InsertGame(db, appID, game.Name)
		if err != nil {
			log.Printf("  Error inserting game %d: %v", appID, err)
			continue
		}

		// Only insert a price record if the game has a price (free games have price 0)
		if err := database.InsertPriceRecord(db, gameID, game.Price, game.Currency); err != nil {
			log.Printf("  Error inserting price for %d: %v", appID, err)
		} else {
			log.Printf("  Successfully added %s (%d %s)", game.Name, game.Price, game.Currency)
		}

		successCount++

		// Be nice to Steam's API – wait 1.5 seconds between requests
		time.Sleep(1500 * time.Millisecond)
	}

	log.Printf("Seeding complete. Successfully added %d out of %d games.", successCount, len(appIDs))
}

// fetchTop100AppIDs calls the Steam Spy API and returns a slice of App IDs.
func fetchTop100AppIDs() ([]int, error) {
	url := "https://steamspy.com/api.php?request=top100in2weeks"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data SteamSpyTop100Response
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("json decode failed: %w", err)
	}

	// Extract App IDs from the map keys
	// Note: The response uses the app ID as the key (as a string)
	var appIDs []int
	for key := range data {
		var appID int
		if _, err := fmt.Sscanf(key, "%d", &appID); err == nil {
			appIDs = append(appIDs, appID)
		}
	}

	return appIDs, nil
}

// Package main is the entry point for the Steam Price Tracker REST API server.
// It initializes the database connection, sets up middleware, and defines routes.
package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/tahir-asif/steam-prices/internal/config"
	"github.com/tahir-asif/steam-prices/internal/database"
	"github.com/tahir-asif/steam-prices/internal/handlers"
    "github.com/tahir-asif/steam-prices/internal/steam"
)


func main() {
	// Loads enviroment variables only if they aren't loaded by Render, Neon, etc.
	config.LoadEnv()

	// Connect to database
	db := database.ConnectDB()
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	go runWorker(db)

	// Create router
	router := SetupRouter(db)

	// Runs the server
	log.Println("Server is running on http://localhost:3000")
	if err := router.Run(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// SetupRouter configures the Gin engine with middleware and routes.
// It is exported so it can be used in tests.
func SetupRouter(db *sql.DB) *gin.Engine {
	router := gin.Default()

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{
			"http://localhost:5173",
			"https://steam-prices.vercel.app/",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	h := &handlers.Handler{DB: db}

	router.GET("/api/health", h.Health)
	router.GET("/api/db-test", h.DBTest)
	router.GET("/api/search", h.Search)
	router.GET("/api/games/:appid/history", h.PriceHistory)

	return router
}

// Worker logic had to be moved here because Render CRON jobs cost money
// runWorker is a background goroutine that periodically fetches prices from Steam.
func runWorker(db *sql.DB) {
    steamClient := steam.NewClient()

    // Run immediately on startup
    doPriceCheck(db, steamClient)

    // Then run every hour
    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()

    for range ticker.C {
        doPriceCheck(db, steamClient)
    }
}

// doPriceCheck performs the actual price fetching and database updates.
func doPriceCheck(db *sql.DB, client *steam.Client) {
    log.Println("Worker: Starting price check cycle...")

    // Get all tracked app IDs from the database
    rows, err := db.Query("SELECT steam_app_id FROM games")
    if err != nil {
        log.Printf("Worker: Error fetching tracked games: %v", err)
        return
    }
    defer rows.Close()

    var appIDs []int
    for rows.Next() {
        var id int
        if err := rows.Scan(&id); err == nil {
            appIDs = append(appIDs, id)
        }
    }

    if len(appIDs) == 0 {
        log.Println("Worker: No tracked games found in database")
        return
    }

    log.Printf("Worker: Checking %d games", len(appIDs))

    for _, appID := range appIDs {
        game, err := client.FetchGameDetails(appID)
        if err != nil {
            log.Printf("Worker: Error fetching game %d: %v", appID, err)
            continue
        }

        // Insert/update game
        gameID, err := database.InsertGame(db, appID, game.Name)
        if err != nil {
            log.Printf("Worker: Error inserting game %d: %v", appID, err)
            continue
        }

        // Get last price
        lastPrice, lastCurrency, err := database.GetLastRecordedPrice(db, gameID)
        priceChanged := false
        if err != nil {
            priceChanged = true // no previous record
        } else {
            priceChanged = (lastPrice != game.Price) || (lastCurrency != game.Currency)
        }

        if priceChanged {
            err = database.InsertPriceRecord(db, gameID, game.Price, game.Currency)
            if err != nil {
                log.Printf("Worker: Error inserting price for game %d: %v", appID, err)
            } else {
                log.Printf("Worker: Price changed for %s: %d %s", game.Name, game.Price, game.Currency)
            }
        }
    }

    log.Println("Worker: Price check cycle complete")
}

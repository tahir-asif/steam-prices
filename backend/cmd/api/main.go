// Package main is the entry point for the Steam Price Tracker REST API server.
// It initializes the database connection, sets up middleware, and defines routes.
package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/tahir-asif/steam-prices/internal/config"
	"github.com/tahir-asif/steam-prices/internal/database"
)

var db *sql.DB

func main() {
	// Loads enviroment variables only if they aren't loaded by Render, Neon, etc.
	config.LoadEnv()

	// Connect to database
	db = database.ConnectDB()
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	router := gin.Default()

	// CORS middleware configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// Database test endpoint
	router.GET("/api/db-test", func(c *gin.Context) {
		var now string
		err := db.QueryRow("SELECT NOW()").Scan(&now)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Database query failed",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"db_time": now,
			"message": "Database connection is working!",
		})
	})

	// Search endpoint – proxies Steam's store search
	router.GET("/api/search", func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing search query 'q'"})
			return
		}

		// Call Steam's public search API
		steamURL := fmt.Sprintf("https://store.steampowered.com/api/storesearch/?term=%s&cc=ca&l=english", query)
		resp, err := http.Get(steamURL)
		if err != nil {
			log.Printf("Steam search error: %v", err)
			c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to reach Steam API"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusBadGateway, gin.H{"error": "Steam API returned an error"})
			return
		}

		// Parse the JSON response
		var result struct {
			Items []struct {
				AppID int    `json:"id"`
				Name  string `json:"name"`
				Icon  string `json:"tiny_image"`
			} `json:"items"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			log.Printf("Failed to parse Steam search response: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid response from Steam"})
			return
		}

		// Transform to a cleaner format for the frontend
		type searchItem struct {
			AppID int    `json:"appid"`
			Name  string `json:"name"`
			Icon  string `json:"icon"`
		}
		items := make([]searchItem, 0, len(result.Items))
		for _, item := range result.Items {
			items = append(items, searchItem{
				AppID: item.AppID,
				Name:  item.Name,
				Icon:  item.Icon,
			})
		}

		c.JSON(http.StatusOK, items)
	})

	// Price history endpoint – returns stored price points
	router.GET("/api/games/:appid/history", func(c *gin.Context) {
		appIDStr := c.Param("appid")
		appID, err := strconv.Atoi(appIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid app ID"})
			return
		}

		// First, get the internal game ID from the database
		var gameID int
		err = db.QueryRow("SELECT id FROM games WHERE steam_app_id = $1", appID).Scan(&gameID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Game not found in database"})
				return
			}
			log.Printf("Database error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
			return
		}

		// Query price history
		rows, err := db.Query(`
			SELECT price, currency, recorded_at
			FROM price_history
			WHERE game_id = $1
			ORDER BY recorded_at ASC
		`, gameID)
		if err != nil {
			log.Printf("Database error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch price history"})
			return
		}
		defer rows.Close()

		type pricePoint struct {
			Price      int       `json:"price"`
			Currency   string    `json:"currency"`
			RecordedAt time.Time `json:"recorded_at"`
		}
		var history []pricePoint

		for rows.Next() {
			var p pricePoint
			if err := rows.Scan(&p.Price, &p.Currency, &p.RecordedAt); err != nil {
				log.Printf("Row scan error: %v", err)
				continue
			}
			history = append(history, p)
		}

		c.JSON(http.StatusOK, history)
	})

	// health check
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// runs the server
	log.Println("Server is running on http://localhost:3000")
	if err := router.Run(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

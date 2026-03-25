package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/tahir-asif/steam-prices/internal/database"
	"github.com/tahir-asif/steam-prices/internal/steam"
)

func (h *Handler) PriceHistory(c *gin.Context) {
	appIDStr := c.Param("appid")
	appID, err := strconv.Atoi(appIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid app ID"})
		return
	}

	var gameID int
	err = h.DB.QueryRow("SELECT id FROM games WHERE steam_app_id = $1", appID).Scan(&gameID)

	// Game not in database – fetch it from Steam and start tracking
	if err == sql.ErrNoRows {
		log.Printf("Game %d not in DB, fetching from Steam...", appID)

		steamClient := steam.NewClient()
		game, fetchErr := steamClient.FetchGameDetails(appID)
		if fetchErr != nil {
			log.Printf("Failed to fetch game %d from Steam: %v", appID, fetchErr)
			c.JSON(http.StatusNotFound, gin.H{"error": "Game not found on Steam"})
			return
		}

		// Insert the game
		gameID, err = database.InsertGame(h.DB, appID, game.Name)
		if err != nil {
			log.Printf("Failed to insert game %d: %v", appID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save game"})
			return
		}

		// Record its current price
		if err := database.InsertPriceRecord(h.DB, gameID, game.Price, game.Currency); err != nil {
			log.Printf("Failed to insert initial price for game %d: %v", appID, err)
			// Non‑fatal – continue to return the data we have
		}

		log.Printf("Started tracking new game: %s (App ID %d)", game.Name, appID)
	} else if err != nil {
		log.Printf("Database error looking up game %d: %v", appID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
		return
	}

	// Fetch price history (now guaranteed to exist)
	rows, err := h.DB.Query(`
		SELECT price, currency, recorded_at
		FROM price_history
		WHERE game_id = $1
		ORDER BY recorded_at ASC
	`, gameID)
	if err != nil {
		log.Printf("Database error fetching history for game %d: %v", gameID, err)
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
}

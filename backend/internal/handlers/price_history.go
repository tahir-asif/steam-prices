package handlers

import (
    "database/sql"
    "log"
    "net/http"
    "strconv"
    "time"
    "github.com/gin-gonic/gin"
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
    if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusNotFound, gin.H{"error": "Game not found in database"})
            return
        }
        log.Printf("Database error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
        return
    }

    rows, err := h.DB.Query(`
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
}

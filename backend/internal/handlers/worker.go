package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tahir-asif/steam-prices/internal/steam"
	"github.com/tahir-asif/steam-prices/internal/worker"
)

func (h *Handler) RunWorker(c *gin.Context) {
	expectedToken := os.Getenv("WORKER_SECRET")
	if expectedToken == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "worker secret not configured"})
		return
	}

	token := c.GetHeader("Authorization")
	if token != "Bearer "+expectedToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Run the price check in a goroutine so the API can respond quickly
	go worker.RunPriceCheck(h.DB, steam.NewClient())

	c.JSON(http.StatusOK, gin.H{"message": "price check started"})
}

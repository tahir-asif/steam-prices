// Package main is the entry point for the Steam Price Tracker REST API server.
// It initializes the database connection, sets up middleware, and defines routes.
package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/tahir-asif/steam-prices/internal/config"
	"github.com/tahir-asif/steam-prices/internal/database"
	"github.com/tahir-asif/steam-prices/internal/handlers"
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

	// Create router
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

	// listen for API requests from front end
	h := &handlers.Handler{DB: db}
	router.GET("/api/health", h.Health)
    router.GET("/api/db-test", h.DBTest)
    router.GET("/api/search", h.Search)
    router.GET("/api/games/:appid/history", h.PriceHistory)

	// runs the server
	log.Println("Server is running on http://localhost:3000")
	if err := router.Run(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

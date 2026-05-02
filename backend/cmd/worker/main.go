package main

import (
	"log"

	"github.com/tahir-asif/steam-prices/internal/config"
	"github.com/tahir-asif/steam-prices/internal/database"
	"github.com/tahir-asif/steam-prices/internal/steam"
	"github.com/tahir-asif/steam-prices/internal/worker"
)

func main() {
	config.LoadEnv()

	db := database.ConnectDB()
	defer db.Close()

	steamClient := steam.NewClient()

	log.Println("Starting scheduled price check...")
	worker.RunPriceCheck(db, steamClient)
	log.Println("Scheduled price check complete.")
}

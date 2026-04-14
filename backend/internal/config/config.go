// Package config handles loading environment variables from .env files
// for local development while remaining compatible with production platforms.
package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from .env file for local development.
// In production, environment variables are already set by the hosting platform
// and so LoadEnv does nothing.
func LoadEnv() {
	// If DATABASE_URL is already set, we're likely in production, do nothing.
	if os.Getenv("DATABASE_URL") != "" {
		log.Println("DATABASE_URL already set, skipping .env loading")
		return
	}

	// Try to find .env in the project root.
	// Since the binary runs from backend/, we need to go up one level.
	envPath := filepath.Join("..", ".env")
	if _, err := os.Stat(envPath); err != nil {
		log.Fatal("No .env file found")
	}
	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf(".env file found but could not be loaded: %v", err)
	}
	log.Println("Loaded environment variables from ../.env")
}

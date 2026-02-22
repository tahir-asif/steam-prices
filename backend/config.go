package main

import (
    "log"
    "os"
    "path/filepath"

    "github.com/joho/godotenv"
)

// LoadEnv loads environment variables from .env file for local development.
// In production, environment variables are already set by the hosting platform.
func LoadEnv() {
    // If DATABASE_URL is already set, we're likely in production.
    // Do nothing.
    if os.Getenv("DATABASE_URL") != "" {
        log.Println("DATABASE_URL already set, skipping .env loading")
        return
    }

    // Try to find .env in the project root.
    // Since the binary runs from backend/, we need to go up one level.
    envPath := filepath.Join("..", ".env")
    if _, err := os.Stat(envPath); err == nil {
        if err := godotenv.Load(envPath); err != nil {
            log.Printf("Warning: .env file found but could not be loaded: %v", err)
        } else {
            log.Println("Loaded environment variables from ../.env")
        }
    } else {
        log.Println("No .env file found, using default local connection string")
    }
}

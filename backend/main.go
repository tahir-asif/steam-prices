package main

import (
    "database/sql"
    "log"
    "net/http"
    "time"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

var db *sql.DB

func main() {
    LoadEnv()

    // Connect to database
    db = ConnectDB()
    defer func() {
        // Ensure the connection is closed when the program exits
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

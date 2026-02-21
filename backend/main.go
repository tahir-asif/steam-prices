package main

import (
    "log"
    "net/http"
    "time"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()

    // CORS middleware configuration
    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5173"},
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    router.GET("/api/health", func(ginContext *gin.Context) {
        ginContext.JSON(http.StatusOK, gin.H{
            "status": "ok",
        })
    })

    log.Println("Server is running on http://localhost:3000")
    router.Run(":3000")
}

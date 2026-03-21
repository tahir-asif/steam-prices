package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func (h *Handler) DBTest(c *gin.Context) {
    var now string
    err := h.DB.QueryRow("SELECT NOW()").Scan(&now)
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
}

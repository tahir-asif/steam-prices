package handlers

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "github.com/gin-gonic/gin"
)

func (h *Handler) Search(c *gin.Context) {
    query := c.Query("q")
    if query == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing search query 'q'"})
        return
    }

    steamURL := fmt.Sprintf("https://store.steampowered.com/api/storesearch/?term=%s&cc=ca&l=english", query)
    resp, err := http.Get(steamURL)
    if err != nil {
        log.Printf("Steam search error: %v", err)
        c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to reach Steam API"})
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        c.JSON(http.StatusBadGateway, gin.H{"error": "Steam API returned an error"})
        return
    }

    var result struct {
        Items []struct {
            AppID int    `json:"id"`
            Name  string `json:"name"`
            Icon  string `json:"tiny_image"`
        } `json:"items"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        log.Printf("Failed to parse Steam search response: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid response from Steam"})
        return
    }

    type searchItem struct {
        AppID int    `json:"appid"`
        Name  string `json:"name"`
        Icon  string `json:"icon"`
    }
    items := make([]searchItem, 0, len(result.Items))
    for _, item := range result.Items {
        items = append(items, searchItem{
            AppID: item.AppID,
            Name:  item.Name,
            Icon:  item.Icon,
        })
    }
    c.JSON(http.StatusOK, items)
}

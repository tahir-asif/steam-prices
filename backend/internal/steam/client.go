// Package steam provides a client for fetching game details and pricing
// information from the public Steam Store API.
package steam

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client handles Steam API requests.
type Client struct {
	HTTPClient *http.Client
	BaseURL    string
}

// NewClient creates a default Steam API client with a 10-second timeout.
func NewClient() *Client {
	return &Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		BaseURL:    "https://store.steampowered.com/api",
	}
}

// FetchGameDetails retrieves game information using the client.
func (c *Client) FetchGameDetails(appID int) (*SteamGame, error) {
	url := fmt.Sprintf("%s/appdetails?appids=%d&cc=us", c.BaseURL, appID)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("steam API returned non-200 status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResp steamAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	appIDStr := fmt.Sprintf("%d", appID)
	details, exists := apiResp[appIDStr]
	if !exists {
		return nil, fmt.Errorf("no data found for app ID %d", appID)
	}

	if !details.Success {
		return nil, fmt.Errorf("steam API returned success=false for app ID %d", appID)
	}

	game := &SteamGame{
		Name:     details.Data.Name,
		IsFree:   details.Data.IsFree,
		Currency: details.Data.PriceOverview.Currency,
	}

	if game.IsFree {
		game.Price = 0
	} else {
		game.Price = details.Data.PriceOverview.Final
		if game.Currency == "" {
			game.Currency = "USD"
		}
	}

	return game, nil
}

type SteamGame struct {
	Name     string
	IsFree   bool
	Price    int
	Currency string
}

type appDetailsResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Name          string `json:"name"`
		IsFree        bool   `json:"is_free"`
		PriceOverview struct {
			Currency string `json:"currency"`
			Final    int    `json:"final"`
		} `json:"price_overview"`
	} `json:"data"`
}

type steamAPIResponse map[string]appDetailsResponse

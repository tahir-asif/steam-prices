package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tahir-asif/steam-prices/internal/database"
	"github.com/tahir-asif/steam-prices/internal/testutils"
)

func TestHealthEndpoint(t *testing.T) {
	db, cleanup := testutils.SetupTestDB(t)
	defer cleanup()

	router := SetupRouter(db)

	req := httptest.NewRequest("GET", "/api/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
	if resp["status"] != "ok" {
		t.Errorf("expected status ok, got %s", resp["status"])
	}
}

func TestDBTestEndpoint(t *testing.T) {
	db, cleanup := testutils.SetupTestDB(t)
	defer cleanup()

	router := SetupRouter(db)

	req := httptest.NewRequest("GET", "/api/db-test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
	if _, ok := resp["db_time"]; !ok {
		t.Error("expected db_time field in response")
	}
}

func TestPriceHistoryEndpoint(t *testing.T) {
	db, cleanup := testutils.SetupTestDB(t)
	defer cleanup()

	// Insert test data
	gameID, err := database.InsertGame(db, 999, "Test Game")
	if err != nil {
		t.Fatalf("failed to insert game: %v", err)
	}
	err = database.InsertPriceRecord(db, gameID, 1999, "CAD")
	if err != nil {
		t.Fatalf("failed to insert price record: %v", err)
	}

	router := SetupRouter(db)

	req := httptest.NewRequest("GET", "/api/games/999/history", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var history []map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &history); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
	if len(history) != 1 {
		t.Errorf("expected 1 history entry, got %d", len(history))
	}
}

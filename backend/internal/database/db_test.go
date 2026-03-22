package database

import (
	"database/sql"
	"testing"

	"github.com/tahir-asif/steam-prices/internal/testutils"
)

// TestInsertGame verifies that InsertGame correctly adds a new game
// and handles conflicts appropriately.
func TestInsertGame(t *testing.T) {
	db, cleanup := testutils.SetupTestDB(t)
	defer cleanup()

	// Insert a new game
	id, err := InsertGame(db, 730, "Counter-Strike 2")
	if err != nil {
		t.Fatalf("InsertGame failed: %v", err)
	}
	if id == 0 {
		t.Error("expected non-zero game ID")
	}

	// Insert the same game again (should update name and return existing ID)
	id2, err := InsertGame(db, 730, "CS2 Updated")
	if err != nil {
		t.Fatalf("InsertGame on conflict failed: %v", err)
	}
	if id != id2 {
		t.Errorf("expected same ID on conflict, got %d and %d", id, id2)
	}

	// Verify the name was updated
	var name string
	err = db.QueryRow(`SELECT name FROM games WHERE steam_app_id = 730`).Scan(&name)
	if err != nil {
		t.Fatalf("failed to query game: %v", err)
	}
	if name != "CS2 Updated" {
		t.Errorf("expected name 'CS2 Updated', got '%s'", name)
	}
}

// TestGetGameIDBySteamAppID tests the retrieval of internal IDs.
func TestGetGameIDBySteamAppID(t *testing.T) {
	db, cleanup := testutils.SetupTestDB(t)
	defer cleanup()

	// Insert a game first
	insertedID, err := InsertGame(db, 440, "Team Fortress 2")
	if err != nil {
		t.Fatalf("InsertGame failed: %v", err)
	}

	// Retrieve by Steam App ID
	retrievedID, err := GetGameIDBySteamAppID(db, 440)
	if err != nil {
		t.Fatalf("GetGameIDBySteamAppID failed: %v", err)
	}
	if retrievedID != insertedID {
		t.Errorf("expected ID %d, got %d", insertedID, retrievedID)
	}

	// Test non-existent game
	_, err = GetGameIDBySteamAppID(db, 999999)
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows for missing game, got %v", err)
	}
}

// TestPriceHistoryFunctions tests inserting and retrieving price history.
func TestPriceHistoryFunctions(t *testing.T) {
	db, cleanup := testutils.SetupTestDB(t)
	defer cleanup()

	// Insert a game
	gameID, err := InsertGame(db, 292030, "The Witcher 3")
	if err != nil {
		t.Fatalf("InsertGame failed: %v", err)
	}

	// Initially, there should be no price history
	_, _, err = GetLastRecordedPrice(db, gameID)
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}

	// Insert first price record
	err = InsertPriceRecord(db, gameID, 3999, "USD")
	if err != nil {
		t.Fatalf("InsertPriceRecord failed: %v", err)
	}

	// Retrieve the last recorded price
	price, currency, err := GetLastRecordedPrice(db, gameID)
	if err != nil {
		t.Fatalf("GetLastRecordedPrice failed: %v", err)
	}
	if price != 3999 {
		t.Errorf("expected price 3999, got %d", price)
	}
	if currency != "USD" {
		t.Errorf("expected currency USD, got %s", currency)
	}

	// Insert a second, different price
	err = InsertPriceRecord(db, gameID, 1999, "USD")
	if err != nil {
		t.Fatalf("InsertPriceRecord (second) failed: %v", err)
	}

	// The last price should now be the new one
	price, _, err = GetLastRecordedPrice(db, gameID)
	if err != nil {
		t.Fatalf("GetLastRecordedPrice after second insert failed: %v", err)
	}
	if price != 1999 {
		t.Errorf("expected price 1999, got %d", price)
	}
}

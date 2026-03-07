package database

import (
	"context"
	"database/sql"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	migratedb "github.com/golang-migrate/migrate/v4/database/postgres" // alias to avoid conflict
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// setupTestDB spins up a fresh PostgreSQL container, runs migrations,
// and returns a connected *sql.DB along with a cleanup function.
func setupTestDB(t *testing.T) (*sql.DB, func()) {
	ctx := context.Background()

	// Start PostgreSQL container
	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	// Get connection string
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	// Connect to the test database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	// Run migrations from the migrations directory
	if err := runMigrations(db, "../migrations"); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	// Cleanup function to close DB and terminate container
	cleanup := func() {
		db.Close()
		if err := testcontainers.TerminateContainer(pgContainer); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}

	return db, cleanup
}

// runMigrations applies all pending migrations from the specified path
// using the golang-migrate library.
func runMigrations(db *sql.DB, migrationsPath string) error {
	driver, err := migratedb.WithInstance(db, &migratedb.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

// TestInsertGame verifies that InsertGame correctly adds a new game
// and handles conflicts appropriately.
func TestInsertGame(t *testing.T) {
	db, cleanup := setupTestDB(t)
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
	db, cleanup := setupTestDB(t)
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
	db, cleanup := setupTestDB(t)
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

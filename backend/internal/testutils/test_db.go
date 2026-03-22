// Package testutils provides shared testing utilities for database setup,
// container management, and migrations.
package testutils

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	migratedb "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	_ "github.com/lib/pq"
)

// findProjectRoot walks up the directory tree until it finds a go.mod file,
// indicating the module root.
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break // reached filesystem root
		}
		dir = parent
	}
	return "", errors.New("go.mod not found")
}

// migrationsPath returns the absolute path to the migrations directory.
func migrationsPath() (string, error) {
	root, err := findProjectRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "internal", "migrations"), nil
}

// SetupTestDB spins up a fresh PostgreSQL container, runs migrations,
// and returns a connected *sql.DB along with a cleanup function.
func SetupTestDB(t *testing.T) (*sql.DB, func()) {
	t.Helper()

	ctx := context.Background()

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

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	// Automatically locate migrations
	mPath, err := migrationsPath()
	if err != nil {
		t.Fatalf("failed to find migrations directory: %v", err)
	}

	if err := RunMigrations(db, mPath); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	cleanup := func() {
		db.Close()
		if err := testcontainers.TerminateContainer(pgContainer); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}

	return db, cleanup
}

// RunMigrations applies all pending migrations from the specified path
// using the golang-migrate library.
func RunMigrations(db *sql.DB, migrationsPath string) error {
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

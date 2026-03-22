package config

import (
	"os"
	"testing"
)

func TestLoadEnv(t *testing.T) {
	// Save original environment and restore after test
	originalDBURL := os.Getenv("DATABASE_URL")
	defer os.Setenv("DATABASE_URL", originalDBURL)

	t.Run("uses existing DATABASE_URL", func(t *testing.T) {
		os.Setenv("DATABASE_URL", "postgres://prod:5432/db")
		LoadEnv()
		if os.Getenv("DATABASE_URL") != "postgres://prod:5432/db" {
			t.Error("DATABASE_URL was overwritten")
		}
	})

	t.Run("falls back when DATABASE_URL missing", func(t *testing.T) {
		os.Unsetenv("DATABASE_URL")
		LoadEnv()
		// Since .env may not exist in test, it should not panic and leave DATABASE_URL empty
		if os.Getenv("DATABASE_URL") != "" {
			t.Error("DATABASE_URL should remain empty when .env not found")
		}
	})
}

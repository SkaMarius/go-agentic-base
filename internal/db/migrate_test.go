package db

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

// TestMigrations_CreateExampleTable applies the SQL files in migrations/ to a
// real Postgres instance and confirms the example table exists afterward.
// It requires TEST_DATABASE_URL (or DATABASE_URL) to point at a disposable
// database; it skips otherwise since no test database is available.
func TestMigrations_CreateExampleTable(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = os.Getenv("DATABASE_URL")
	}
	if databaseURL == "" {
		t.Skip("set TEST_DATABASE_URL or DATABASE_URL to run migration tests against a Postgres instance")
	}

	ctx := context.Background()
	conn, err := Connect(ctx, databaseURL)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer conn.Close()

	if err := applyMigrations(ctx, conn, "../../migrations"); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}

	var exists bool
	err = conn.QueryRowContext(ctx, `SELECT EXISTS (
		SELECT 1 FROM information_schema.tables
		WHERE table_schema = 'public' AND table_name = 'example_items'
	)`).Scan(&exists)
	if err != nil {
		t.Fatalf("query table existence: %v", err)
	}
	if !exists {
		t.Fatal("expected example_items table to exist after migrations")
	}
}

// applyMigrations runs each *.up.sql file in dir, in lexical order.
func applyMigrations(ctx context.Context, conn *sql.DB, dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.up.sql"))
	if err != nil {
		return err
	}
	sort.Strings(files)
	for _, f := range files {
		sqlBytes, err := os.ReadFile(f)
		if err != nil {
			return err
		}
		if _, err := conn.ExecContext(ctx, string(sqlBytes)); err != nil {
			return err
		}
	}
	return nil
}

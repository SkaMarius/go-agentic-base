package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Connect opens a connection pool to Postgres using databaseURL and verifies
// connectivity with a ping. It fails fast with a clear error if the URL is
// missing/invalid or the database is unreachable.
func Connect(ctx context.Context, databaseURL string) (*sql.DB, error) {
	if databaseURL == "" {
		return nil, fmt.Errorf("db: DATABASE_URL is not set")
	}

	pool, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("db: invalid DATABASE_URL: %w", err)
	}

	if err := pool.PingContext(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("db: failed to connect to database: %w", err)
	}

	return pool, nil
}

// ConnectFromEnv connects using the DATABASE_URL environment variable.
func ConnectFromEnv(ctx context.Context) (*sql.DB, error) {
	return Connect(ctx, os.Getenv("DATABASE_URL"))
}

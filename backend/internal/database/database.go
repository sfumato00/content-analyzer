package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Database represents the database connection
type Database struct {
	Pool *pgxpool.Pool
}

// New creates a new database connection pool
func New(ctx context.Context, databaseURL string) (*Database, error) {
	// Configure connection pool
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %w", err)
	}

	// Connection pool settings
	config.MaxConns = 25                      // Maximum number of connections
	config.MinConns = 5                       // Minimum number of connections
	config.MaxConnLifetime = time.Hour        // Maximum connection lifetime
	config.MaxConnIdleTime = 30 * time.Minute // Maximum idle time
	config.HealthCheckPeriod = time.Minute    // Health check frequency

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	slog.Info("Database connection pool created",
		"max_conns", config.MaxConns,
		"min_conns", config.MinConns,
	)

	return &Database{Pool: pool}, nil
}

// RunMigrations runs pending database migrations
func RunMigrations(databaseURL string, migrationsPath string) error {
	slog.Info("Running database migrations", "path", migrationsPath)

	// Create migration instance
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		databaseURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	defer m.Close()

	// Run migrations
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			slog.Info("No new migrations to apply")
			return nil
		}
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Get current version
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	slog.Info("Migrations completed successfully",
		"version", version,
		"dirty", dirty,
	)

	return nil
}

// Close closes the database connection pool
func (db *Database) Close() {
	slog.Info("Closing database connection pool")
	db.Pool.Close()
}

// Ping checks if the database is reachable
func (db *Database) Ping(ctx context.Context) error {
	return db.Pool.Ping(ctx)
}

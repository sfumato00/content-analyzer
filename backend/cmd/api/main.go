package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/sfumato00/content-analyzer/internal/cache"
	"github.com/sfumato00/content-analyzer/internal/config"
	"github.com/sfumato00/content-analyzer/internal/database"
	"github.com/sfumato00/content-analyzer/internal/server"
)

func main() {
	// Load configuration from environment variables
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Configure structured logging
	setupLogging(cfg)

	// Run migrations in development mode
	if cfg.IsDevelopment() {
		slog.Info("Running database migrations (development mode)")
		if err := database.RunMigrations(cfg.DatabaseURL, "./migrations"); err != nil {
			slog.Warn("Failed to run migrations", "error", err)
		}
	}

	// Initialize database connection
	ctx := context.Background()
	db, err := database.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	slog.Info("Database connection established")

	// Initialize Redis cache
	redisCache, err := cache.New(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisCache.Close()

	// Print startup banner
	printBanner(cfg)

	// Create and start HTTP server
	srv := server.New(cfg, db, redisCache)

	slog.Info("Application starting",
		"environment", cfg.Environment,
		"port", cfg.Port,
	)

	// Start server (blocks until shutdown)
	if err := srv.Start(); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}

	slog.Info("Application stopped")
}

// setupLogging configures the structured logger
func setupLogging(cfg *config.Config) {
	var handler slog.Handler

	if cfg.IsProduction() {
		// JSON logging for production
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		// Text logging for development
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

// printBanner prints a startup banner
func printBanner(cfg *config.Config) {
	fmt.Println()
	fmt.Println(`   ____            _             _     _                _
  / ___|___  _ __ | |_ ___ _ __ | |_  / \   _ __   __ _| |_   _ ______ _ __
 | |   / _ \| '_ \| __/ _ \ '_ \| __| / _ \ | '_ \ / _` + "`" + ` | | | | |_  / _ \ '__|
 | |__| (_) | | | | ||  __/ | | | |_ / ___ \| | | | (_| | | |_| |/ /  __/ |
  \____\___/|_| |_|\__\___|_| |_|\__/_/   \_\_| |_|\__,_|_|\__, /___\___|_|
                                                            |___/`)
	fmt.Println()

	fmt.Println("  AI-Powered Content Analysis Platform")
	fmt.Println("  =====================================")
	fmt.Printf("  Environment: %s\n", cfg.Environment)
	fmt.Printf("  Port:        %s\n", cfg.Port)
	fmt.Printf("  URL:         http://localhost:%s\n", cfg.Port)
	fmt.Println("  =====================================")
	fmt.Println()
}

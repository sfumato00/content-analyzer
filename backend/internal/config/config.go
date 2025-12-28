package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	// API Keys
	GeminiAPIKey string

	// Database
	DatabaseURL string

	// Redis
	RedisURL string

	// Authentication
	JWTSecret string

	// Server
	Port           string
	Environment    string
	AllowedOrigins []string
}

// Load reads configuration from environment variables
// In development, it loads from .env file
func Load() (*Config, error) {
	// Load .env file only in non-production environments
	env := os.Getenv("ENV")
	if env != "production" {
		// Attempt to load from current directory, then parent directory
		_ = godotenv.Load()
		_ = godotenv.Load("../.env")
	}

	cfg := &Config{
		GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		RedisURL:     os.Getenv("REDIS_URL"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
		Port:         getEnvOrDefault("PORT", "8080"),
		Environment:  getEnvOrDefault("ENV", "development"),
	}

	// Parse allowed origins (comma-separated)
	if origins := os.Getenv("ALLOWED_ORIGINS"); origins != "" {
		cfg.AllowedOrigins = parseCommaSeparated(origins)
	} else {
		// Default for development
		cfg.AllowedOrigins = []string{"http://localhost:3000", "http://localhost:8080"}
	}

	// Validate required configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks that all required configuration is present
func (c *Config) Validate() error {
	if c.GeminiAPIKey == "" {
		return fmt.Errorf("GEMINI_API_KEY environment variable is required")
	}

	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL environment variable is required")
	}

	if c.RedisURL == "" {
		return fmt.Errorf("REDIS_URL environment variable is required")
	}

	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET environment variable is required")
	}

	// Validate JWT secret length (should be at least 32 characters for security)
	if len(c.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters long")
	}

	return nil
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// getEnvOrDefault returns the value of an environment variable or a default value
func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// getEnvAsBool returns an environment variable as a boolean
func getEnvAsBool(key string, defaultVal bool) bool {
	if val := os.Getenv(key); val != "" {
		b, err := strconv.ParseBool(val)
		if err != nil {
			return defaultVal
		}
		return b
	}
	return defaultVal
}

// parseCommaSeparated parses a comma-separated string into a slice
func parseCommaSeparated(s string) []string {
	var result []string
	for _, item := range splitAndTrim(s, ',') {
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}

// splitAndTrim splits a string and trims whitespace from each element
func splitAndTrim(s string, sep rune) []string {
	var result []string
	var current string

	for _, char := range s {
		if char == sep {
			trimmed := trimSpace(current)
			if trimmed != "" {
				result = append(result, trimmed)
			}
			current = ""
		} else {
			current += string(char)
		}
	}

	// Add the last element
	if trimmed := trimSpace(current); trimmed != "" {
		result = append(result, trimmed)
	}

	return result
}

// trimSpace removes leading and trailing whitespace
func trimSpace(s string) string {
	start := 0
	end := len(s)

	// Trim leading spaces
	for start < end && isSpace(s[start]) {
		start++
	}

	// Trim trailing spaces
	for end > start && isSpace(s[end-1]) {
		end--
	}

	return s[start:end]
}

// isSpace checks if a byte is a whitespace character
func isSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

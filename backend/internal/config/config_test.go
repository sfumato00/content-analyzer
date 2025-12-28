package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Set required environment variables for testing
	os.Setenv("GEMINI_API_KEY", "test-api-key-1234567890")
	os.Setenv("DATABASE_URL", "postgresql://postgres:password@localhost:5432/test_db")
	os.Setenv("REDIS_URL", "redis://localhost:6379")
	os.Setenv("JWT_SECRET", "this-is-a-test-secret-that-is-at-least-32-characters-long")
	os.Setenv("PORT", "8080")
	os.Setenv("ENV", "test")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify all values are loaded correctly
	if cfg.GeminiAPIKey != "test-api-key-1234567890" {
		t.Errorf("Expected GeminiAPIKey to be 'test-api-key-1234567890', got '%s'", cfg.GeminiAPIKey)
	}

	if cfg.DatabaseURL != "postgresql://postgres:password@localhost:5432/test_db" {
		t.Errorf("Expected DatabaseURL to match, got '%s'", cfg.DatabaseURL)
	}

	if cfg.RedisURL != "redis://localhost:6379" {
		t.Errorf("Expected RedisURL to be 'redis://localhost:6379', got '%s'", cfg.RedisURL)
	}

	if cfg.Port != "8080" {
		t.Errorf("Expected Port to be '8080', got '%s'", cfg.Port)
	}

	if cfg.Environment != "test" {
		t.Errorf("Expected Environment to be 'test', got '%s'", cfg.Environment)
	}
}

func TestValidate_MissingGeminiAPIKey(t *testing.T) {
	cfg := &Config{
		DatabaseURL: "postgresql://localhost/test",
		RedisURL:    "redis://localhost:6379",
		JWTSecret:   "this-is-a-test-secret-at-least-32-chars",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation to fail when GEMINI_API_KEY is missing")
	}

	if err.Error() != "GEMINI_API_KEY environment variable is required" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestValidate_MissingDatabaseURL(t *testing.T) {
	cfg := &Config{
		GeminiAPIKey: "test-key",
		RedisURL:     "redis://localhost:6379",
		JWTSecret:    "this-is-a-test-secret-at-least-32-chars",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation to fail when DATABASE_URL is missing")
	}

	if err.Error() != "DATABASE_URL environment variable is required" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestValidate_ShortJWTSecret(t *testing.T) {
	cfg := &Config{
		GeminiAPIKey: "test-key",
		DatabaseURL:  "postgresql://localhost/test",
		RedisURL:     "redis://localhost:6379",
		JWTSecret:    "short", // Less than 32 characters
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation to fail when JWT_SECRET is too short")
	}

	if err.Error() != "JWT_SECRET must be at least 32 characters long" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestIsDevelopment(t *testing.T) {
	cfg := &Config{Environment: "development"}
	if !cfg.IsDevelopment() {
		t.Error("Expected IsDevelopment() to return true for development environment")
	}

	cfg.Environment = "production"
	if cfg.IsDevelopment() {
		t.Error("Expected IsDevelopment() to return false for production environment")
	}
}

func TestIsProduction(t *testing.T) {
	cfg := &Config{Environment: "production"}
	if !cfg.IsProduction() {
		t.Error("Expected IsProduction() to return true for production environment")
	}

	cfg.Environment = "development"
	if cfg.IsProduction() {
		t.Error("Expected IsProduction() to return false for development environment")
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	os.Setenv("TEST_VAR", "custom_value")
	defer os.Unsetenv("TEST_VAR")

	result := getEnvOrDefault("TEST_VAR", "default_value")
	if result != "custom_value" {
		t.Errorf("Expected 'custom_value', got '%s'", result)
	}

	result = getEnvOrDefault("NON_EXISTENT_VAR", "default_value")
	if result != "default_value" {
		t.Errorf("Expected 'default_value', got '%s'", result)
	}
}

func TestParseCommaSeparated(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			input:    "http://localhost:3000,http://localhost:8080",
			expected: []string{"http://localhost:3000", "http://localhost:8080"},
		},
		{
			input:    "http://localhost:3000, http://localhost:8080",
			expected: []string{"http://localhost:3000", "http://localhost:8080"},
		},
		{
			input:    "single-value",
			expected: []string{"single-value"},
		},
		{
			input:    "",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		result := parseCommaSeparated(tt.input)
		if len(result) != len(tt.expected) {
			t.Errorf("For input '%s', expected %d items, got %d", tt.input, len(tt.expected), len(result))
			continue
		}

		for i, val := range result {
			if val != tt.expected[i] {
				t.Errorf("For input '%s', expected item %d to be '%s', got '%s'", tt.input, i, tt.expected[i], val)
			}
		}
	}
}

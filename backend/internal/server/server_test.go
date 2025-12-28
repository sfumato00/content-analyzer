package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/content-analyzer/internal/config"
)

func TestHealthEndpoint(t *testing.T) {
	// Create test configuration
	cfg := &config.Config{
		Port:        "8080",
		Environment: "test",
		GeminiAPIKey: "test-key-for-testing-purposes-only-12345",
		DatabaseURL: "postgresql://localhost/test",
		RedisURL:    "redis://localhost:6379",
		JWTSecret:   "test-secret-that-is-at-least-32-characters-long",
		AllowedOrigins: []string{"http://localhost:3000"},
	}

	// Create server
	srv := New(cfg)

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	// Serve request
	srv.Router().ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Parse response
	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Check response fields
	if status, ok := response["status"].(string); !ok || status != "healthy" {
		t.Errorf("Expected status 'healthy', got %v", response["status"])
	}

	if _, ok := response["uptime"].(string); !ok {
		t.Error("Expected uptime field in response")
	}

	if version, ok := response["version"].(string); !ok || version == "" {
		t.Error("Expected non-empty version field in response")
	}
}

func TestRootEndpoint(t *testing.T) {
	cfg := &config.Config{
		Port:        "8080",
		Environment: "test",
		GeminiAPIKey: "test-key-for-testing-purposes-only-12345",
		DatabaseURL: "postgresql://localhost/test",
		RedisURL:    "redis://localhost:6379",
		JWTSecret:   "test-secret-that-is-at-least-32-characters-long",
		AllowedOrigins: []string{"http://localhost:3000"},
	}

	srv := New(cfg)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	srv.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if name, ok := response["name"].(string); !ok || name == "" {
		t.Error("Expected non-empty name field in response")
	}
}

func TestNotFoundEndpoint(t *testing.T) {
	cfg := &config.Config{
		Port:        "8080",
		Environment: "test",
		GeminiAPIKey: "test-key-for-testing-purposes-only-12345",
		DatabaseURL: "postgresql://localhost/test",
		RedisURL:    "redis://localhost:6379",
		JWTSecret:   "test-secret-that-is-at-least-32-characters-long",
		AllowedOrigins: []string{"http://localhost:3000"},
	}

	srv := New(cfg)

	req := httptest.NewRequest(http.MethodGet, "/does-not-exist", nil)
	w := httptest.NewRecorder()

	srv.Router().ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestCORS(t *testing.T) {
	cfg := &config.Config{
		Port:        "8080",
		Environment: "test",
		GeminiAPIKey: "test-key-for-testing-purposes-only-12345",
		DatabaseURL: "postgresql://localhost/test",
		RedisURL:    "redis://localhost:6379",
		JWTSecret:   "test-secret-that-is-at-least-32-characters-long",
		AllowedOrigins: []string{"http://localhost:3000"},
	}

	srv := New(cfg)

	req := httptest.NewRequest(http.MethodOptions, "/health", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "GET")

	w := httptest.NewRecorder()
	srv.Router().ServeHTTP(w, req)

	// Check CORS headers are set
	if w.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
		t.Errorf("Expected CORS origin header to be set")
	}
}

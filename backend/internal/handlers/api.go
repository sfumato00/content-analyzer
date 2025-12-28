package handlers

import (
	"net/http"

	"github.com/sfumato00/content-analyzer/internal/config"
	"github.com/sfumato00/content-analyzer/internal/response"
)

// APIHandler handles general API requests
type APIHandler struct {
	config *config.Config
}

// NewAPIHandler creates a new API handler
func NewAPIHandler(cfg *config.Config) *APIHandler {
	return &APIHandler{
		config: cfg,
	}
}

// Index returns API information
func (h *APIHandler) Index(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]interface{}{
		"name":        "Content Analyzer API",
		"version":     "1.0.0",
		"environment": h.config.Environment,
		"endpoints": map[string]string{
			"health":   "/health",
			"ready":    "/ready",
			"live":     "/live",
			"api_root": "/api/v1",
		},
	})
}

// NotFound handles 404 errors
func (h *APIHandler) NotFound(w http.ResponseWriter, r *http.Request) {
	response.NotFound(w, "The requested resource was not found")
}

// MethodNotAllowed handles 405 errors
func (h *APIHandler) MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
}

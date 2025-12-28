package handlers

import (
	"net/http"
	"time"

	"github.com/yourusername/content-analyzer/internal/response"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	startTime time.Time
}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
	}
}

// Health returns the health status of the application
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(h.startTime)

	response.Success(w, map[string]interface{}{
		"status":  "healthy",
		"uptime":  uptime.String(),
		"version": "1.0.0",
	})
}

// Ready returns readiness status (useful for Kubernetes readiness probes)
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	// TODO: Add checks for database, Redis, etc.
	// For now, just return OK
	response.Success(w, map[string]interface{}{
		"status": "ready",
	})
}

// Live returns liveness status (useful for Kubernetes liveness probes)
func (h *HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]interface{}{
		"status": "alive",
	})
}

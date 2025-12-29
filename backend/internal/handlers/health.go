package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/sfumato00/content-analyzer/internal/cache"
	"github.com/sfumato00/content-analyzer/internal/database"
	"github.com/sfumato00/content-analyzer/internal/response"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	startTime time.Time
	db        *database.Database
	cache     *cache.Cache
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *database.Database, cache *cache.Cache) *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
		db:        db,
		cache:     cache,
	}
}

// Health returns the health status of the application
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	uptime := time.Since(h.startTime)

	// Check component health
	components := make(map[string]string)

	// Check database
	if err := h.db.Ping(ctx); err != nil {
		components["database"] = "disconnected"
	} else {
		components["database"] = "connected"
	}

	// Check Redis
	if err := h.cache.Ping(ctx); err != nil {
		components["redis"] = "disconnected"
	} else {
		components["redis"] = "connected"
	}

	// Overall status is healthy only if all components are connected
	status := "healthy"
	if components["database"] != "connected" || components["redis"] != "connected" {
		status = "degraded"
	}

	response.Success(w, map[string]interface{}{
		"status":     status,
		"uptime":     uptime.String(),
		"version":    "1.0.0",
		"components": components,
	})
}

// Ready returns readiness status (useful for Kubernetes readiness probes)
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	// Check if database is ready
	if err := h.db.Ping(ctx); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		response.Error(w, http.StatusServiceUnavailable, "database not ready")
		return
	}

	// Check if Redis is ready
	if err := h.cache.Ping(ctx); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		response.Error(w, http.StatusServiceUnavailable, "redis not ready")
		return
	}

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

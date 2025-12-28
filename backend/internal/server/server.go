package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"

	"github.com/sfumato00/content-analyzer/internal/config"
	"github.com/sfumato00/content-analyzer/internal/handlers"
	custommw "github.com/sfumato00/content-analyzer/internal/middleware"
)

// Server represents the HTTP server
type Server struct {
	config     *config.Config
	router     *chi.Mux
	httpServer *http.Server
}

// New creates a new server instance
func New(cfg *config.Config) *Server {
	s := &Server{
		config: cfg,
		router: chi.NewRouter(),
	}

	s.setupMiddleware()
	s.setupRoutes()

	s.httpServer = &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return s
}

// setupMiddleware configures all middleware
func (s *Server) setupMiddleware() {
	// Logger middleware
	logger := httplog.NewLogger("content-analyzer", httplog.Options{
		JSON:             s.config.IsProduction(),
		LogLevel:         slog.LevelInfo,
		Concise:          true,
		RequestHeaders:   true,
		MessageFieldName: "message",
		Tags: map[string]string{
			"env": s.config.Environment,
		},
	})

	s.router.Use(httplog.RequestLogger(logger))

	// Recoverer - recover from panics
	s.router.Use(middleware.Recoverer)

	// Request ID
	s.router.Use(middleware.RequestID)

	// Real IP
	s.router.Use(middleware.RealIP)

	// Timeout
	s.router.Use(middleware.Timeout(30 * time.Second))

	// Compress responses
	s.router.Use(middleware.Compress(5))

	// Security headers
	s.router.Use(custommw.SecurityHeaders)

	// CORS
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   s.config.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Heartbeat endpoint (doesn't log)
	s.router.Use(middleware.Heartbeat("/ping"))
}

// setupRoutes configures all routes
func (s *Server) setupRoutes() {
	// Create handlers
	healthHandler := handlers.NewHealthHandler()
	apiHandler := handlers.NewAPIHandler(s.config)

	// Root endpoint
	s.router.Get("/", apiHandler.Index)

	// Health check endpoints
	s.router.Get("/health", healthHandler.Health)
	s.router.Get("/ready", healthHandler.Ready)
	s.router.Get("/live", healthHandler.Live)

	// API v1 routes
	s.router.Route("/api/v1", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "API v1", http.StatusOK)
		})

		// Auth routes (TODO: implement)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "TODO: Register endpoint", http.StatusNotImplemented)
			})
			r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "TODO: Login endpoint", http.StatusNotImplemented)
			})
			r.Post("/logout", func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "TODO: Logout endpoint", http.StatusNotImplemented)
			})
		})

		// Submissions routes (TODO: implement)
		r.Route("/submissions", func(r chi.Router) {
			// TODO: Add JWT middleware here
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "TODO: List submissions", http.StatusNotImplemented)
			})
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "TODO: Create submission", http.StatusNotImplemented)
			})
			r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "TODO: Get submission", http.StatusNotImplemented)
			})
			r.Get("/{id}/analysis", func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "TODO: Get analysis", http.StatusNotImplemented)
			})
		})

		// User routes (TODO: implement)
		r.Route("/me", func(r chi.Router) {
			// TODO: Add JWT middleware here
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "TODO: Get current user", http.StatusNotImplemented)
			})
			r.Get("/stats", func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "TODO: Get user stats", http.StatusNotImplemented)
			})
		})
	})

	// 404 handler
	s.router.NotFound(apiHandler.NotFound)

	// 405 handler
	s.router.MethodNotAllowed(apiHandler.MethodNotAllowed)
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Print routes in development
	if s.config.IsDevelopment() {
		s.printRoutes()
	}

	slog.Info("Starting HTTP server",
		"port", s.config.Port,
		"env", s.config.Environment,
	)

	// Channel to listen for errors from the server
	serverErrors := make(chan error, 1)

	// Start the server in a goroutine
	go func() {
		serverErrors <- s.httpServer.ListenAndServe()
	}()

	// Channel to listen for interrupt signal
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block until we receive a signal or error
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		slog.Info("Shutdown signal received", "signal", sig.String())

		// Give outstanding requests 30 seconds to complete
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Shutdown the server gracefully
		if err := s.httpServer.Shutdown(ctx); err != nil {
			// Force close if graceful shutdown fails
			s.httpServer.Close()
			return fmt.Errorf("failed to gracefully shutdown server: %w", err)
		}

		slog.Info("Server stopped gracefully")
	}

	return nil
}

// printRoutes prints all registered routes (development only)
func (s *Server) printRoutes() {
	fmt.Println("\n📍 Registered Routes:")
	fmt.Println("====================")

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if route != "/ping" { // Skip heartbeat
			fmt.Printf("%-6s %s\n", method, route)
		}
		return nil
	}

	if err := chi.Walk(s.router, walkFunc); err != nil {
		slog.Error("Failed to walk routes", "error", err)
	}

	fmt.Println("====================")
	fmt.Println()
}

// Router returns the Chi router (useful for testing)
func (s *Server) Router() *chi.Mux {
	return s.router
}

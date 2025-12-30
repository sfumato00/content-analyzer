package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/sfumato00/content-analyzer/internal/auth"
	"github.com/sfumato00/content-analyzer/internal/models"
	"github.com/sfumato00/content-analyzer/internal/response"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	userStore  *models.UserStore
	jwtManager *auth.JWTManager
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userStore *models.UserStore, jwtManager *auth.JWTManager) *AuthHandler {
	return &AuthHandler{
		userStore:  userStore,
		jwtManager: jwtManager,
	}
}

// RegisterRequest represents the registration request
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	User  *UserResponse       `json:"user"`
	Token *auth.TokenPair     `json:"token"`
}

// UserResponse represents the user data in responses (without sensitive fields)
type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	// Normalize email
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// Create user
	user, err := h.userStore.Create(r.Context(), req.Email, req.Password)
	if err != nil {
		// Check for duplicate email error
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			response.BadRequest(w, "Email already exists")
			return
		}

		// Check for validation errors
		if strings.Contains(err.Error(), "email") || strings.Contains(err.Error(), "password") {
			response.BadRequest(w, err.Error())
			return
		}

		slog.Error("Failed to create user", "error", err)
		response.InternalServerError(w, "Failed to create user")
		return
	}

	// Generate JWT token
	tokenPair, err := h.jwtManager.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		slog.Error("Failed to generate token", "error", err)
		response.InternalServerError(w, "Failed to generate authentication token")
		return
	}

	// Return user and token
	authResp := AuthResponse{
		User: &UserResponse{
			ID:        user.ID.String(),
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
		Token: tokenPair,
	}

	response.Created(w, authResp)
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	// Normalize email
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// Get user by email
	user, err := h.userStore.GetByEmail(r.Context(), req.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			response.Unauthorized(w, "Invalid email or password")
			return
		}

		slog.Error("Failed to get user", "error", err)
		response.InternalServerError(w, "Failed to authenticate")
		return
	}

	// Compare password
	if err := user.ComparePassword(req.Password); err != nil {
		response.Unauthorized(w, "Invalid email or password")
		return
	}

	// Generate JWT token
	tokenPair, err := h.jwtManager.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		slog.Error("Failed to generate token", "error", err)
		response.InternalServerError(w, "Failed to generate authentication token")
		return
	}

	// Return user and token
	authResp := AuthResponse{
		User: &UserResponse{
			ID:        user.ID.String(),
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
		Token: tokenPair,
	}

	response.Success(w, authResp)
}

// Logout handles user logout
// Note: Since we're using JWT, logout is primarily client-side
// The client should remove the token from storage
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// For JWT, logout is handled client-side by removing the token
	// In the future, we could implement token blacklisting using Redis
	response.Success(w, map[string]string{
		"message": "Logged out successfully",
	})
}

// Me returns the current authenticated user
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	// Get user from database
	user, err := h.userStore.GetByID(r.Context(), userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			response.NotFound(w, "User not found")
			return
		}

		slog.Error("Failed to get user", "error", err)
		response.InternalServerError(w, "Failed to get user")
		return
	}

	// Return user
	userResp := UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	response.Success(w, userResp)
}

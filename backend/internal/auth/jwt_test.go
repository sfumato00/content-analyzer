package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWTManager_GenerateTokenPair(t *testing.T) {
	secret := "test-secret-key-at-least-32-characters-long"
	jwtManager := NewJWTManager(secret)

	userID := uuid.New()
	email := "test@example.com"

	tokenPair, err := jwtManager.GenerateTokenPair(userID, email)
	if err != nil {
		t.Fatalf("GenerateTokenPair() error = %v", err)
	}

	if tokenPair.AccessToken == "" {
		t.Error("GenerateTokenPair() returned empty access token")
	}

	if tokenPair.TokenType != "Bearer" {
		t.Errorf("GenerateTokenPair() token type = %v, want Bearer", tokenPair.TokenType)
	}

	if tokenPair.ExpiresAt.Before(time.Now()) {
		t.Error("GenerateTokenPair() ExpiresAt is in the past")
	}
}

func TestJWTManager_ValidateToken(t *testing.T) {
	secret := "test-secret-key-at-least-32-characters-long"
	jwtManager := NewJWTManager(secret)

	userID := uuid.New()
	email := "test@example.com"

	// Generate token
	tokenPair, err := jwtManager.GenerateTokenPair(userID, email)
	if err != nil {
		t.Fatalf("GenerateTokenPair() error = %v", err)
	}

	// Validate token
	claims, err := jwtManager.ValidateToken(tokenPair.AccessToken)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("ValidateToken() UserID = %v, want %v", claims.UserID, userID)
	}

	if claims.Email != email {
		t.Errorf("ValidateToken() Email = %v, want %v", claims.Email, email)
	}
}

func TestJWTManager_ValidateToken_Invalid(t *testing.T) {
	secret := "test-secret-key-at-least-32-characters-long"
	jwtManager := NewJWTManager(secret)

	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "empty token",
			token: "",
		},
		{
			name:  "invalid token",
			token: "invalid.token.here",
		},
		{
			name:  "malformed token",
			token: "not-a-jwt-token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := jwtManager.ValidateToken(tt.token)
			if err == nil {
				t.Error("ValidateToken() expected error for invalid token")
			}
		})
	}
}

func TestJWTManager_ValidateToken_WrongSecret(t *testing.T) {
	secret1 := "test-secret-key-at-least-32-characters-long"
	secret2 := "different-secret-key-at-least-32-chars-long"

	jwtManager1 := NewJWTManager(secret1)
	jwtManager2 := NewJWTManager(secret2)

	userID := uuid.New()
	email := "test@example.com"

	// Generate token with first secret
	tokenPair, err := jwtManager1.GenerateTokenPair(userID, email)
	if err != nil {
		t.Fatalf("GenerateTokenPair() error = %v", err)
	}

	// Try to validate with second secret (should fail)
	_, err = jwtManager2.ValidateToken(tokenPair.AccessToken)
	if err == nil {
		t.Error("ValidateToken() should fail with different secret")
	}
}

func TestJWTManager_ExtractUserID(t *testing.T) {
	secret := "test-secret-key-at-least-32-characters-long"
	jwtManager := NewJWTManager(secret)

	expectedUserID := uuid.New()
	email := "test@example.com"

	// Generate token
	tokenPair, err := jwtManager.GenerateTokenPair(expectedUserID, email)
	if err != nil {
		t.Fatalf("GenerateTokenPair() error = %v", err)
	}

	// Extract user ID
	userID, err := jwtManager.ExtractUserID(tokenPair.AccessToken)
	if err != nil {
		t.Fatalf("ExtractUserID() error = %v", err)
	}

	if userID != expectedUserID {
		t.Errorf("ExtractUserID() = %v, want %v", userID, expectedUserID)
	}
}

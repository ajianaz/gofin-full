package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/gofiber/fiber/v2"

	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors"
)

// KeyLookup is implemented by the API key repository.
// This interface breaks the import cycle between auth and repository.
type KeyLookup interface {
	FindByHash(ctx context.Context, keyHash string) (userID int64, keyID int64, err error)
	UpdateLastUsed(ctx context.Context, keyID int64) error
}

// APIKeyMiddleware checks for API key authentication.
// It checks the Authorization header (Bearer gofin_*) or X-API-Key header.
// If a valid key is found, it sets the user context.
// If a gofin_ key is invalid, it rejects immediately (does not fall through to JWT).
// If no gofin_ key is present, it passes through to the next middleware (JWT).
func APIKeyMiddleware(lookup KeyLookup) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rawKey := extractAPIKey(c)
		if rawKey == "" {
			return c.Next()
		}

		keyHash := HashAPIKey(rawKey)
		userID, keyID, err := lookup.FindByHash(c.Context(), keyHash)
		if err != nil {
			// Key had gofin_ prefix but lookup failed — reject
			return apperrors.ErrUnauthorized
		}

		// Update last used (best-effort)
		go func() {
			_ = lookup.UpdateLastUsed(context.Background(), keyID)
		}()

		SetUser(c, &UserIdentity{
			ID:    userID,
			Email: "",
		})
		c.Locals("auth_method", "api_key")

		return c.Next()
	}
}

// extractAPIKey extracts an API key from the request headers.
func extractAPIKey(c *fiber.Ctx) string {
	authHeader := c.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if strings.HasPrefix(token, "gofin_") {
			return token
		}
	}

	apiKey := c.Get("X-API-Key")
	if strings.HasPrefix(apiKey, "gofin_") {
		return apiKey
	}

	return ""
}

// HashAPIKey computes the SHA-256 hash of an API key string.
func HashAPIKey(rawKey string) string {
	hash := sha256.Sum256([]byte(rawKey))
	return hex.EncodeToString(hash[:])
}

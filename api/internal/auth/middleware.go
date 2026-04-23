package auth

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors"
)

// TokenVersionLookup is implemented by repositories that can check a user's token version.
// This interface breaks the import cycle between auth and repository.
type TokenVersionLookup interface {
	GetTokenVersion(ctx context.Context, userID uuid.UUID) (int, error)
}

// ErrTokenInvalidated is returned when a JWT's token_version does not match the DB value.
var ErrTokenInvalidated = apperrors.NewWithDetail(401, "Unauthenticated", "Token has been invalidated. Please log in again.")

// AuthMiddleware creates a Fiber middleware that validates JWT tokens.
func AuthMiddleware(jwtMgr *JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip if already authenticated by API key middleware
		if c.Locals("auth_method") == "api_key" {
			return c.Next()
		}

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return apperrors.ErrUnauthorized
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return apperrors.ErrUnauthorized
		}

		claims, err := jwtMgr.ValidateAccessToken(parts[1])
		if err != nil {
			return apperrors.ErrUnauthorized
		}

		// Validate token_version against the DB to detect invalidated tokens.
		// The version lookup is injected via c.Locals("token_version_lookup") which is set
		// by the TokenVersionMiddleware that runs before this middleware.
		if lookup, ok := c.Locals("token_version_lookup").(TokenVersionLookup); ok {
			dbVersion, err := lookup.GetTokenVersion(c.Context(), claims.UserID)
			if err == nil && dbVersion != claims.TokenVersion {
				return ErrTokenInvalidated
			}
		}

		// Store claims and user identity in context
		SetClaims(c, claims)
		SetUser(c, &UserIdentity{
			ID:       claims.UserID,
			Email:    claims.Email,
			DemoUser: claims.DemoUser,
		})

		// Set active group from claims only (not from query param)
		// Group must be set via POST /groups/switch which validates membership
		if claims.GroupID != nil {
			SetActiveGroupID(c, *claims.GroupID)
		}

		return c.Next()
	}
}

// OptionalAuthMiddleware validates JWT if present, but doesn't reject unauthenticated requests.
func OptionalAuthMiddleware(jwtMgr *JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return c.Next()
		}

		claims, err := jwtMgr.ValidateAccessToken(parts[1])
		if err != nil {
			return c.Next()
		}

		// Validate token_version against the DB to detect invalidated tokens.
		if lookup, ok := c.Locals("token_version_lookup").(TokenVersionLookup); ok {
			dbVersion, err := lookup.GetTokenVersion(c.Context(), claims.UserID)
			if err == nil && dbVersion != claims.TokenVersion {
				return c.Next() // Optional auth: just skip, don't reject
			}
		}

		SetClaims(c, claims)
		SetUser(c, &UserIdentity{
			ID:       claims.UserID,
			Email:    claims.Email,
			DemoUser: claims.DemoUser,
		})

		if claims.GroupID != nil {
			SetActiveGroupID(c, *claims.GroupID)
		}

		return c.Next()
	}
}

// DemoUserMiddleware blocks demo users from destructive operations.
func DemoUserMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := GetUser(c)
		if user != nil && user.DemoUser {
			return apperrors.DemoUserBlocked()
		}
		return c.Next()
	}
}

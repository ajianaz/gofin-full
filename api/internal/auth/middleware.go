package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors"
)

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

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

// TokenInvalidator is implemented by caches that need to be invalidated on auth events.
type TokenInvalidator interface {
	InvalidateTokenVersion(ctx context.Context, userID uuid.UUID)
}

// ErrTokenInvalidated is returned when a JWT's token_version does not match the DB value.
var ErrTokenInvalidated = apperrors.NewWithDetail(401, "Unauthenticated", "Token has been invalidated. Please log in again.")

// Cookie names for httpOnly token storage.
const (
	// AccessTokenCookieName is the name of the httpOnly cookie holding the access token.
	AccessTokenCookieName = "gofin_access_token"
	// RefreshTokenCookieName is the name of the httpOnly cookie holding the refresh token.
	RefreshTokenCookieName = "gofin_refresh_token"
	// AccessTokenMaxAge is the access token cookie MaxAge in seconds (15 minutes).
	AccessTokenMaxAge = 900
	// RefreshTokenMaxAge is the refresh token cookie MaxAge in seconds (7 days).
	RefreshTokenMaxAge = 604800
)


// extractToken retrieves the JWT access token from the request.
// Priority: httpOnly cookie first, then Authorization: Bearer ***
// This allows gradual migration from header-based to cookie-based auth.
// Note: Token validation (including token_version invalidation check) still
// happens in AuthMiddleware/OptionalAuthMiddleware after extractToken returns.
// The cookie contains the same JWT — only the transport mechanism differs.
func extractToken(c *fiber.Ctx) string {
	// 1. Try httpOnly cookie
	if tok := c.Cookies(AccessTokenCookieName); tok != "" {
		return tok
	}
	// 2. Fall back to Authorization header (backward compat)
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return ""
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}
	return ""
}

// SetTokenCookies sets httpOnly, Secure, SameSite=Lax cookies for both access and refresh tokens.
// The response body still includes the tokens for backward compatibility.
// The refresh token cookie is set for future use if the refresh endpoint is
// updated to read from cookies. Currently the refresh endpoint reads from
// the request body (existing API contract). The access token cookie is the
// primary mechanism — sent automatically by the browser on every request.
func SetTokenCookies(c *fiber.Ctx, accessToken, refreshToken string) {
	c.Cookie(&fiber.Cookie{
		Name:     AccessTokenCookieName,
		Value:    accessToken,
		MaxAge:   AccessTokenMaxAge,
		Path:     "/",
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Lax",
	})
	c.Cookie(&fiber.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    refreshToken,
		MaxAge:   RefreshTokenMaxAge,
		Path:     "/",
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Lax",
	})
}

// ClearTokenCookies removes the httpOnly token cookies (used on logout).
func ClearTokenCookies(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     AccessTokenCookieName,
		Value:    "",
		MaxAge:   -1, // delete immediately
		Path:     "/",
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Lax",
	})
	c.Cookie(&fiber.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    "",
		MaxAge:   -1, // delete immediately
		Path:     "/",
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Lax",
	})
}

// AuthMiddleware creates a Fiber middleware that validates JWT tokens.
func AuthMiddleware(jwtMgr *JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip if already authenticated by API key middleware
		if c.Locals("auth_method") == "api_key" {
			return c.Next()
		}

		token := extractToken(c)
		if token == "" {
			return apperrors.ErrUnauthorized
		}

		claims, err := jwtMgr.ValidateAccessToken(token)
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
		token := extractToken(c)
		if token == "" {
			return c.Next()
		}

		claims, err := jwtMgr.ValidateAccessToken(token)
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

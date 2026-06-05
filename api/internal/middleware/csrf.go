package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
)

// CSRFConfig holds configuration for the CSRF middleware wrapper.
type CSRFConfig struct {
	Secret  string
	IsProd  bool
	IsDebug bool
}

// CSRF returns a Fiber CSRF middleware configured for the double-submit cookie pattern.
//
// Cookie: "gofin_csrf" (non-httpOnly so JS can read it; Secure in production; SameSite=Lax).
// Header: "X-CSRF-Token".
//
// API-key authenticated requests are exempt (c.Locals("auth_method") == "api_key").
// Note: The CSRF middleware runs BEFORE auth middleware, so auth_method is set by a
// lightweight pre-check that inspects the Authorization header for API key patterns.
func CSRF(cfg CSRFConfig) fiber.Handler {
	h := csrf.New(csrf.Config{
		KeyLookup:   "header:X-CSRF-Token",
		CookieName:  "gofin_csrf",
		CookieSameSite: "Lax",
		CookieSecure:   cfg.IsProd,
		CookieHTTPOnly: false,
		CookiePath:     "/",
		Expiration:     1 * time.Hour,
		ContextKey:     "csrf_token",
		// Use the configured secret for deterministic key generation.
		// This ensures CSRF cookies survive server restarts.
		KeyGenerator: func() string {
			if cfg.Secret != "" {
				return cfg.Secret
			}
			// Fallback: random 32-byte hex key (cookies invalidate on restart)
			b := make([]byte, 32)
			_, _ = rand.Read(b)
			return hex.EncodeToString(b)
		},
		// Skip CSRF for API-key requests — detect via Authorization header
		// pattern BEFORE auth middleware runs. API keys use "gofin_" prefix
		// in the Authorization header or X-API-Key header.
		Next: func(c *fiber.Ctx) bool {
			// Check for API key in Authorization header (format: "Bearer gofin_..." or "gofin_...")
			if authHeader := c.Get("Authorization"); authHeader != "" {
				// Strip "Bearer " prefix if present
				token := authHeader
				if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
					token = authHeader[7:]
				}
				if len(token) > 6 && token[:6] == "gofin_" {
					c.Locals("auth_method", "api_key")
					return true
				}
			}
			// Check X-API-Key header
			if apiKey := c.Get("X-API-Key"); apiKey != "" {
				if len(apiKey) > 6 && apiKey[:6] == "gofin_" {
					c.Locals("auth_method", "api_key")
					return true
				}
			}
			return false
		},
	})
	return h
}

// CSRFTokenHandler returns a handler that responds with the current CSRF token.
// This endpoint is used by the frontend to obtain an initial CSRF token/cookie.
func CSRFTokenHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.Locals("csrf_token").(string)
		if !ok || token == "" {
			return c.Status(500).JSON(fiber.Map{
				"error": "CSRF token not available",
			})
		}
		return c.JSON(fiber.Map{
			"token": token,
		})
	}
}

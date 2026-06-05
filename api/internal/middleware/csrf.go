package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
)

// CSRFConfig holds configuration for the CSRF middleware wrapper.
type CSRFConfig struct {
	Secret string
	IsProd bool
}

// CSRF returns a Fiber CSRF middleware configured for the double-submit cookie pattern.
//
// Cookie: "gofin_csrf" (non-httpOnly so JS can read it; Secure in production; SameSite=Lax).
// Header: "X-CSRF-Token".
//
// API-key authenticated requests are exempt (c.Locals("auth_method") == "api_key").
// Safe methods (GET, HEAD, OPTIONS) are exempt by Fiber's CSRF middleware automatically.
func CSRF(cfg CSRFConfig) fiber.Handler {
	return csrf.New(csrf.Config{
		KeyLookup:      "header:X-CSRF-Token",
		CookieName:     "gofin_csrf",
		CookieSameSite: "Lax",
		CookieSecure:   cfg.IsProd,
		CookieHTTPOnly: false,
		CookiePath:     "/",
		Expiration:     1 * time.Hour,
		ContextKey:     "csrf_token",
		// Skip CSRF for API-key requests — they use Bearer token / X-API-Key auth
		// and are not vulnerable to browser-based CSRF attacks.
		Next: func(c *fiber.Ctx) bool {
			if method := c.Method(); method == fiber.MethodGet || method == fiber.MethodHead || method == fiber.MethodOptions {
				return true
			}
			if authMethod, ok := c.Locals("auth_method").(string); ok && authMethod == "api_key" {
				return true
			}
			return false
		},
	})
}

// CSRFTokenHandler returns a handler that responds with the current CSRF token.
// This endpoint is used by the frontend to obtain an initial CSRF token/cookie.
func CSRFTokenHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.Locals("csrf_token").(string)
		if !ok || token == "" {
			// The CSRF middleware should have already set the token in locals.
			// This should not happen if the middleware is properly chained.
			return c.Status(500).JSON(fiber.Map{
				"error": "CSRF token not available",
			})
		}
		return c.JSON(fiber.Map{
			"token": token,
		})
	}
}

package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// apiKeyBlockedGetPrefixes are path prefixes that API keys are never
// allowed to access, even via GET.
var apiKeyBlockedGetPrefixes = []string{
	"/auth/logout",
	"/admin/",
	"/notifications/stream",
	"/metrics",
}

// APIKeyScopeMiddleware restricts API-key-authenticated requests to
// read-only (GET) access and blocks sensitive endpoints entirely.
//
// It must be placed in the middleware chain AFTER APIKeyMiddleware
// (which sets auth_method) so that JWT-authenticated requests are
// left untouched.
func APIKeyScopeMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Only gate requests authenticated via API key.
		if c.Locals("auth_method") != "api_key" {
			return c.Next()
		}

		path := c.Path()

		// Block sensitive GET endpoints even for read access.
		for _, prefix := range apiKeyBlockedGetPrefixes {
			if strings.HasPrefix(path, prefix) {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "API keys do not have access to this endpoint",
				})
			}
		}

		// Allow all GET / HEAD methods; reject everything else.
		method := c.Method()
		if method != fiber.MethodGet && method != fiber.MethodHead {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "API keys only have read access",
			})
		}

		return c.Next()
	}
}

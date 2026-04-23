package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// SecurityHeaders adds common security headers to responses.
func SecurityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		// Prevent MIME type sniffing
		c.Set("X-Content-Type-Options", "nosniff")
		// Prevent clickjacking
		c.Set("X-Frame-Options", "DENY")
		// XSS protection
		c.Set("X-XSS-Protection", "1; mode=block")
		// Referrer policy
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		// HSTS (enable in production with HTTPS)
		c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		// Content Security Policy
		c.Set("Content-Security-Policy", "default-src 'self'")

		return err
	}
}

// RequestSizeLimit rejects requests whose Content-Length header exceeds maxBytes.
// Note: Fiber's own body parser handles the actual body limit; this middleware
// rejects obviously oversized requests early.
func RequestSizeLimit(maxBytes int64) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() == "GET" || c.Method() == "HEAD" || c.Method() == "DELETE" {
			return c.Next()
		}
		cl := int64(c.Request().Header.ContentLength())
		if cl > maxBytes {
			return c.Status(413).JSON(fiber.Map{
				"error": "Request body too large",
			})
		}
		return c.Next()
	}
}

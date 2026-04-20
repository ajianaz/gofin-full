package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

// Logger returns a zerolog-based request logging middleware.
func Logger(logger zerolog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		reqID, _ := c.Locals("request_id").(string)

		err := c.Next()

		duration := time.Since(start)
		status := c.Response().StatusCode()

		evt := logger.Info()
		if status >= 500 {
			evt = logger.Error()
		} else if status >= 400 {
			evt = logger.Warn()
		}

		evt.
			Str("request_id", reqID).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Str("query", string(c.Request().URI().QueryString())).
			Int("status", status).
			Dur("duration", duration).
			Str("client_ip", c.IP()).
			Str("user_agent", c.Get("User-Agent")).
			Msg("request")

		return err
	}
}

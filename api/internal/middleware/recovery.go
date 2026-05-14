package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

// Recovery recovers from panics and returns 500.
// Exception details are only included when APP_DEBUG=true (via DEBUG_PANICS env var).
func Recovery(logger zerolog.Logger) fiber.Handler {
	debugPanic := os.Getenv("APP_DEBUG") == "true"
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				reqID, _ := c.Locals("request_id").(string)
				logger.Error().
					Str("request_id", reqID).
					Str("method", c.Method()).
					Str("path", c.Path()).
					Interface("panic", r).
					Msg("panic recovered")

				response := fiber.Map{
					"message": "Internal Server Error",
				}
				if debugPanic {
					response["exception"] = fmt.Sprintf("%v", r)
				}

				c.Status(500).JSON(response)
			}
		}()
		return c.Next()
	}
}

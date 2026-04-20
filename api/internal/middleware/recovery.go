package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

// Recovery recovers from panics and returns 500.
func Recovery(logger zerolog.Logger) fiber.Handler {
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

				c.Status(500).JSON(fiber.Map{
					"message":   "Internal Server Error",
					"exception": fmt.Sprintf("%v", r),
				})
			}
		}()
		return c.Next()
	}
}

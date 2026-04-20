package middleware

import (
	"errors"
	"reflect"

	"github.com/gofiber/fiber/v2"
	apperrors "github.com/azfirazka/gofin-full/api/pkg/errors"
)

// ErrorHandler is the centralized error handler for all Fiber errors.
// Implements the 5 error response shapes from research/12-api-error-responses.md.
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Handle Fiber native errors
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return c.Status(fiberErr.Code).JSON(apperrors.New(fiberErr.Code, fiberErr.Message))
	}

	// Shape 1: Standard AppError
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		if appErr.StatusCode >= 500 {
			// Don't expose details in production for 500 errors
			return c.Status(appErr.StatusCode).JSON(fiber.Map{
				"message":   appErr.Title,
				"exception": appErr.Exception,
			})
		}
		return c.Status(appErr.StatusCode).JSON(appErr)
	}

	// Shape 2: ValidationError
	var valErr *apperrors.ValidationError
	if errors.As(err, &valErr) {
		return c.Status(422).JSON(valErr)
	}

	// Shape 3: RateLimitError
	var rateErr *apperrors.RateLimitError
	if errors.As(err, &rateErr) {
		c.Set("Retry-After", rateErr.ResetAt)
		return c.Status(429).JSON(fiber.Map{
			"error": rateErr,
		})
	}

	// Shape 5: Internal error fallback
	return c.Status(500).JSON(fiber.Map{
		"message":   "Internal Server Error",
		"exception": reflect.TypeOf(err).String(),
	})
}

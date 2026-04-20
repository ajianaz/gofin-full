package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	apperrors "github.com/azfirazka/gofin-full/api/pkg/errors"
)

var validAcceptHeaders = []string{
	"application/json",
	"application/vnd.api+json",
	"application/x-www-form-urlencoded",
	"application/octet-stream",
	"*/*",
}

var validContentTypes = []string{
	"application/json",
	"application/vnd.api+json",
	"application/x-www-form-urlencoded",
	"application/octet-stream",
}

// AcceptHeaders validates Accept and Content-Type headers per Firefly III spec.
func AcceptHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Validate Accept header
		accept := c.Get("Accept")
		if accept != "" && !isValidHeader(accept, validAcceptHeaders) {
			return apperrors.NotAcceptableHeader(accept)
		}

		// Validate Content-Type for POST/PUT/PATCH
		method := c.Method()
		if method == "POST" || method == "PUT" || method == "PATCH" {
			contentType := c.Get("Content-Type")
			if contentType == "" {
				return apperrors.EmptyContentType()
			}
			// Strip boundary for multipart
			ct := strings.Split(contentType, ";")[0]
			ct = strings.TrimSpace(ct)
			if !isValidHeader(ct, validContentTypes) {
				return apperrors.UnsupportedContentType(ct)
			}
		}

		return c.Next()
	}
}

func isValidHeader(value string, valid []string) bool {
	for _, v := range valid {
		if strings.EqualFold(value, v) {
			return true
		}
	}
	return false
}

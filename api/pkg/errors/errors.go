package errors

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

// AppError is the base application error with HTTP status code.
type AppError struct {
	StatusCode int    `json:"status"`
	Title      string `json:"title"`
	Detail     string `json:"detail,omitempty"`
	Exception  string `json:"exception,omitempty"`
}

func (e *AppError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("%s: %s", e.Title, e.Detail)
	}
	return e.Title
}

func (e *AppError) Unwrap() error {
	return nil
}

// New creates a new AppError.
func New(statusCode int, title string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Title:      title,
		Exception:  title,
	}
}

// NewWithDetail creates a new AppError with detail message.
func NewWithDetail(statusCode int, title, detail string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Title:      title,
		Detail:     detail,
		Exception:  title,
	}
}

// ValidationError represents field-level validation errors (Shape 2).
type ValidationError struct {
	Message    string              `json:"message"`
	FieldErrors map[string][]string `json:"errors"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

// NewValidationError creates a new validation error.
func NewValidationError(fieldErrors map[string][]string) *ValidationError {
	return &ValidationError{
		Message:    "The given data was invalid.",
		FieldErrors: fieldErrors,
	}
}

// RateLimitError represents rate limit exceeded (Shape 3).
type RateLimitError struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Limit     int    `json:"limit"`
	Remaining int    `json:"remaining"`
	ResetAt   string `json:"reset_at"`
}

func (e *RateLimitError) Error() string {
	return e.Message
}

// Predefined errors.
var (
	ErrBadRequest       = New(http.StatusBadRequest, "Bad Request")
	ErrUnauthorized     = New(http.StatusUnauthorized, "Unauthenticated")
	ErrForbidden        = New(http.StatusForbidden, "Forbidden")
	ErrNotFound         = New(http.StatusNotFound, "Resource not found")
	ErrMethodNotAllowed = New(http.StatusMethodNotAllowed, "Method Not Allowed")
	ErrNotAcceptable    = New(http.StatusNotAcceptable, "Not Acceptable")
	ErrConflict         = New(http.StatusConflict, "Conflict")
	ErrGone             = New(http.StatusGone, "Gone")
	ErrUnsupportedMedia = New(http.StatusUnsupportedMediaType, "Unsupported Media Type")
	ErrInternal         = New(http.StatusInternalServerError, "Internal Server Error")
	ErrTooManyRequests  = New(http.StatusTooManyRequests, "Too Many Requests")
)

// NotFoundResource creates a not-found error for a specific resource type.
func NotFoundResource(resourceType string, id uuid.UUID) *AppError {
	return &AppError{
		StatusCode: http.StatusNotFound,
		Title:      "Resource not found",
		Detail:     fmt.Sprintf("%s with ID %s not found", resourceType, id),
		Exception:  "NotFoundHttpException",
	}
}

// GoneTransaction creates a 410 Gone error for rule-deleted transactions.
func GoneTransaction() *AppError {
	return &AppError{
		StatusCode: http.StatusGone,
		Title:      "Gone",
		Detail:     "200032: Cannot find transaction. Possibly, a rule deleted this transaction after its creation.",
		Exception:  "GoneException",
	}
}

// DemoUserBlocked creates a 403 error for demo user protection.
func DemoUserBlocked() *AppError {
	return &AppError{
		StatusCode: http.StatusForbidden,
		Title:      "Forbidden",
		Detail:     "Demo user is not allowed to perform this action.",
		Exception:  "DemoUserException",
	}
}

// NotAcceptableHeader creates a 406 error for invalid Accept header.
func NotAcceptableHeader(header string) *AppError {
	return &AppError{
		StatusCode: http.StatusNotAcceptable,
		Title:      "Not Acceptable",
		Detail:     fmt.Sprintf("Accept header \"%s\" is not something this server can provide.", header),
		Exception:  "NotAcceptableException",
	}
}

// UnsupportedContentType creates a 415 error for invalid Content-Type.
func UnsupportedContentType(contentType string) *AppError {
	return &AppError{
		StatusCode: http.StatusUnsupportedMediaType,
		Title:      "Unsupported Media Type",
		Detail:     fmt.Sprintf("Content-Type cannot be \"%s\".", contentType),
		Exception:  "UnsupportedMediaTypeException",
	}
}

// EmptyContentType creates a 415 error for missing Content-Type.
func EmptyContentType() *AppError {
	return &AppError{
		StatusCode: http.StatusUnsupportedMediaType,
		Title:      "Unsupported Media Type",
		Detail:     "Content-Type header cannot be empty.",
		Exception:  "UnsupportedMediaTypeException",
	}
}

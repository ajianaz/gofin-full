package validation

import (
	"fmt"
	"strings"
	"unicode"

	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors"
)

// FieldErrors is a convenience type for building validation error maps.
type FieldErrors map[string][]string

// Add appends an error message for a field.
func (f FieldErrors) Add(field, msg string) {
	f[field] = append(f[field], msg)
}

// Has returns true if any errors exist.
func (f FieldErrors) Has() bool {
	return len(f) > 0
}

// ToAppError converts to the app-level ValidationError.
func (f FieldErrors) ToAppError() *apperrors.ValidationError {
	if len(f) == 0 {
		return nil
	}
	return apperrors.NewValidationError(f)
}

// --- Reusable validation rules ---

// Required checks that a string is non-empty after trimming whitespace.
func Required(field, value string, errs FieldErrors) {
	if strings.TrimSpace(value) == "" {
		errs.Add(field, field+" is required.")
	}
}

// MinLength checks minimum string length.
func MinLength(field, value string, min int, errs FieldErrors) {
	if len(value) < min {
		errs.Add(field, fmt.Sprintf("%s must be at least %d characters.", field, min))
	}
}

// MaxLength checks maximum string length.
func MaxLength(field, value string, max int, errs FieldErrors) {
	if len(value) > max {
		errs.Add(field, fmt.Sprintf("%s must be at most %d characters.", field, max))
	}
}

// Email checks basic email format (at-sign + domain with dot).
func Email(field, value string, errs FieldErrors) {
	if value == "" {
		return
	}
	at := strings.LastIndex(value, "@")
	if at < 1 || at == len(value)-1 {
		errs.Add(field, "Invalid email address.")
		return
	}
	dot := strings.LastIndex(value[at:], ".")
	if dot < 2 {
		errs.Add(field, "Invalid email address.")
	}
}

// OneOf checks if a value is one of the allowed values.
func OneOf(field, value string, allowed []string, errs FieldErrors) {
	for _, a := range allowed {
		if value == a {
			return
		}
	}
	errs.Add(field, "Must be one of: "+strings.Join(allowed, ", "))
}

// PasswordStrength checks common password requirements.
// Skips validation if value is empty (use Required separately).
func PasswordStrength(field, value string, errs FieldErrors) {
	if value == "" {
		return
	}
	var issues []string
	if len(value) < 8 {
		issues = append(issues, "at least 8 characters")
	}
	var hasUpper, hasLower, hasDigit bool
	for _, c := range value {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasDigit = true
		}
	}
	if !hasLower {
		issues = append(issues, "a lowercase letter")
	}
	if !hasUpper {
		issues = append(issues, "an uppercase letter")
	}
	if !hasDigit {
		issues = append(issues, "a digit")
	}
	if len(issues) > 0 {
		errs.Add(field, "Password must contain "+strings.Join(issues, ", "))
	}
}

// IntRange checks that an integer is within [min, max].
func IntRange(field string, value, min, max int, errs FieldErrors) {
	if value < min || value > max {
		errs.Add(field, fmt.Sprintf("%s must be between %d and %d.", field, min, max))
	}
}

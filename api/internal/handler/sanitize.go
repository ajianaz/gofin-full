package handler

import "html"

// sanitizeStr escapes HTML entities in user-provided strings to prevent XSS.
// Apply this to all user-supplied text fields (names, titles, descriptions, notes)
// before persisting to the database or returning in API responses.
func sanitizeStr(s string) string {
	return html.EscapeString(s)
}

// sanitizePtr escapes HTML entities in a pointer string, returning a new pointer.
// Returns nil if the input is nil.
func sanitizePtr(s *string) *string {
	if s == nil {
		return nil
	}
	escaped := html.EscapeString(*s)
	return &escaped
}

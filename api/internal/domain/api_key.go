package domain

import "time"

// APIKey represents a long-lived API key for programmatic access.
type APIKey struct {
	ID          int64
	UserID      int64
	Name        string
	KeyHash     string
	KeyPrefix   string
	LastUsedAt  *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

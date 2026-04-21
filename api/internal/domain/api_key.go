package domain

import (
	"time"

	"github.com/google/uuid"
)

// APIKey represents a long-lived API key for programmatic access.
type APIKey struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        string
	KeyHash     string
	KeyPrefix   string
	LastUsedAt  *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

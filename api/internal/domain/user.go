package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	Email        string     `json:"email" db:"email"`
	Password     string     `json:"-" db:"password"`
	Blocked      bool       `json:"blocked" db:"blocked"`
	BlockedCode  *string    `json:"blocked_code,omitempty" db:"blocked_code"`
	DemoUser     bool       `json:"demo_user" db:"demo_user"`
	UserGroupID  *uuid.UUID `json:"user_group_id,omitempty" db:"user_group_id"`
	TokenVersion int        `json:"token_version" db:"token_version"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time `json:"-" db:"deleted_at"`

	// Joined (not in DB)
	Group *UserGroup `json:"group,omitempty" db:"-"`
	Roles []Role     `json:"roles,omitempty" db:"-"`
}

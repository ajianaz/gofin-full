package domain

import "time"

type Category struct {
	ID         int64     `json:"id" db:"id"`
	UserID     int64     `json:"user_id" db:"user_id"`
	UserGroupID int64    `json:"user_group_id" db:"user_group_id"`
	Name       string    `json:"name" db:"name"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt  *time.Time `json:"-" db:"deleted_at"`
}

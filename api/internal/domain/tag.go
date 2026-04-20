package domain

import "time"

type Tag struct {
	ID         int64     `json:"id" db:"id"`
	UserID     int64     `json:"user_id" db:"user_id"`
	UserGroupID int64    `json:"user_group_id" db:"user_group_id"`
	Tag        string    `json:"tag" db:"tag"`
	Date       *time.Time `json:"date,omitempty" db:"date"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt  *time.Time `json:"-" db:"deleted_at"`
}

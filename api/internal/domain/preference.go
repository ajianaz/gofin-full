package domain

import "time"

type Preference struct {
	ID         int64     `json:"id" db:"id"`
	UserID     int64     `json:"user_id" db:"user_id"`
	Name       string    `json:"name" db:"name"`
	Data       string    `json:"data" db:"data"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type Configuration struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Value     string    `json:"value" db:"value"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type ObjectGroup struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	UserGroupID int64   `json:"user_group_id" db:"user_group_id"`
	Title     string    `json:"title" db:"title"`
	Order     int       `json:"order" db:"order"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Note struct {
	ID             int64     `json:"id" db:"id"`
	NoteableType   string    `json:"noteable_type" db:"noteable_type"`
	NoteableID     int64     `json:"noteable_id" db:"noteable_id"`
	Note           string    `json:"note" db:"note"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type Location struct {
	ID            int64     `json:"id" db:"id"`
	LocatableType string    `json:"locatable_type" db:"locatable_type"`
	LocatableID   int64     `json:"locatable_id" db:"locatable_id"`
	Latitude      *float64  `json:"latitude,omitempty" db:"latitude"`
	Longitude     *float64  `json:"longitude,omitempty" db:"longitude"`
	ZoomLevel     *int      `json:"zoom_level,omitempty" db:"zoom_level"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

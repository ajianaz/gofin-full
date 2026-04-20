package domain

import "time"

type Currency struct {
	ID             int64     `json:"id" db:"id"`
	Code           string    `json:"code" db:"code"`
	Name           string    `json:"name" db:"name"`
	Symbol         string    `json:"symbol" db:"symbol"`
	DecimalPlaces  int       `json:"decimal_places" db:"decimal_places"`
	Enabled        bool      `json:"enabled" db:"enabled"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt      *time.Time `json:"-" db:"deleted_at"`
}

type AccountType struct {
	ID        int64     `json:"id" db:"id"`
	Type      string    `json:"type" db:"type"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

package domain

import (
	"github.com/shopspring/decimal"
	"time"
)

type Bill struct {
	ID                 int64           `json:"id" db:"id"`
	UserID             int64           `json:"user_id" db:"user_id"`
	UserGroupID        int64           `json:"user_group_id" db:"user_group_id"`
	Name               string          `json:"name" db:"name"`
	AmountMin          decimal.Decimal `json:"amount_min" db:"amount_min"`
	AmountMax          decimal.Decimal `json:"amount_max" db:"amount_max"`
	Date               time.Time       `json:"date" db:"date"`
	EndDate            *time.Time      `json:"end_date,omitempty" db:"end_date"`
	RepeatFreq         string          `json:"repeat_freq" db:"repeat_freq"`
	Skip               int             `json:"skip" db:"skip"`
	Active             bool            `json:"active" db:"active"`
	Order              int             `json:"order" db:"order"`
	Notes              *string         `json:"notes,omitempty" db:"notes"`
	CurrencyID         string          `json:"currency_id" db:"currency_id"`
	CreatedAt          time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at" db:"updated_at"`
	DeletedAt          *time.Time      `json:"-" db:"deleted_at"`

	// Joined
	ObjectGroups []ObjectGroup `json:"object_groups,omitempty" db:"-"`
}

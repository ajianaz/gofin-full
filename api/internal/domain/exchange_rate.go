package domain

import (
	"github.com/shopspring/decimal"
	"time"
)

type ExchangeRate struct {
	ID             int64           `json:"id" db:"id"`
	UserID         int64           `json:"user_id" db:"user_id"`
	UserGroupID    int64           `json:"user_group_id" db:"user_group_id"`
	FromCurrencyID string          `json:"from_currency_id" db:"from_currency_id"`
	ToCurrencyID   string          `json:"to_currency_id" db:"to_currency_id"`
	Rate           decimal.Decimal `json:"rate" db:"rate"`
	Date           time.Time       `json:"date" db:"date"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at" db:"updated_at"`
}

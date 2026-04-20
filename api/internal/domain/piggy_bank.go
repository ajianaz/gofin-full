package domain

import (
	"github.com/shopspring/decimal"
	"time"
)

type PiggyBank struct {
	ID                  int64           `json:"id" db:"id"`
	AccountID           int64           `json:"wallet_id" db:"account_id"`
	Name                string          `json:"name" db:"name"`
	TargetAmount        decimal.Decimal `json:"target_amount" db:"target_amount"`
	StartDate           *time.Time      `json:"start_date,omitempty" db:"start_date"`
	TargetDate          *time.Time      `json:"target_date,omitempty" db:"target_date"`
	Order               int             `json:"order" db:"order"`
	Notes               *string         `json:"notes,omitempty" db:"notes"`
	CreatedAt           time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at" db:"updated_at"`
	DeletedAt           *time.Time      `json:"-" db:"deleted_at"`

	// Computed
	CurrentAmount       decimal.Decimal `json:"current_amount" db:"-"`
	NativeCurrentAmount decimal.Decimal `json:"native_current_amount" db:"-"`
	LeftToTarget        decimal.Decimal `json:"left_to_target" db:"-"`
	Percentage          float64         `json:"percentage" db:"-"`

	// Joined
	Account             *Wallet         `json:"wallet,omitempty" db:"-"`
}

type PiggyBankEvent struct {
	ID          int64           `json:"id" db:"id"`
	PiggyBankID int64           `json:"piggy_bank_id" db:"piggy_bank_id"`
	Amount      decimal.Decimal `json:"amount" db:"amount"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}

type PiggyBankRepetition struct {
	ID              int64           `json:"id" db:"id"`
	PiggyBankID     int64           `json:"piggy_bank_id" db:"piggy_bank_id"`
	TargetAmount    decimal.Decimal `json:"target_amount" db:"target_amount"`
	CurrentAmount   decimal.Decimal `json:"current_amount" db:"current_amount"`
	StartDate       time.Time       `json:"start_date" db:"start_date"`
	TargetDate      *time.Time      `json:"target_date,omitempty" db:"target_date"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
}

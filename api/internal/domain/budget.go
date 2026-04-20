package domain

import (
	"github.com/shopspring/decimal"
	"time"
)

type AutoBudgetType string

const (
	AutoBudgetTypeNone      AutoBudgetType = "none"
	AutoBudgetTypeReset     AutoBudgetType = "reset"
	AutoBudgetTypeRollover  AutoBudgetType = "rollover"
	AutoBudgetTypeAdjusted  AutoBudgetType = "adjusted"
)

type Budget struct {
	ID         int64     `json:"id" db:"id"`
	UserID     int64     `json:"user_id" db:"user_id"`
	UserGroupID int64    `json:"user_group_id" db:"user_group_id"`
	Name       string    `json:"name" db:"name"`
	Active     bool      `json:"active" db:"active"`
	Order      int       `json:"order" db:"order"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt  *time.Time `json:"-" db:"deleted_at"`

	// Joined
	Limits    []BudgetLimit   `json:"limits,omitempty" db:"-"`
	AutoBudget *AutoBudget    `json:"auto_budget,omitempty" db:"-"`
}

type BudgetLimit struct {
	ID          int64           `json:"id" db:"id"`
	BudgetID    int64           `json:"budget_id" db:"budget_id"`
	Start       time.Time       `json:"start" db:"start"`
	End         time.Time       `json:"end" db:"end"`
	Amount      decimal.Decimal `json:"amount" db:"amount"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}

type AutoBudget struct {
	ID               int64          `json:"id" db:"id"`
	BudgetID         int64          `json:"budget_id" db:"budget_id"`
	AutoBudgetType   AutoBudgetType `json:"auto_budget_type" db:"auto_budget_type"`
	Period           string         `json:"period" db:"period"`
	CreatedAt        time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at" db:"updated_at"`
}

type AvailableBudget struct {
	ID           int64           `json:"id" db:"id"`
	BudgetID     int64           `json:"budget_id" db:"budget_id"`
	CurrencyID   int64           `json:"currency_id" db:"currency_id"`
	Start        time.Time       `json:"start" db:"start"`
	End          time.Time       `json:"end" db:"end"`
	Amount       decimal.Decimal `json:"amount" db:"amount"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at" db:"updated_at"`
}

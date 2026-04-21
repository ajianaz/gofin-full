package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AutoBudgetType string

const (
	AutoBudgetTypeNone      AutoBudgetType = "none"
	AutoBudgetTypeReset     AutoBudgetType = "reset"
	AutoBudgetTypeRollover  AutoBudgetType = "rollover"
	AutoBudgetTypeAdjusted  AutoBudgetType = "adjusted"
)

type Budget struct {
	ID         uuid.UUID       `json:"id" db:"id"`
	UserID     uuid.UUID       `json:"user_id" db:"user_id"`
	UserGroupID uuid.UUID      `json:"user_group_id" db:"user_group_id"`
	Name       string          `json:"name" db:"name"`
	Active     bool            `json:"active" db:"active"`
	Order      int             `json:"order" db:"order"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at" db:"updated_at"`
	DeletedAt  *time.Time      `json:"-" db:"deleted_at"`

	// Joined
	Limits    []BudgetLimit   `json:"limits,omitempty" db:"-"`
	AutoBudget *AutoBudget    `json:"auto_budget,omitempty" db:"-"`
}

type BudgetLimit struct {
	ID          uuid.UUID       `json:"id" db:"id"`
	BudgetID    uuid.UUID       `json:"budget_id" db:"budget_id"`
	Start       time.Time       `json:"start" db:"start"`
	End         time.Time       `json:"end" db:"end"`
	Amount      decimal.Decimal `json:"amount" db:"amount"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}

type AutoBudget struct {
	ID               uuid.UUID       `json:"id" db:"id"`
	BudgetID         uuid.UUID       `json:"budget_id" db:"budget_id"`
	AutoBudgetType   AutoBudgetType  `json:"auto_budget_type" db:"auto_budget_type"`
	Period           string          `json:"period" db:"period"`
	CreatedAt        time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at" db:"updated_at"`
}

type AvailableBudget struct {
	ID           uuid.UUID       `json:"id" db:"id"`
	BudgetID     uuid.UUID       `json:"budget_id" db:"budget_id"`
	CurrencyID   uuid.UUID       `json:"currency_id" db:"currency_id"`
	Start        time.Time       `json:"start" db:"start"`
	End          time.Time       `json:"end" db:"end"`
	Amount       decimal.Decimal `json:"amount" db:"amount"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at" db:"updated_at"`
}

package domain

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

// Recurrence represents a recurring transaction schedule (maps to recurrences table).
type Recurrence struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	UserID             uuid.UUID `json:"user_id" db:"user_id"`
	UserGroupID        uuid.UUID `json:"user_group_id" db:"user_group_id"`
	Title              string    `json:"title" db:"title"`
	Description        *string   `json:"description,omitempty" db:"description"`
	FirstDate          time.Time `json:"first_date" db:"first_date"`
	LatestDate         *time.Time `json:"latest_date,omitempty" db:"latest_date"`
	RepeatUntil        *time.Time `json:"repeat_until,omitempty" db:"repeat_until"`
	RepeatFreq         string    `json:"repeat_freq" db:"repeat_freq"`
	Skip               int       `json:"skip" db:"skip"`
	Active             bool      `json:"active" db:"active"`
	ApplyRules         bool      `json:"apply_rules" db:"apply_rules"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt          *time.Time `json:"-" db:"deleted_at"`

	// Joined
	Transactions []RecurringTransaction `json:"transactions,omitempty" db:"-"`
	Repetitions  []RecurringRepetition `json:"repetitions,omitempty" db:"-"`
	Meta         []RecurrenceMeta      `json:"meta,omitempty" db:"-"`
}

// RecurringTransaction represents a transaction template within a recurrence.
type RecurringTransaction struct {
	ID             uuid.UUID       `json:"id" db:"id"`
	RecurrenceID   uuid.UUID       `json:"recurrence_id" db:"recurrence_id"`
	Type           string          `json:"type" db:"type"`
	Description    string          `json:"description" db:"description"`
	Amount         decimal.Decimal `json:"amount" db:"amount"`
	CurrencyID     uuid.UUID       `json:"currency_id" db:"transaction_currency_id"`
	SourceID       uuid.UUID       `json:"source_id" db:"source_id"`
	DestinationID  uuid.UUID       `json:"destination_id" db:"destination_id"`
	BudgetID       *uuid.UUID      `json:"budget_id,omitempty" db:"budget_id"`
	CategoryID     *uuid.UUID      `json:"category_id,omitempty" db:"category_id"`
	PiggyBankID    *uuid.UUID      `json:"piggy_bank_id,omitempty" db:"piggy_bank_id"`
	Order          int             `json:"order" db:"order"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at" db:"updated_at"`

	// Joined
	Meta []RecurringTransactionMeta `json:"meta,omitempty" db:"-"`
}

type RecurringRepetition struct {
	ID           uuid.UUID `json:"id" db:"id"`
	RecurrenceID uuid.UUID `json:"recurrence_id" db:"recurrence_id"`
	RelevantDate time.Time `json:"relevant_date" db:"relevant_date"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type RecurrenceMeta struct {
	ID           uuid.UUID `json:"id" db:"id"`
	RecurrenceID uuid.UUID `json:"recurrence_id" db:"recurrence_id"`
	Name         string    `json:"name" db:"name"`
	Value        string    `json:"value" db:"value"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type RecurringTransactionMeta struct {
	ID                     uuid.UUID `json:"id" db:"id"`
	RecurringTransactionID uuid.UUID `json:"recurring_transaction_id" db:"recurring_transaction_id"`
	Name                   string    `json:"name" db:"name"`
	Value                  string    `json:"value" db:"value"`
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time `json:"updated_at" db:"updated_at"`
}

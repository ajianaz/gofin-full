package domain

import (
	"github.com/shopspring/decimal"
	"time"
)

// TransactionType represents the 7 transaction types from Firefly III.
type TransactionType string

const (
	TransactionTypeWithdrawal      TransactionType = "withdrawal"
	TransactionTypeDeposit         TransactionType = "deposit"
	TransactionTypeTransfer        TransactionType = "transfer"
	TransactionTypeOpeningBalance  TransactionType = "opening-balance"
	TransactionTypeReconciliation  TransactionType = "reconciliation"
	TransactionTypeLiabilityCredit TransactionType = "liability-credit"
	TransactionTypeInvalid         TransactionType = "invalid"
)

// TransactionGroup is the top-level transaction container.
// Groups one or more TransactionJournals (split transactions).
type TransactionGroup struct {
	ID         int64                `json:"id" db:"id"`
	UserID     int64                `json:"user_id" db:"user_id"`
	UserGroupID int64               `json:"user_group_id" db:"user_group_id"`
	GroupTitle string               `json:"group_title" db:"group_title"`
	CreatedAt  time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time            `json:"updated_at" db:"updated_at"`
	DeletedAt  *time.Time           `json:"-" db:"deleted_at"`
	Journals   []TransactionJournal `json:"transactions,omitempty" db:"-"`
}

// TransactionJournal represents a single journal within a transaction group.
type TransactionJournal struct {
	ID                     int64                     `json:"transaction_journal_id" db:"id"`
	TransactionGroupID     int64                     `json:"-" db:"transaction_group_id"`
	UserID                 int64                     `json:"user" db:"user_id"`
	UserGroupID            int64                     `json:"user_group" db:"user_group_id"`
	TransactionTypeID      int64                     `json:"-" db:"transaction_type_id"`
	Type                   TransactionType           `json:"type" db:"-"`
	Date                   time.Time                 `json:"date" db:"date"`
	Order                  int                       `json:"order" db:"order"`
	Description            string                    `json:"description" db:"description"`
	CurrencyID             string                    `json:"currency_id" db:"transaction_currency_id"`
	ForeignCurrencyID      *string                   `json:"foreign_currency_id,omitempty" db:"foreign_currency_id"`
	BudgetID               *int64                    `json:"budget_id,omitempty" db:"budget_id"`
	BillID                 *int64                    `json:"bill_id,omitempty" db:"bill_id"`
	PiggyBankID            *int64                    `json:"piggy_bank_id,omitempty" db:"piggy_bank_id"`
	Reconciled             bool                      `json:"reconciled" db:"reconciled"`
	Notes                  *string                   `json:"notes,omitempty" db:"notes"`
	InterestDate           *time.Time                `json:"interest_date,omitempty" db:"interest_date"`
	BookDate               *time.Time                `json:"book_date,omitempty" db:"book_date"`
	ProcessDate            *time.Time                `json:"process_date,omitempty" db:"process_date"`
	DueDate                *time.Time                `json:"due_date,omitempty" db:"due_date"`
	PaymentDate            *time.Time                `json:"payment_date,omitempty" db:"payment_date"`
	InvoiceDate            *time.Time                `json:"invoice_date,omitempty" db:"invoice_date"`
	ExternalID             *string                   `json:"external_id,omitempty" db:"external_id"`
	ExternalURL            *string                   `json:"external_url,omitempty" db:"external_url"`
	InternalReference      *string                   `json:"internal_reference,omitempty" db:"internal_reference"`
	RecurrenceID           *int64                    `json:"recurrence_id,omitempty" db:"recurrence_id"`
	RecurrenceTotal        *int                      `json:"recurrence_total,omitempty" db:"recurrence_total"`
	RecurrenceCount        *int                      `json:"recurrence_count,omitempty" db:"recurrence_count"`
	ImportHashV2           *string                   `json:"import_hash_v2,omitempty" db:"import_hash_v2"`
	OriginalSource         *string                   `json:"original_source,omitempty" db:"original_source"`
	BalanceDirty           bool                      `json:"-" db:"-"`
	HasAttachments         bool                      `json:"has_attachments" db:"-"`
	CreatedAt              time.Time                 `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time                 `json:"updated_at" db:"updated_at"`
	DeletedAt              *time.Time                `json:"-" db:"deleted_at"`

	// SEPA fields
	SepaCC       *string `json:"sepa_cc,omitempty" db:"sepa_cc"`
	SepaCTOp     *string `json:"sepa_ct_op,omitempty" db:"sepa_ct_op"`
	SepaCTID     *string `json:"sepa_ct_id,omitempty" db:"sepa_ct_id"`
	SepaDB       *string `json:"sepa_db,omitempty" db:"sepa_db"`
	SepaCountry  *string `json:"sepa_country,omitempty" db:"sepa_country"`
	SepaEP       *string `json:"sepa_ep,omitempty" db:"sepa_ep"`
	SepaCI       *string `json:"sepa_ci,omitempty" db:"sepa_ci"`
	SepaBatchID  *string `json:"sepa_batch_id,omitempty" db:"sepa_batch_id"`

	// Joined (not in DB directly)
	Tags                    []Tag              `json:"tags,omitempty" db:"-"`
	Categories              []Category         `json:"categories,omitempty" db:"-"`
	Source                  *Wallet            `json:"source,omitempty" db:"-"`
	Destination             *Wallet            `json:"destination,omitempty" db:"-"`
	Currency                *Currency          `json:"currency,omitempty" db:"-"`
	ForeignCurrency         *Currency          `json:"foreign_currency,omitempty" db:"-"`
	SourceTransactions      []Transaction      `json:"-" db:"-"`
	DestinationTransactions []Transaction      `json:"-" db:"-"`
}

// Transaction represents the actual monetary movement (debit or credit).
// Each journal has exactly 2 transactions: one debit, one credit.
type Transaction struct {
	ID                     int64           `json:"id" db:"id"`
	TransactionJournalID   int64           `json:"transaction_journal_id" db:"transaction_journal_id"`
	AccountID              int64           `json:"-" db:"account_id"`
	Amount                 decimal.Decimal `json:"amount" db:"amount"`
	NativeAmount           decimal.Decimal `json:"native_amount,omitempty" db:"native_amount"`
	ForeignAmount          *decimal.Decimal `json:"foreign_amount,omitempty" db:"foreign_amount"`
	NativeForeignAmount    *decimal.Decimal `json:"pc_foreign_amount,omitempty" db:"native_foreign_amount"`
	ForeignCurrencyID      *string         `json:"foreign_currency_id,omitempty" db:"foreign_currency_id"`
	Reconciled             bool            `json:"reconciled" db:"reconciled"`
	CreatedAt              time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time       `json:"updated_at" db:"updated_at"`

	// Joined
	Account *Wallet `json:"-" db:"-"`
}

// TransactionJournalMeta stores key-value metadata for journals.
type TransactionJournalMeta struct {
	ID                    int64     `json:"id" db:"id"`
	TransactionJournalID  int64     `json:"-" db:"transaction_journal_id"`
	Name                  string    `json:"name" db:"name"`
	Value                 string    `json:"value" db:"value"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
}

// TransactionJournalLink represents a link between two journals.
type TransactionJournalLink struct {
	ID                  int64     `json:"id" db:"id"`
	LinkTypeID          int64     `json:"link_type_id" db:"link_type_id"`
	SourceID            int64     `json:"source_id" db:"source_id"`
	DestinationID       int64     `json:"destination_id" db:"destination_id"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

// LinkType defines the type of journal link.
type LinkType struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Inward    string    `json:"inward" db:"inward"`
	Outward   string    `json:"outward" db:"outward"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

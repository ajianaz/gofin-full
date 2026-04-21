package domain

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

// WalletType represents the 14 wallet types from Firefly III.
type WalletType string

const (
	WalletTypeAsset           WalletType = "asset"
	WalletTypeDefault         WalletType = "defaultAsset"
	WalletTypeBeneficiary     WalletType = "beneficiary"
	WalletTypeCash            WalletType = "cash"
	WalletTypeCreditCard      WalletType = "credit-card"
	WalletTypeDebt            WalletType = "debt"
	WalletTypeExpense         WalletType = "expense"
	WalletTypeImport          WalletType = "import"
	WalletTypeInitialBalance  WalletType = "initial-balance"
	WalletTypeLiabilityCredit WalletType = "liability-credit"
	WalletTypeLoan            WalletType = "loan"
	WalletTypeMortgage        WalletType = "mortgage"
	WalletTypeReconciliation  WalletType = "reconciliation"
	WalletTypeRevenue         WalletType = "revenue"
)

// WalletRole represents the wallet role for asset accounts.
type WalletRole string

const (
	WalletRoleDefaultAsset  WalletRole = "defaultAsset"
	WalletRoleSharedAsset   WalletRole = "sharedAsset"
	WalletRoleSavingAsset   WalletRole = "savingAsset"
	WalletRoleCcAsset       WalletRole = "ccAsset"
	WalletRoleCashWallet    WalletRole = "cashWalletAsset"
)

// CreditCardType represents the type of credit card.
type CreditCardType string

const (
	CreditCardTypeMonthlyFull CreditCardType = "monthlyFull"
)

// LiabilityDirection represents the direction of a liability.
type LiabilityDirection string

const (
	LiabilityDirectionCredit LiabilityDirection = "credit"
	LiabilityDirectionDebit  LiabilityDirection = "debit"
)

// LiabilityType represents the type of liability.
type LiabilityType string

const (
	LiabilityTypeLoan     LiabilityType = "loan"
	LiabilityTypeDebt     LiabilityType = "debt"
	LiabilityTypeMortgage LiabilityType = "mortgage"
)

// Wallet (financial wallet) represents a financial account.
type Wallet struct {
	ID                   uuid.UUID          `json:"id" db:"id"`
	UserID               uuid.UUID          `json:"user_id" db:"user_id"`
	UserGroupID          uuid.UUID          `json:"user_group_id" db:"user_group_id"`
	AccountType          string             `json:"wallet_type" db:"account_type"`
	Name                 string             `json:"name" db:"name"`
	Active               bool               `json:"active" db:"active"`
	CurrencyID           *string            `json:"currency_id,omitempty" db:"currency_id"`
	VirtualBalance       decimal.Decimal    `json:"virtual_balance" db:"virtual_balance"`
	IBAN                 *string            `json:"iban,omitempty" db:"iban"`
	BIC                  *string            `json:"bic,omitempty" db:"bic"`
	IncludeNetWorth      bool               `json:"include_net_worth" db:"include_net_worth"`
	Notes                *string            `json:"notes,omitempty" db:"notes"`
	Latitude             *float64           `json:"latitude,omitempty" db:"latitude"`
	Longitude            *float64           `json:"longitude,omitempty" db:"longitude"`
	CreatedAt            time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time          `json:"updated_at" db:"updated_at"`
	DeletedAt            *time.Time         `json:"-" db:"deleted_at"`

	// Liability fields
	LiabilityType      *string          `json:"liability_type,omitempty" db:"liability_type"`
	LiabilityDirection *string          `json:"liability_direction,omitempty" db:"liability_direction"`
	InterestRate       *decimal.Decimal `json:"interest,omitempty" db:"interest"`
	InterestPeriod     *string          `json:"interest_period,omitempty" db:"interest_period"`
	CurrentDebt        *decimal.Decimal `json:"current_debt,omitempty" db:"current_debt"`

	// Credit card fields
	CreditCardType       *string     `json:"credit_card_type,omitempty" db:"credit_card_type"`
	MonthlyPaymentDate   *time.Time  `json:"monthly_payment_date,omitempty" db:"monthly_payment_date"`
	MonthlyPaymentAmount *decimal.Decimal `json:"monthly_payment_amount,omitempty" db:"monthly_payment_amount"`
}

// IsSourceValid returns true if the wallet type can be a transaction source.
func IsSourceValid(wt WalletType) bool {
	switch wt {
	case WalletTypeAsset, WalletTypeDefault, WalletTypeCash, WalletTypeDebt,
		WalletTypeInitialBalance, WalletTypeLoan, WalletTypeMortgage, WalletTypeReconciliation:
		return true
	default:
		return false
	}
}

// IsDestinationValid returns true if the wallet type can be a transaction destination.
func IsDestinationValid(wt WalletType) bool {
	switch wt {
	case WalletTypeAsset, WalletTypeDefault, WalletTypeCash, WalletTypeDebt,
		WalletTypeExpense, WalletTypeInitialBalance, WalletTypeLoan, WalletTypeMortgage:
		return true
	default:
		return false
	}
}

// CanHoldPiggyBanks returns true if the wallet can have piggy banks.
func CanHoldPiggyBanks(wt WalletType) bool {
	switch wt {
	case WalletTypeAsset, WalletTypeDefault, WalletTypeLoan, WalletTypeMortgage, WalletTypeDebt:
		return true
	default:
		return false
	}
}

// CanHaveOpeningBalance returns true if the wallet can have an opening balance.
func CanHaveOpeningBalance(wt WalletType) bool {
	switch wt {
	case WalletTypeAsset, WalletTypeDefault, WalletTypeLoan, WalletTypeMortgage, WalletTypeDebt:
		return true
	default:
		return false
	}
}

// CanHaveCurrency returns true if the wallet can have its own currency setting.
func CanHaveCurrency(wt WalletType) bool {
	switch wt {
	case WalletTypeAsset, WalletTypeDefault, WalletTypeCash, WalletTypeCreditCard,
		WalletTypeDebt, WalletTypeInitialBalance, WalletTypeLiabilityCredit,
		WalletTypeLoan, WalletTypeMortgage, WalletTypeReconciliation:
		return true
	default:
		return false
	}
}

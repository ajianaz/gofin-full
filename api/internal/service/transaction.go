package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/ajianaz/gofin-full/api/internal/domain"
	"github.com/ajianaz/gofin-full/api/internal/repository"
)

// TransactionService handles the core business logic for transactions.
// It enforces double-entry bookkeeping and balance calculation.
type TransactionService struct {
	txRepo     *repository.TransactionRepository
	walletRepo *repository.WalletRepository
}

func NewTransactionService(txRepo *repository.TransactionRepository, walletRepo *repository.WalletRepository) *TransactionService {
	return &TransactionService{txRepo: txRepo, walletRepo: walletRepo}
}

// CreateTransactionInput is the API input for creating a transaction.
type CreateTransactionInput struct {
	Type          string      `json:"type"`
	Description   string      `json:"description"`
	Date          time.Time   `json:"date"`
	Amount        string      `json:"amount"` // decimal string, positive
	SourceID      uuid.UUID   `json:"source_id"`
	DestinationID uuid.UUID   `json:"destination_id"`
	CurrencyID    string      `json:"currency_id"`
	CategoryIDs   []uuid.UUID `json:"category_ids,omitempty"`
	TagIDs        []uuid.UUID `json:"tag_ids,omitempty"`
	Notes         *string     `json:"notes,omitempty"`
	BudgetID      *uuid.UUID  `json:"budget_id,omitempty"`
	PiggyBankID   *uuid.UUID  `json:"piggy_bank_id,omitempty"`
}

// SplitJournalInput represents one journal in a split transaction.
type SplitJournalInput struct {
	Description   string      `json:"description"`
	Amount        string      `json:"amount"` // decimal string, positive
	SourceID      uuid.UUID   `json:"source_id"`
	DestinationID uuid.UUID   `json:"destination_id"`
	CategoryIDs   []uuid.UUID `json:"category_ids,omitempty"`
	TagIDs        []uuid.UUID `json:"tag_ids,omitempty"`
}

// CreateResult is the output of a transaction creation.
type CreateResult struct {
	GroupID   uuid.UUID `json:"group_id"`
	JournalID uuid.UUID `json:"journal_id"`
}

// CreateTransaction creates a triple-layer transaction atomically.
// It enforces double-entry: source gets debited, destination gets credited.
// All DB operations (group + journal + transactions + balance updates + category/tag links)
// are wrapped in a single database transaction for consistency.
func (s *TransactionService) CreateTransaction(ctx context.Context, userID, groupID uuid.UUID, input CreateTransactionInput) (*CreateResult, error) {
	amount, err := decimal.NewFromString(input.Amount)
	if err != nil {
		return nil, fmt.Errorf("invalid amount: %w", err)
	}
	if amount.IsZero() || amount.IsNegative() {
		return nil, fmt.Errorf("amount must be positive")
	}

	// Resolve transaction type
	typeID, err := s.txRepo.GetTransactionTypeID(ctx, input.Type)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction type: %w", err)
	}

	// Validate wallet types
	sourceWallet, err := s.walletRepo.FindByID(ctx, input.SourceID, groupID)
	if err != nil {
		return nil, fmt.Errorf("source wallet not found: %w", err)
	}
	destWallet, err := s.walletRepo.FindByID(ctx, input.DestinationID, groupID)
	if err != nil {
		return nil, fmt.Errorf("destination wallet not found: %w", err)
	}

	// Calculate debit/credit based on transaction type
	sourceAmount, destAmount, err := s.calculateAmounts(input.Type, amount)
	if err != nil {
		return nil, err
	}

	// Build the group title
	title := ""
	if len(input.Description) > 0 {
		title = input.Description
	}

	result, err := s.txRepo.CreateFullTransaction(ctx, repository.CreateFullTransactionInput{
		UserID:     userID,
		GroupID:    groupID,
		GroupTitle: title,
		Journal: &domain.TransactionJournal{
			TransactionTypeID: typeID,
			Date:              input.Date,
			Description:       input.Description,
			CurrencyID:        input.CurrencyID,
			BudgetID:          input.BudgetID,
			PiggyBankID:       input.PiggyBankID,
			Notes:             input.Notes,
		},
		SourceTxn: &domain.Transaction{
			AccountID:    sourceWallet.ID,
			Amount:       sourceAmount,
			NativeAmount: sourceAmount,
		},
		DestTxn: &domain.Transaction{
			AccountID:    destWallet.ID,
			Amount:       destAmount,
			NativeAmount: destAmount,
		},
		SourceWalletID: sourceWallet.ID,
		DestWalletID:   destWallet.ID,
		SourceAmount:   sourceAmount,
		DestAmount:     destAmount,
		CategoryIDs:    input.CategoryIDs,
		TagIDs:         input.TagIDs,
	})
	if err != nil {
		return nil, err
	}

	return &CreateResult{GroupID: result.GroupID, JournalID: result.JournalID}, nil
}

// CreateSplitTransaction creates a single group with multiple journals.
// All journals share the same type but have different amounts/descriptions.
// The entire operation (group + all journals + transactions + balance updates) is atomic.
func (s *TransactionService) CreateSplitTransaction(
	ctx context.Context, userID, groupID uuid.UUID,
	txType string, date time.Time, groupTitle string,
	journals []SplitJournalInput,
) (*CreateResult, error) {
	if len(journals) < 2 {
		return nil, fmt.Errorf("split transaction requires at least 2 journals")
	}

	typeID, err := s.txRepo.GetTransactionTypeID(ctx, txType)
	if err != nil {
		return nil, err
	}

	// Validate all amounts
	var totalAmount decimal.Decimal
	for _, j := range journals {
		amt, err := decimal.NewFromString(j.Amount)
		if err != nil {
			return nil, fmt.Errorf("invalid amount in split journal: %w", err)
		}
		if amt.IsZero() || amt.IsNegative() {
			return nil, fmt.Errorf("split journal amount must be positive")
		}
		totalAmount = totalAmount.Add(amt)
	}

	// Build split journal entries
	entries := make([]repository.SplitJournalEntryInput, len(journals))
	for i, jInput := range journals {
		amount, _ := decimal.NewFromString(jInput.Amount)
		sourceAmt, destAmt, err := s.calculateAmounts(txType, amount)
		if err != nil {
			return nil, err
		}

		entries[i] = repository.SplitJournalEntryInput{
			Journal: &domain.TransactionJournal{
				TransactionTypeID: typeID,
				Date:              date,
				Description:       jInput.Description,
				CurrencyID:        "",
			},
			SourceTxn: &domain.Transaction{
				AccountID:    jInput.SourceID,
				Amount:       sourceAmt,
				NativeAmount: sourceAmt,
			},
			DestTxn: &domain.Transaction{
				AccountID:    jInput.DestinationID,
				Amount:       destAmt,
				NativeAmount: destAmt,
			},
			SourceWalletID: jInput.SourceID,
			DestWalletID:   jInput.DestinationID,
			SourceAmount:   sourceAmt,
			DestAmount:     destAmt,
			CategoryIDs:    jInput.CategoryIDs,
			TagIDs:         jInput.TagIDs,
		}
	}

	result, err := s.txRepo.CreateSplitTransactionInTx(ctx, repository.CreateSplitTransactionInput{
		UserID:     userID,
		GroupID:    groupID,
		GroupTitle: groupTitle,
		Journals:   entries,
	})
	if err != nil {
		return nil, err
	}

	return &CreateResult{GroupID: result.GroupID, JournalID: result.JournalID}, nil
}

// calculateAmounts determines the debit/credit amounts based on transaction type.
// Returns (sourceAmount, destinationAmount, error).
//
// Withdrawal:  source gets -amount (expense), destination gets +amount
// Deposit:     source gets -amount (revenue), destination gets +amount (asset)
// Transfer:    source gets -amount, destination gets +amount
func (s *TransactionService) calculateAmounts(txType string, amount decimal.Decimal) (decimal.Decimal, decimal.Decimal, error) {
	switch domain.TransactionType(txType) {
	case domain.TransactionTypeWithdrawal, domain.TransactionTypeDeposit,
		domain.TransactionTypeTransfer, domain.TransactionTypeOpeningBalance,
		domain.TransactionTypeReconciliation, domain.TransactionTypeLiabilityCredit:
		return amount.Neg(), amount, nil
	default:
		return decimal.Zero, decimal.Zero, fmt.Errorf("unsupported transaction type: %s", txType)
	}
}

// DeleteTransaction soft-deletes a transaction group and reverses wallet balances.
// The balance reversal and soft delete are performed in a single database transaction.
func (s *TransactionService) DeleteTransaction(ctx context.Context, groupID, userID, groupIDScope uuid.UUID) error {
	return s.txRepo.DeleteFullTransaction(ctx, groupID, userID, groupIDScope)
}

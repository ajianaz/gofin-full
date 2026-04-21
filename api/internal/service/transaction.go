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

	// Create the group
	title := ""
	if len(input.Description) > 0 {
		title = input.Description
	}
	group, err := s.txRepo.CreateGroup(ctx, userID, groupID, title)
	if err != nil {
		return nil, err
	}

	// Create the journal
	journal := &domain.TransactionJournal{
		TransactionGroupID: group.ID,
		UserID:             userID,
		UserGroupID:        groupID,
		TransactionTypeID:  typeID,
		Date:               input.Date,
		Description:        input.Description,
		CurrencyID:         input.CurrencyID,
		BudgetID:           input.BudgetID,
		PiggyBankID:        input.PiggyBankID,
		Notes:              input.Notes,
	}
	journal, err = s.txRepo.CreateJournal(ctx, journal)
	if err != nil {
		return nil, err
	}

	// Create source transaction (debit)
	srcTxn := &domain.Transaction{
		TransactionJournalID: journal.ID,
		AccountID:            sourceWallet.ID,
		Amount:               sourceAmount,
		NativeAmount:         sourceAmount, // same currency for now
	}
	if _, err := s.txRepo.CreateTransaction(ctx, srcTxn); err != nil {
		return nil, err
	}

	// Create destination transaction (credit)
	dstTxn := &domain.Transaction{
		TransactionJournalID: journal.ID,
		AccountID:            destWallet.ID,
		Amount:               destAmount,
		NativeAmount:         destAmount,
	}
	if _, err := s.txRepo.CreateTransaction(ctx, dstTxn); err != nil {
		return nil, err
	}

	// Update wallet balances
	if err := s.txRepo.UpdateWalletBalance(ctx, sourceWallet.ID, sourceAmount); err != nil {
		return nil, fmt.Errorf("failed to update source balance: %w", err)
	}
	if err := s.txRepo.UpdateWalletBalance(ctx, destWallet.ID, destAmount); err != nil {
		return nil, fmt.Errorf("failed to update destination balance: %w", err)
	}

	// Link categories and tags
	if len(input.CategoryIDs) > 0 {
		s.txRepo.SetJournalCategories(ctx, journal.ID, input.CategoryIDs)
	}
	if len(input.TagIDs) > 0 {
		s.txRepo.SetJournalTags(ctx, journal.ID, input.TagIDs)
	}

	return &CreateResult{GroupID: group.ID, JournalID: journal.ID}, nil
}

// CreateSplitTransaction creates a single group with multiple journals.
// All journals share the same type but have different amounts/descriptions.
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

	// Create group
	group, err := s.txRepo.CreateGroup(ctx, userID, groupID, groupTitle)
	if err != nil {
		return nil, err
	}

	var firstJournalID uuid.UUID

	for i, jInput := range journals {
		amount, _ := decimal.NewFromString(jInput.Amount)
		sourceAmt, destAmt, err := s.calculateAmounts(txType, amount)
		if err != nil {
			return nil, err
		}

		journal := &domain.TransactionJournal{
			TransactionGroupID: group.ID,
			UserID:             userID,
			UserGroupID:        groupID,
			TransactionTypeID:  typeID,
			Date:               date,
			Order:              i,
			Description:        jInput.Description,
			CurrencyID:         "",
		}
		journal, err = s.txRepo.CreateJournal(ctx, journal)
		if err != nil {
			return nil, err
		}
		if i == 0 {
			firstJournalID = journal.ID
		}

		// Source transaction
		srcTxn := &domain.Transaction{
			TransactionJournalID: journal.ID,
			AccountID:            jInput.SourceID,
			Amount:               sourceAmt,
			NativeAmount:         sourceAmt,
		}
		if _, err := s.txRepo.CreateTransaction(ctx, srcTxn); err != nil {
			return nil, err
		}

		// Destination transaction
		dstTxn := &domain.Transaction{
			TransactionJournalID: journal.ID,
			AccountID:            jInput.DestinationID,
			Amount:               destAmt,
			NativeAmount:         destAmt,
		}
		if _, err := s.txRepo.CreateTransaction(ctx, dstTxn); err != nil {
			return nil, err
		}

		// Update balances
		if err := s.txRepo.UpdateWalletBalance(ctx, jInput.SourceID, sourceAmt); err != nil {
			return nil, err
		}
		if err := s.txRepo.UpdateWalletBalance(ctx, jInput.DestinationID, destAmt); err != nil {
			return nil, err
		}

		// Link categories and tags
		if len(jInput.CategoryIDs) > 0 {
			s.txRepo.SetJournalCategories(ctx, journal.ID, jInput.CategoryIDs)
		}
		if len(jInput.TagIDs) > 0 {
			s.txRepo.SetJournalTags(ctx, journal.ID, jInput.TagIDs)
		}
	}

	return &CreateResult{GroupID: group.ID, JournalID: firstJournalID}, nil
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
func (s *TransactionService) DeleteTransaction(ctx context.Context, groupID, userID, groupIDScope uuid.UUID) error {
	// Load the full transaction to reverse balances
	group, err := s.txRepo.FindGroupByID(ctx, groupID, groupIDScope)
	if err != nil {
		return err
	}

	// Reverse balances for all journals
	for _, journal := range group.Journals {
		for _, t := range journal.SourceTransactions {
			if err := s.txRepo.UpdateWalletBalance(ctx, t.AccountID, t.Amount.Neg()); err != nil {
				return fmt.Errorf("failed to reverse source balance: %w", err)
			}
		}
		for _, t := range journal.DestinationTransactions {
			if err := s.txRepo.UpdateWalletBalance(ctx, t.AccountID, t.Amount.Neg()); err != nil {
				return fmt.Errorf("failed to reverse destination balance: %w", err)
			}
		}
	}

	return s.txRepo.DeleteGroup(ctx, groupID, groupIDScope)
}

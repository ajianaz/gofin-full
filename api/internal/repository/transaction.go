package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"github.com/ajianaz/gofin-full/api/internal/domain"
)

// TransactionRepository handles the triple-layer transaction data access.
type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// CreateGroup inserts a new transaction group and returns it with ID.
func (r *TransactionRepository) CreateGroup(ctx context.Context, userID, groupID uuid.UUID, title string) (*domain.TransactionGroup, error) {
	now := time.Now().UTC()
	var g domain.TransactionGroup
	err := r.db.QueryRow(ctx,
		`INSERT INTO transaction_groups (user_id, user_group_id, group_title, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5) RETURNING id, user_id, user_group_id, group_title, created_at, updated_at`,
		userID, groupID, title, now, now,
	).Scan(&g.ID, &g.UserID, &g.UserGroupID, &g.GroupTitle, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction group: %w", err)
	}
	return &g, nil
}

// CreateJournal inserts a new transaction journal.
func (r *TransactionRepository) CreateJournal(ctx context.Context, j *domain.TransactionJournal) (*domain.TransactionJournal, error) {
	now := time.Now().UTC()
	var id uuid.UUID
	err := r.db.QueryRow(ctx,
		`INSERT INTO transaction_journals
		 (transaction_group_id, user_id, user_group_id, transaction_type_id,
		  date, "order", description, transaction_currency_id, foreign_currency_id,
		  budget_id, bill_id, piggy_bank_id, reconciled, notes,
		  interest_date, book_date, process_date, due_date, payment_date, invoice_date,
		  external_id, external_url, internal_reference,
		  sepa_cc, sepa_ct_op, sepa_ct_id, sepa_db, sepa_country, sepa_ep, sepa_ci, sepa_batch_id,
		  created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33)
		 RETURNING id`,
		j.TransactionGroupID, j.UserID, j.UserGroupID, j.TransactionTypeID,
		j.Date, j.Order, j.Description, j.CurrencyID, j.ForeignCurrencyID,
		j.BudgetID, j.BillID, j.PiggyBankID, j.Reconciled, j.Notes,
		j.InterestDate, j.BookDate, j.ProcessDate, j.DueDate, j.PaymentDate, j.InvoiceDate,
		j.ExternalID, j.ExternalURL, j.InternalReference,
		j.SepaCC, j.SepaCTOp, j.SepaCTID, j.SepaDB, j.SepaCountry, j.SepaEP, j.SepaCI, j.SepaBatchID,
		now, now,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction journal: %w", err)
	}
	j.ID = id
	j.CreatedAt = now
	j.UpdatedAt = now
	return j, nil
}

// CreateTransaction inserts a transaction record (debit/credit line).
func (r *TransactionRepository) CreateTransaction(ctx context.Context, t *domain.Transaction) (*domain.Transaction, error) {
	now := time.Now().UTC()
	var id uuid.UUID
	err := r.db.QueryRow(ctx,
		`INSERT INTO transactions
		 (transaction_journal_id, account_id, amount, native_amount,
		  foreign_amount, native_foreign_amount, foreign_currency_id, reconciled,
		  created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		 RETURNING id`,
		t.TransactionJournalID, t.AccountID, t.Amount, t.NativeAmount,
		t.ForeignAmount, t.NativeForeignAmount, t.ForeignCurrencyID, t.Reconciled,
		now, now,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	t.ID = id
	t.CreatedAt = now
	t.UpdatedAt = now
	return t, nil
}

// UpdateWalletBalance adjusts a wallet's virtual_balance by delta.
func (r *TransactionRepository) UpdateWalletBalance(ctx context.Context, walletID uuid.UUID, delta decimal.Decimal) error {
	_, err := r.db.Exec(ctx,
		`UPDATE wallets SET virtual_balance = virtual_balance + $1, updated_at = $2
		 WHERE id = $3 AND deleted_at IS NULL`,
		delta, time.Now().UTC(), walletID)
	if err != nil {
		return fmt.Errorf("failed to update wallet balance: %w", err)
	}
	return nil
}

// FindGroupByID fetches a transaction group with all its journals and transactions.
func (r *TransactionRepository) FindGroupByID(ctx context.Context, id, groupID uuid.UUID) (*domain.TransactionGroup, error) {
	var g domain.TransactionGroup
	var deletedAt *time.Time
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, user_group_id, group_title, created_at, updated_at, deleted_at
		 FROM transaction_groups WHERE id = $1 AND user_group_id = $2`, id, groupID,
	).Scan(&g.ID, &g.UserID, &g.UserGroupID, &g.GroupTitle, &g.CreatedAt, &g.UpdatedAt, &deletedAt)
	if err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}
	if deletedAt != nil {
		return nil, fmt.Errorf("transaction not found")
	}

	journals, err := r.findJournalsByGroupID(ctx, g.ID)
	if err != nil {
		return nil, err
	}
	g.Journals = journals
	return &g, nil
}

// findJournalsByGroupID loads all journals for a group with their transactions.
func (r *TransactionRepository) findJournalsByGroupID(ctx context.Context, groupID uuid.UUID) ([]domain.TransactionJournal, error) {
	rows, err := r.db.Query(ctx,
		`SELECT j.id, j.transaction_group_id, j.user_id, j.user_group_id,
		  j.transaction_type_id, COALESCE(tt.type, 'invalid'),
		  j.date, j."order", j.description, j.transaction_currency_id,
		  j.foreign_currency_id, j.budget_id, j.bill_id, j.piggy_bank_id,
		  j.reconciled, j.notes,
		  j.interest_date, j.book_date, j.process_date, j.due_date, j.payment_date, j.invoice_date,
		  j.external_id, j.external_url, j.internal_reference,
		  j.created_at, j.updated_at
		 FROM transaction_journals j
		 LEFT JOIN transaction_types tt ON tt.id = j.transaction_type_id
		 WHERE j.transaction_group_id = $1 AND j.deleted_at IS NULL
		 ORDER BY j."order"`, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to load journals: %w", err)
	}
	defer rows.Close()

	var journals []domain.TransactionJournal
	for rows.Next() {
		var j domain.TransactionJournal
		var tt string
		if err := rows.Scan(
			&j.ID, &j.TransactionGroupID, &j.UserID, &j.UserGroupID,
			&j.TransactionTypeID, &tt,
			&j.Date, &j.Order, &j.Description, &j.CurrencyID,
			&j.ForeignCurrencyID, &j.BudgetID, &j.BillID, &j.PiggyBankID,
			&j.Reconciled, &j.Notes,
			&j.InterestDate, &j.BookDate, &j.ProcessDate, &j.DueDate, &j.PaymentDate, &j.InvoiceDate,
			&j.ExternalID, &j.ExternalURL, &j.InternalReference,
			&j.CreatedAt, &j.UpdatedAt,
		); err != nil {
			return nil, err
		}
		j.Type = domain.TransactionType(tt)

		// Load transactions for this journal
		txns, err := r.findTransactionsByJournalID(ctx, j.ID)
		if err != nil {
			return nil, err
		}
		// Split into source and destination
		for _, t := range txns {
			// First transaction is source, second is destination
			if len(j.SourceTransactions) == 0 {
				j.SourceTransactions = append(j.SourceTransactions, t)
			} else {
				j.DestinationTransactions = append(j.DestinationTransactions, t)
			}
		}
		journals = append(journals, j)
	}
	return journals, rows.Err()
}

func (r *TransactionRepository) findTransactionsByJournalID(ctx context.Context, journalID uuid.UUID) ([]domain.Transaction, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, transaction_journal_id, account_id, amount, native_amount,
		  foreign_amount, native_foreign_amount, foreign_currency_id, reconciled,
		  created_at, updated_at
		 FROM transactions WHERE transaction_journal_id = $1`, journalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txns []domain.Transaction
	for rows.Next() {
		var t domain.Transaction
		if err := rows.Scan(
			&t.ID, &t.TransactionJournalID, &t.AccountID,
			&t.Amount, &t.NativeAmount,
			&t.ForeignAmount, &t.NativeForeignAmount, &t.ForeignCurrencyID,
			&t.Reconciled, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		txns = append(txns, t)
	}
	return txns, rows.Err()
}

// ListGroups returns paginated transaction groups for a group with optional filters.
type TransactionFilter struct {
	DateFrom *time.Time
	DateTo   *time.Time
	Type     string // transaction type filter
	WalletID *uuid.UUID // filter by source or destination wallet
	Page     int
	PerPage  int
}

func (f *TransactionFilter) defaults() {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PerPage < 1 || f.PerPage > 100 {
		f.PerPage = 50
	}
}

// ListGroups returns filtered transaction groups with total count.
func (r *TransactionRepository) ListGroups(ctx context.Context, groupID uuid.UUID, f TransactionFilter) ([]domain.TransactionGroup, int64, error) {
	f.defaults()
	offset := (f.Page - 1) * f.PerPage

	where := "g.user_group_id = $1 AND g.deleted_at IS NULL"
	args := []interface{}{groupID}
	argN := 2

	if f.DateFrom != nil {
		where += fmt.Sprintf(" AND j.date >= $%d", argN)
		args = append(args, f.DateFrom)
		argN++
	}
	if f.DateTo != nil {
		where += fmt.Sprintf(" AND j.date <= $%d", argN)
		args = append(args, f.DateTo)
		argN++
	}
	if f.Type != "" {
		where += fmt.Sprintf(" AND tt.type = $%d", argN)
		args = append(args, f.Type)
		argN++
	}
	if f.WalletID != nil {
		where += fmt.Sprintf(" AND EXISTS (SELECT 1 FROM transactions t WHERE t.transaction_journal_id = j.id AND t.account_id = $%d)", argN)
		args = append(args, *f.WalletID)
		argN++
	}

	// Count
	var total int64
	countQuery := fmt.Sprintf(
		`SELECT COUNT(DISTINCT g.id) FROM transaction_groups g
		 JOIN transaction_journals j ON j.transaction_group_id = g.id AND j.deleted_at IS NULL
		 LEFT JOIN transaction_types tt ON tt.id = j.transaction_type_id
		 WHERE %s`, where)
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count transactions: %w", err)
	}

	// Fetch groups
	dataQuery := fmt.Sprintf(
		`SELECT DISTINCT g.id, g.user_id, g.user_group_id, g.group_title, g.created_at, g.updated_at
		 FROM transaction_groups g
		 JOIN transaction_journals j ON j.transaction_group_id = g.id AND j.deleted_at IS NULL
		 LEFT JOIN transaction_types tt ON tt.id = j.transaction_type_id
		 WHERE %s
		 ORDER BY g.created_at DESC
		 LIMIT $%d OFFSET $%d`, where, argN, argN+1)
	args = append(args, f.PerPage, offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list transactions: %w", err)
	}
	defer rows.Close()

	var groups []domain.TransactionGroup
	for rows.Next() {
		var g domain.TransactionGroup
		if err := rows.Scan(&g.ID, &g.UserID, &g.UserGroupID, &g.GroupTitle, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, 0, err
		}
		groups = append(groups, g)
	}

	return groups, total, rows.Err()
}

// UpdateJournal updates mutable fields on a journal.
func (r *TransactionRepository) UpdateJournal(ctx context.Context, id, groupID uuid.UUID, description string, date *time.Time, notes *string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE transaction_journals SET
		  description = COALESCE(NULLIF($1, ''), description),
		  date = COALESCE($2, date),
		  notes = $3,
		  updated_at = $4
		 WHERE id = $5 AND user_group_id = $6 AND deleted_at IS NULL`,
		description, date, notes, time.Now().UTC(), id, groupID)
	return err
}

// DeleteGroup soft-deletes a transaction group and all its journals.
func (r *TransactionRepository) DeleteGroup(ctx context.Context, id, groupID uuid.UUID) error {
	now := time.Now().UTC()
	// Soft-delete journals first
	_, err := r.db.Exec(ctx,
		`UPDATE transaction_journals SET deleted_at = $1
		 WHERE transaction_group_id = $2 AND user_group_id = $3 AND deleted_at IS NULL`,
		now, id, groupID)
	if err != nil {
		return fmt.Errorf("failed to delete transaction journals: %w", err)
	}
	// Soft-delete group
	_, err = r.db.Exec(ctx,
		`UPDATE transaction_groups SET deleted_at = $1
		 WHERE id = $2 AND user_group_id = $3 AND deleted_at IS NULL`,
		now, id, groupID)
	if err != nil {
		return fmt.Errorf("failed to delete transaction group: %w", err)
	}
	return nil
}

// GetTransactionTypeID resolves a transaction type string to its DB ID.
func (r *TransactionRepository) GetTransactionTypeID(ctx context.Context, typeName string) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.db.QueryRow(ctx,
		`SELECT id FROM transaction_types WHERE type = $1`, typeName).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("unknown transaction type '%s'", typeName)
	}
	return id, nil
}

// SetJournalCategories replaces category links for a journal.
func (r *TransactionRepository) SetJournalCategories(ctx context.Context, journalID uuid.UUID, categoryIDs []uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `DELETE FROM category_transaction WHERE transaction_journal_id = $1`, journalID)
	if err != nil {
		return err
	}
	for _, catID := range categoryIDs {
		_, err = tx.Exec(ctx,
			`INSERT INTO category_transaction (category_id, transaction_journal_id) VALUES ($1, $2)`,
			catID, journalID)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

// SetJournalTags replaces tag links for a journal.
func (r *TransactionRepository) SetJournalTags(ctx context.Context, journalID uuid.UUID, tagIDs []uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `DELETE FROM journal_tag WHERE transaction_journal_id = $1`, journalID)
	if err != nil {
		return err
	}
	for _, tagID := range tagIDs {
		_, err = tx.Exec(ctx,
			`INSERT INTO journal_tag (tag_id, transaction_journal_id) VALUES ($1, $2)`,
			tagID, journalID)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

// Ensure strings is used.
var _ = strings.TrimSpace

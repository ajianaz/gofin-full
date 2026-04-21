package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"github.com/ajianaz/gofin-full/api/internal/domain"
)

type RecurrenceRepository struct {
	db *pgxpool.Pool
}

func NewRecurrenceRepository(db *pgxpool.Pool) *RecurrenceRepository {
	return &RecurrenceRepository{db: db}
}

func (r *RecurrenceRepository) Create(ctx context.Context, userID, groupID uuid.UUID, title string, firstDate time.Time, repeatFreq string) (*domain.Recurrence, error) {
	now := time.Now().UTC()
	var rec domain.Recurrence
	err := r.db.QueryRow(ctx,
		`INSERT INTO recurrences (user_id, user_group_id, title, first_date, repeat_freq, active, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,TRUE,$6,$7)
		 RETURNING id, user_id, user_group_id, title, description, first_date, latest_date, repeat_until, repeat_freq, skip, active, apply_rules, created_at, updated_at`,
		userID, groupID, title, firstDate, repeatFreq, now, now,
	).Scan(&rec.ID, &rec.UserID, &rec.UserGroupID, &rec.Title, &rec.Description,
		&rec.FirstDate, &rec.LatestDate, &rec.RepeatUntil,
		&rec.RepeatFreq, &rec.Skip, &rec.Active, &rec.ApplyRules,
		&rec.CreatedAt, &rec.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create recurrence: %w", err)
	}
	return &rec, nil
}

func (r *RecurrenceRepository) FindByID(ctx context.Context, id, groupID uuid.UUID) (*domain.Recurrence, error) {
	var rec domain.Recurrence
	var deletedAt *time.Time
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, user_group_id, title, description, first_date, latest_date, repeat_until,
		  repeat_freq, skip, active, apply_rules, created_at, updated_at, deleted_at
		 FROM recurrences WHERE id = $1 AND user_group_id = $2`, id, groupID,
	).Scan(&rec.ID, &rec.UserID, &rec.UserGroupID, &rec.Title, &rec.Description,
		&rec.FirstDate, &rec.LatestDate, &rec.RepeatUntil,
		&rec.RepeatFreq, &rec.Skip, &rec.Active, &rec.ApplyRules,
		&rec.CreatedAt, &rec.UpdatedAt, &deletedAt)
	if err != nil {
		return nil, fmt.Errorf("recurrence not found: %w", err)
	}
	if deletedAt != nil {
		return nil, fmt.Errorf("recurrence not found")
	}

	txns, _ := r.findTransactions(ctx, rec.ID)
	reps, _ := r.findRepetitions(ctx, rec.ID)
	rec.Transactions = txns
	rec.Repetitions = reps
	return &rec, nil
}

func (r *RecurrenceRepository) List(ctx context.Context, groupID uuid.UUID) ([]domain.Recurrence, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, user_group_id, title, description, first_date, latest_date, repeat_until,
		  repeat_freq, skip, active, apply_rules, created_at, updated_at
		 FROM recurrences WHERE user_group_id = $1 AND deleted_at IS NULL ORDER BY title`, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to list recurrences: %w", err)
	}
	defer rows.Close()

	var recurrences []domain.Recurrence
	for rows.Next() {
		var rec domain.Recurrence
		if err := rows.Scan(&rec.ID, &rec.UserID, &rec.UserGroupID, &rec.Title, &rec.Description,
			&rec.FirstDate, &rec.LatestDate, &rec.RepeatUntil,
			&rec.RepeatFreq, &rec.Skip, &rec.Active, &rec.ApplyRules,
			&rec.CreatedAt, &rec.UpdatedAt); err != nil {
			return nil, err
		}
		recurrences = append(recurrences, rec)
	}
	return recurrences, rows.Err()
}

func (r *RecurrenceRepository) Update(ctx context.Context, id, groupID uuid.UUID, title string, repeatFreq string, active *bool, description *string, repeatUntil *time.Time) error {
	_, err := r.db.Exec(ctx,
		`UPDATE recurrences SET
		  title = COALESCE(NULLIF($1, ''), title),
		  repeat_freq = COALESCE(NULLIF($2, ''), repeat_freq),
		  active = COALESCE($3, active),
		  description = $4,
		  repeat_until = $5,
		  updated_at = $6
		 WHERE id = $7 AND user_group_id = $8 AND deleted_at IS NULL`,
		title, repeatFreq, active, description, repeatUntil, time.Now().UTC(), id, groupID)
	return err
}

func (r *RecurrenceRepository) Delete(ctx context.Context, id, groupID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`UPDATE recurrences SET deleted_at = $1 WHERE id = $2 AND user_group_id = $3 AND deleted_at IS NULL`,
		time.Now().UTC(), id, groupID)
	return err
}

// AddTransaction adds a transaction template to a recurrence.
func (r *RecurrenceRepository) AddTransaction(ctx context.Context, recID uuid.UUID, txn *domain.RecurringTransaction) error {
	now := time.Now().UTC()
	_, err := r.db.Exec(ctx,
		`INSERT INTO recurring_transactions
		 (recurrence_id, type, description, amount, transaction_currency_id, source_id, destination_id, budget_id, category_id, piggy_bank_id, "order", created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		recID, txn.Type, txn.Description, txn.Amount, txn.CurrencyID,
		txn.SourceID, txn.DestinationID, txn.BudgetID, txn.CategoryID, txn.PiggyBankID,
		txn.Order, now, now)
	return err
}

func (r *RecurrenceRepository) findTransactions(ctx context.Context, recID uuid.UUID) ([]domain.RecurringTransaction, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, recurrence_id, type, description, amount, transaction_currency_id, source_id, destination_id,
		  budget_id, category_id, piggy_bank_id, "order", created_at, updated_at
		 FROM recurring_transactions WHERE recurrence_id = $1 ORDER BY "order"`, recID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txns []domain.RecurringTransaction
	for rows.Next() {
		var t domain.RecurringTransaction
		if err := rows.Scan(&t.ID, &t.RecurrenceID, &t.Type, &t.Description, &t.Amount,
			&t.CurrencyID, &t.SourceID, &t.DestinationID, &t.BudgetID, &t.CategoryID, &t.PiggyBankID,
			&t.Order, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		txns = append(txns, t)
	}
	return txns, rows.Err()
}

func (r *RecurrenceRepository) findRepetitions(ctx context.Context, recID uuid.UUID) ([]domain.RecurringRepetition, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, recurrence_id, relevant_date, created_at, updated_at
		 FROM recurring_repetitions WHERE recurrence_id = $1 ORDER BY relevant_date DESC`, recID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reps []domain.RecurringRepetition
	for rows.Next() {
		var rep domain.RecurringRepetition
		if err := rows.Scan(&rep.ID, &rep.RecurrenceID, &rep.RelevantDate, &rep.CreatedAt, &rep.UpdatedAt); err != nil {
			return nil, err
		}
		reps = append(reps, rep)
	}
	return reps, rows.Err()
}

// Ensure unused import is not flagged.
var _ = decimal.Decimal{}

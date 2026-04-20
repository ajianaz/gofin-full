package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"github.com/ajianaz/gofin-full/api/internal/domain"
)

type BillRepository struct {
	db *pgxpool.Pool
}

func NewBillRepository(db *pgxpool.Pool) *BillRepository {
	return &BillRepository{db: db}
}

func (r *BillRepository) Create(ctx context.Context, userID, groupID int64, name string, amountMin, amountMax decimal.Decimal, date time.Time, repeatFreq, currencyID string, order int) (*domain.Bill, error) {
	now := time.Now().UTC()
	var b domain.Bill
	err := r.db.QueryRow(ctx,
		`INSERT INTO bills (user_id, user_group_id, name, amount_min, amount_max, date, repeat_freq, currency_id, "order", active, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,TRUE,$10,$11)
		 RETURNING id, user_id, user_group_id, name, amount_min, amount_max, date, end_date, repeat_freq, skip, active, notes, currency_id, created_at, updated_at`,
		userID, groupID, name, amountMin, amountMax, date, repeatFreq, currencyID, order, now, now,
	).Scan(&b.ID, &b.UserID, &b.UserGroupID, &b.Name, &b.AmountMin, &b.AmountMax, &b.Date, &b.EndDate,
		&b.RepeatFreq, &b.Skip, &b.Active, &b.Notes, &b.CurrencyID, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create bill: %w", err)
	}
	return &b, nil
}

func (r *BillRepository) FindByID(ctx context.Context, id, groupID int64) (*domain.Bill, error) {
	var b domain.Bill
	var deletedAt *time.Time
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, user_group_id, name, amount_min, amount_max, date, end_date, repeat_freq, skip, active, notes, currency_id, created_at, updated_at, deleted_at
		 FROM bills WHERE id = $1 AND user_group_id = $2`, id, groupID,
	).Scan(&b.ID, &b.UserID, &b.UserGroupID, &b.Name, &b.AmountMin, &b.AmountMax, &b.Date, &b.EndDate,
		&b.RepeatFreq, &b.Skip, &b.Active, &b.Notes, &b.CurrencyID, &b.CreatedAt, &b.UpdatedAt, &deletedAt)
	if err != nil {
		return nil, fmt.Errorf("bill not found: %w", err)
	}
	if deletedAt != nil {
		return nil, fmt.Errorf("bill not found")
	}
	return &b, nil
}

func (r *BillRepository) List(ctx context.Context, groupID int64) ([]domain.Bill, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, user_group_id, name, amount_min, amount_max, date, end_date, repeat_freq, skip, active, notes, currency_id, created_at, updated_at
		 FROM bills WHERE user_group_id = $1 AND deleted_at IS NULL ORDER BY "order", name`, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to list bills: %w", err)
	}
	defer rows.Close()

	var bills []domain.Bill
	for rows.Next() {
		var b domain.Bill
		if err := rows.Scan(&b.ID, &b.UserID, &b.UserGroupID, &b.Name, &b.AmountMin, &b.AmountMax, &b.Date, &b.EndDate,
			&b.RepeatFreq, &b.Skip, &b.Active, &b.Notes, &b.CurrencyID, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		bills = append(bills, b)
	}
	return bills, rows.Err()
}

func (r *BillRepository) Update(ctx context.Context, id, groupID int64, name string, active *bool, notes *string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE bills SET name = COALESCE(NULLIF($1, ''), name), active = COALESCE($2, active), notes = $3, updated_at = $4
		 WHERE id = $5 AND user_group_id = $6 AND deleted_at IS NULL`,
		name, active, notes, time.Now().UTC(), id, groupID)
	return err
}

func (r *BillRepository) Delete(ctx context.Context, id, groupID int64) error {
	_, err := r.db.Exec(ctx,
		`UPDATE bills SET deleted_at = $1 WHERE id = $2 AND user_group_id = $3 AND deleted_at IS NULL`,
		time.Now().UTC(), id, groupID)
	return err
}

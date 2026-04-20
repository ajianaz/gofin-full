package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"github.com/azfirazka/gofin-full/api/internal/domain"
)

type PiggyBankRepository struct {
	db *pgxpool.Pool
}

func NewPiggyBankRepository(db *pgxpool.Pool) *PiggyBankRepository {
	return &PiggyBankRepository{db: db}
}

func (r *PiggyBankRepository) Create(ctx context.Context, pb *domain.PiggyBank, groupID int64) (*domain.PiggyBank, error) {
	now := time.Now().UTC()

	// Verify the wallet belongs to the group
	var walletGroupID int64
	err := r.db.QueryRow(ctx,
		`SELECT user_group_id FROM wallets WHERE id = $1 AND deleted_at IS NULL`, pb.AccountID,
	).Scan(&walletGroupID)
	if err != nil {
		return nil, fmt.Errorf("wallet not found")
	}
	if walletGroupID != groupID {
		return nil, fmt.Errorf("wallet does not belong to this group")
	}

	var id int64
	err = r.db.QueryRow(ctx,
		`INSERT INTO piggy_banks (account_id, name, target_amount, start_date, target_date, "order", notes, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		 RETURNING id`,
		pb.AccountID, pb.Name, pb.TargetAmount, pb.StartDate, pb.TargetDate, pb.Order, pb.Notes, now, now,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to create piggy bank: %w", err)
	}
	pb.ID = id
	pb.CreatedAt = now
	pb.UpdatedAt = now
	return pb, nil
}

func (r *PiggyBankRepository) FindByID(ctx context.Context, id, groupID int64) (*domain.PiggyBank, error) {
	var pb domain.PiggyBank
	var deletedAt *time.Time
	err := r.db.QueryRow(ctx,
		`SELECT pb.id, pb.account_id, pb.name, pb.target_amount, pb.start_date, pb.target_date, pb."order", pb.notes, pb.created_at, pb.updated_at, pb.deleted_at
		 FROM piggy_banks pb
		 JOIN wallets w ON w.id = pb.account_id
		 WHERE pb.id = $1 AND w.user_group_id = $2`, id, groupID,
	).Scan(&pb.ID, &pb.AccountID, &pb.Name, &pb.TargetAmount, &pb.StartDate, &pb.TargetDate, &pb.Order, &pb.Notes, &pb.CreatedAt, &pb.UpdatedAt, &deletedAt)
	if err != nil {
		return nil, fmt.Errorf("piggy bank not found: %w", err)
	}
	if deletedAt != nil {
		return nil, fmt.Errorf("piggy bank not found")
	}

	// Calculate current amount from events
	currentAmount, err := r.getCurrentAmount(ctx, pb.ID)
	if err == nil {
		pb.CurrentAmount = currentAmount
		pb.LeftToTarget = pb.TargetAmount.Sub(currentAmount)
		if pb.TargetAmount.IsPositive() {
			pb.Percentage = float64(currentAmount.IntPart()) / float64(pb.TargetAmount.IntPart()) * 100
		}
	}

	return &pb, nil
}

func (r *PiggyBankRepository) List(ctx context.Context, walletID, groupID int64) ([]domain.PiggyBank, error) {
	rows, err := r.db.Query(ctx,
		`SELECT pb.id, pb.account_id, pb.name, pb.target_amount, pb.start_date, pb.target_date, pb."order", pb.notes, pb.created_at, pb.updated_at
		 FROM piggy_banks pb
		 JOIN wallets w ON w.id = pb.account_id
		 WHERE pb.account_id = $1 AND w.user_group_id = $2 AND pb.deleted_at IS NULL
		 ORDER BY pb."order", pb.name`, walletID, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to list piggy banks: %w", err)
	}
	defer rows.Close()

	var piggyBanks []domain.PiggyBank
	for rows.Next() {
		var pb domain.PiggyBank
		if err := rows.Scan(&pb.ID, &pb.AccountID, &pb.Name, &pb.TargetAmount, &pb.StartDate, &pb.TargetDate, &pb.Order, &pb.Notes, &pb.CreatedAt, &pb.UpdatedAt); err != nil {
			return nil, err
		}
		piggyBanks = append(piggyBanks, pb)
	}
	return piggyBanks, rows.Err()
}

func (r *PiggyBankRepository) Update(ctx context.Context, id, groupID int64, name string, targetAmount *decimal.Decimal, startDate, targetDate *time.Time, notes *string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE piggy_banks SET
		  name = COALESCE(NULLIF($1, ''), name),
		  target_amount = COALESCE($2, target_amount),
		  start_date = $3,
		  target_date = $4,
		  notes = $5,
		  updated_at = $6
		 WHERE id = $7 AND deleted_at IS NULL
		 AND EXISTS (SELECT 1 FROM wallets w WHERE w.id = piggy_banks.account_id AND w.user_group_id = $8)`,
		name, targetAmount, startDate, targetDate, notes, time.Now().UTC(), id, groupID)
	return err
}

func (r *PiggyBankRepository) Delete(ctx context.Context, id, groupID int64) error {
	_, err := r.db.Exec(ctx,
		`UPDATE piggy_banks SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL
		 AND EXISTS (SELECT 1 FROM wallets w WHERE w.id = piggy_banks.account_id AND w.user_group_id = $3)`,
		time.Now().UTC(), id, groupID)
	return err
}

// AddMoney creates a piggy bank event and updates the associated wallet balance.
func (r *PiggyBankRepository) AddMoney(ctx context.Context, piggyBankID, groupID int64, amount decimal.Decimal) (*domain.PiggyBankEvent, error) {
	now := time.Now().UTC()

	// Verify piggy bank belongs to group
	var accountID int64
	err := r.db.QueryRow(ctx,
		`SELECT pb.account_id FROM piggy_banks pb
		 JOIN wallets w ON w.id = pb.account_id
		 WHERE pb.id = $1 AND w.user_group_id = $2 AND pb.deleted_at IS NULL`,
		piggyBankID, groupID,
	).Scan(&accountID)
	if err != nil {
		return nil, fmt.Errorf("piggy bank not found in group")
	}

	var evt domain.PiggyBankEvent
	err = r.db.QueryRow(ctx,
		`INSERT INTO piggy_bank_events (piggy_bank_id, amount, created_at, updated_at)
		 VALUES ($1,$2,$3,$4)
		 RETURNING id, piggy_bank_id, amount, created_at, updated_at`,
		piggyBankID, amount, now, now,
	).Scan(&evt.ID, &evt.PiggyBankID, &evt.Amount, &evt.CreatedAt, &evt.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to add money to piggy bank: %w", err)
	}

	// Update the wallet balance (reduce it, money moved to piggy bank)
	_, err = r.db.Exec(ctx,
		`UPDATE wallets SET virtual_balance = virtual_balance - $1, updated_at = $2
		 WHERE id = $3 AND deleted_at IS NULL`,
		amount, now, accountID)

	return &evt, err
}

// RemoveMoney creates a negative piggy bank event and returns money to the wallet.
func (r *PiggyBankRepository) RemoveMoney(ctx context.Context, piggyBankID, groupID int64, amount decimal.Decimal) (*domain.PiggyBankEvent, error) {
	return r.AddMoney(ctx, piggyBankID, groupID, amount.Neg())
}

func (r *PiggyBankRepository) getCurrentAmount(ctx context.Context, piggyBankID int64) (decimal.Decimal, error) {
	var total decimal.Decimal
	err := r.db.QueryRow(ctx,
		`SELECT COALESCE(SUM(amount), 0) FROM piggy_bank_events WHERE piggy_bank_id = $1`, piggyBankID).Scan(&total)
	return total, err
}

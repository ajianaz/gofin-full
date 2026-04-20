package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"github.com/azfirazka/gofin-full/api/internal/domain"
)

type BudgetRepository struct {
	db *pgxpool.Pool
}

func NewBudgetRepository(db *pgxpool.Pool) *BudgetRepository {
	return &BudgetRepository{db: db}
}

func (r *BudgetRepository) Create(ctx context.Context, userID, groupID int64, name string, order int) (*domain.Budget, error) {
	now := time.Now().UTC()
	var b domain.Budget
	err := r.db.QueryRow(ctx,
		`INSERT INTO budgets (user_id, user_group_id, name, active, "order", created_at, updated_at)
		 VALUES ($1,$2,$3,TRUE,$4,$5,$6)
		 RETURNING id, user_id, user_group_id, name, active, "order", created_at, updated_at`,
		userID, groupID, name, order, now, now,
	).Scan(&b.ID, &b.UserID, &b.UserGroupID, &b.Name, &b.Active, &b.Order, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create budget: %w", err)
	}
	return &b, nil
}

func (r *BudgetRepository) FindByID(ctx context.Context, id, groupID int64) (*domain.Budget, error) {
	var b domain.Budget
	var deletedAt *time.Time
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, user_group_id, name, active, "order", created_at, updated_at, deleted_at
		 FROM budgets WHERE id = $1 AND user_group_id = $2`, id, groupID,
	).Scan(&b.ID, &b.UserID, &b.UserGroupID, &b.Name, &b.Active, &b.Order, &b.CreatedAt, &b.UpdatedAt, &deletedAt)
	if err != nil {
		return nil, fmt.Errorf("budget not found: %w", err)
	}
	if deletedAt != nil {
		return nil, fmt.Errorf("budget not found")
	}

	limits, _ := r.findLimits(ctx, b.ID)
	b.Limits = limits
	return &b, nil
}

func (r *BudgetRepository) List(ctx context.Context, groupID int64) ([]domain.Budget, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, user_group_id, name, active, "order", created_at, updated_at
		 FROM budgets WHERE user_group_id = $1 AND deleted_at IS NULL ORDER BY "order", name`, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to list budgets: %w", err)
	}
	defer rows.Close()

	var budgets []domain.Budget
	for rows.Next() {
		var b domain.Budget
		if err := rows.Scan(&b.ID, &b.UserID, &b.UserGroupID, &b.Name, &b.Active, &b.Order, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		budgets = append(budgets, b)
	}
	return budgets, rows.Err()
}

func (r *BudgetRepository) Update(ctx context.Context, id, groupID int64, name string, active *bool) error {
	_, err := r.db.Exec(ctx,
		`UPDATE budgets SET
		  name = COALESCE(NULLIF($1, ''), name),
		  active = COALESCE($2, active),
		  updated_at = $3
		 WHERE id = $4 AND user_group_id = $5 AND deleted_at IS NULL`,
		name, active, time.Now().UTC(), id, groupID)
	return err
}

func (r *BudgetRepository) Delete(ctx context.Context, id, groupID int64) error {
	_, err := r.db.Exec(ctx,
		`UPDATE budgets SET deleted_at = $1 WHERE id = $2 AND user_group_id = $3 AND deleted_at IS NULL`,
		time.Now().UTC(), id, groupID)
	return err
}

func (r *BudgetRepository) findLimits(ctx context.Context, budgetID int64) ([]domain.BudgetLimit, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, budget_id, start, "end", amount, created_at, updated_at
		 FROM budget_limits WHERE budget_id = $1 ORDER BY start`, budgetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var limits []domain.BudgetLimit
	for rows.Next() {
		var l domain.BudgetLimit
		if err := rows.Scan(&l.ID, &l.BudgetID, &l.Start, &l.End, &l.Amount, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		limits = append(limits, l)
	}
	return limits, rows.Err()
}

func (r *BudgetRepository) CreateLimit(ctx context.Context, budgetID int64, start, end time.Time, amount decimal.Decimal) (*domain.BudgetLimit, error) {
	now := time.Now().UTC()
	var l domain.BudgetLimit
	err := r.db.QueryRow(ctx,
		`INSERT INTO budget_limits (budget_id, start, "end", amount, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 RETURNING id, budget_id, start, "end", amount, created_at, updated_at`,
		budgetID, start, end, amount, now, now,
	).Scan(&l.ID, &l.BudgetID, &l.Start, &l.End, &l.Amount, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create budget limit: %w", err)
	}
	return &l, nil
}

func (r *BudgetRepository) DeleteLimit(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM budget_limits WHERE id = $1`, id)
	return err
}

func (r *BudgetRepository) SetAutoBudget(ctx context.Context, budgetID int64, autoType, period string) error {
	now := time.Now().UTC()
	_, err := r.db.Exec(ctx,
		`INSERT INTO auto_budgets (budget_id, auto_budget_type, period, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5)
		 ON CONFLICT (budget_id) DO UPDATE SET auto_budget_type = $2, period = $3, updated_at = $5`,
		budgetID, autoType, period, now, now)
	return err
}

func (r *BudgetRepository) RemoveAutoBudget(ctx context.Context, budgetID int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM auto_budgets WHERE budget_id = $1`, budgetID)
	return err
}

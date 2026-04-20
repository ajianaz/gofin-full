package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ajianaz/gofin-full/api/internal/domain"
)

type AccountTypeRepository struct {
	db *pgxpool.Pool
}

func NewAccountTypeRepository(db *pgxpool.Pool) *AccountTypeRepository {
	return &AccountTypeRepository{db: db}
}

func (r *AccountTypeRepository) List(ctx context.Context) ([]domain.AccountType, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, type, created_at, updated_at FROM account_types ORDER BY type`)
	if err != nil {
		return nil, fmt.Errorf("failed to list account types: %w", err)
	}
	defer rows.Close()

	var types []domain.AccountType
	for rows.Next() {
		var at domain.AccountType
		if err := rows.Scan(&at.ID, &at.Type, &at.CreatedAt, &at.UpdatedAt); err != nil {
			return nil, err
		}
		types = append(types, at)
	}
	return types, rows.Err()
}

func (r *AccountTypeRepository) FindByType(ctx context.Context, t string) (*domain.AccountType, error) {
	var at domain.AccountType
	err := r.db.QueryRow(ctx,
		`SELECT id, type, created_at, updated_at FROM account_types WHERE type = $1`, t,
	).Scan(&at.ID, &at.Type, &at.CreatedAt, &at.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("account type not found: %w", err)
	}
	return &at, nil
}

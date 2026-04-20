package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ajianaz/gofin-full/api/internal/domain"
)

type CurrencyRepository struct {
	db *pgxpool.Pool
}

func NewCurrencyRepository(db *pgxpool.Pool) *CurrencyRepository {
	return &CurrencyRepository{db: db}
}

func (r *CurrencyRepository) List(ctx context.Context) ([]domain.Currency, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, code, name, symbol, decimal_places, enabled, created_at, updated_at
		 FROM currencies WHERE deleted_at IS NULL AND enabled = TRUE ORDER BY code`)
	if err != nil {
		return nil, fmt.Errorf("failed to list currencies: %w", err)
	}
	defer rows.Close()

	var currencies []domain.Currency
	for rows.Next() {
		var c domain.Currency
		if err := rows.Scan(&c.ID, &c.Code, &c.Name, &c.Symbol, &c.DecimalPlaces, &c.Enabled, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		currencies = append(currencies, c)
	}
	return currencies, rows.Err()
}

func (r *CurrencyRepository) FindByCode(ctx context.Context, code string) (*domain.Currency, error) {
	var c domain.Currency
	var deletedAt *time.Time
	err := r.db.QueryRow(ctx,
		`SELECT id, code, name, symbol, decimal_places, enabled, created_at, updated_at, deleted_at
		 FROM currencies WHERE code = $1`, code,
	).Scan(&c.ID, &c.Code, &c.Name, &c.Symbol, &c.DecimalPlaces, &c.Enabled, &c.CreatedAt, &c.UpdatedAt, &deletedAt)
	if err != nil {
		return nil, fmt.Errorf("currency not found: %w", err)
	}
	if deletedAt != nil {
		return nil, fmt.Errorf("currency not found")
	}
	return &c, nil
}

func (r *CurrencyRepository) FindByID(ctx context.Context, id int64) (*domain.Currency, error) {
	var c domain.Currency
	var deletedAt *time.Time
	err := r.db.QueryRow(ctx,
		`SELECT id, code, name, symbol, decimal_places, enabled, created_at, updated_at, deleted_at
		 FROM currencies WHERE id = $1`, id,
	).Scan(&c.ID, &c.Code, &c.Name, &c.Symbol, &c.DecimalPlaces, &c.Enabled, &c.CreatedAt, &c.UpdatedAt, &deletedAt)
	if err != nil {
		return nil, fmt.Errorf("currency not found: %w", err)
	}
	return &c, nil
}

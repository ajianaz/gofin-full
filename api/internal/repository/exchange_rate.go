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

type ExchangeRateRepository struct {
	db *pgxpool.Pool
}

func NewExchangeRateRepository(db *pgxpool.Pool) *ExchangeRateRepository {
	return &ExchangeRateRepository{db: db}
}

func (r *ExchangeRateRepository) Create(ctx context.Context, userID, groupID uuid.UUID, from, to string, rate decimal.Decimal, date time.Time) (*domain.ExchangeRate, error) {
	now := time.Now().UTC()
	var er domain.ExchangeRate
	err := r.db.QueryRow(ctx,
		`INSERT INTO exchange_rates (user_id, user_group_id, from_currency_id, to_currency_id, rate, date, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		 RETURNING id, user_id, user_group_id, from_currency_id, to_currency_id, rate, date, created_at, updated_at`,
		userID, groupID, from, to, rate, date, now, now,
	).Scan(&er.ID, &er.UserID, &er.UserGroupID, &er.FromCurrencyID, &er.ToCurrencyID, &er.Rate, &er.Date, &er.CreatedAt, &er.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create exchange rate: %w", err)
	}
	return &er, nil
}

func (r *ExchangeRateRepository) List(ctx context.Context, groupID uuid.UUID) ([]domain.ExchangeRate, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, user_group_id, from_currency_id, to_currency_id, rate, date, created_at, updated_at
		 FROM exchange_rates WHERE user_group_id = $1 ORDER BY date DESC`, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to list exchange rates: %w", err)
	}
	defer rows.Close()

	var rates []domain.ExchangeRate
	for rows.Next() {
		var er domain.ExchangeRate
		if err := rows.Scan(&er.ID, &er.UserID, &er.UserGroupID, &er.FromCurrencyID, &er.ToCurrencyID, &er.Rate, &er.Date, &er.CreatedAt, &er.UpdatedAt); err != nil {
			return nil, err
		}
		rates = append(rates, er)
	}
	return rates, rows.Err()
}

func (r *ExchangeRateRepository) FindRate(ctx context.Context, groupID uuid.UUID, from, to string, date time.Time) (decimal.Decimal, error) {
	var rate decimal.Decimal
	err := r.db.QueryRow(ctx,
		`SELECT COALESCE(rate, 0) FROM exchange_rates
		 WHERE user_group_id = $1 AND from_currency_id = $2 AND to_currency_id = $3 AND date <= $4
		 ORDER BY date DESC LIMIT 1`,
		groupID, from, to, date,
	).Scan(&rate)
	if err != nil {
		return decimal.Zero, err
	}
	return rate, nil
}

func (r *ExchangeRateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM exchange_rates WHERE id = $1`, id)
	return err
}

package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/ajianaz/gofin-full/api/internal/repository"
)

type ExchangeRateService struct {
	repo *repository.ExchangeRateRepository
}

func NewExchangeRateService(repo *repository.ExchangeRateRepository) *ExchangeRateService {
	return &ExchangeRateService{repo: repo}
}

// GetRate implements the rate lookup chain: DB → reverse → cross-rate via EUR.
func (s *ExchangeRateService) GetRate(ctx context.Context, groupID uuid.UUID, from, to string, date time.Time) (decimal.Decimal, error) {
	if from == to {
		return decimal.NewFromInt(1), nil
	}

	// 1. Direct rate
	rate, err := s.repo.FindRate(ctx, groupID, from, to, date)
	if err == nil && !rate.IsZero() {
		return rate, nil
	}

	// 2. Reverse rate (1 / rate)
	reverseRate, err := s.repo.FindRate(ctx, groupID, to, from, date)
	if err == nil && !reverseRate.IsZero() {
		if result, err := decimal.NewFromString("1"); err == nil {
			return result.Div(reverseRate).Round(6), nil
		}
	}

	// 3. Cross-rate via EUR
	if from != "EUR" && to != "EUR" {
		fromEUR, err := s.repo.FindRate(ctx, groupID, from, "EUR", date)
		if err == nil && !fromEUR.IsZero() {
			toEUR, err := s.repo.FindRate(ctx, groupID, "EUR", to, date)
			if err == nil && !toEUR.IsZero() {
				return fromEUR.Mul(toEUR).Round(6), nil
			}
		}
	}

	return decimal.Zero, fmt.Errorf("no exchange rate found for %s → %s", from, to)
}

// Convert converts an amount from one currency to another.
func (s *ExchangeRateService) Convert(ctx context.Context, groupID uuid.UUID, amount decimal.Decimal, from, to string, date time.Time) (decimal.Decimal, error) {
	if from == to {
		return amount, nil
	}

	rate, err := s.GetRate(ctx, groupID, from, to, date)
	if err != nil {
		return decimal.Zero, err
	}

	return amount.Mul(rate).Round(2), nil
}

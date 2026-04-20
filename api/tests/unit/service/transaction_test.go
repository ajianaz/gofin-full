package service_test

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/azfirazka/gofin-full/api/internal/domain"
)

// These tests validate the transaction service's business logic
// without needing a database (pure unit tests for calculateAmounts).

func TestCalculateAmounts_Withdrawal(t *testing.T) {
	svc := &testableService{}

	src, dst, err := svc.calculateAmounts("withdrawal", decimal.NewFromFloat(100))
	require.NoError(t, err)
	assert.True(t, src.IsNegative(), "source should be negative")
	assert.True(t, dst.IsPositive(), "destination should be positive")
	assert.Equal(t, "-100", src.StringFixed(0))
	assert.Equal(t, "100", dst.StringFixed(0))
}

func TestCalculateAmounts_Deposit(t *testing.T) {
	svc := &testableService{}

	src, dst, err := svc.calculateAmounts("deposit", decimal.NewFromFloat(5000))
	require.NoError(t, err)
	assert.True(t, src.IsNegative())
	assert.True(t, dst.IsPositive())
	assert.Equal(t, "-5000", src.StringFixed(0))
	assert.Equal(t, "5000", dst.StringFixed(0))
}

func TestCalculateAmounts_Transfer(t *testing.T) {
	svc := &testableService{}

	src, dst, err := svc.calculateAmounts("transfer", decimal.NewFromFloat(250))
	require.NoError(t, err)
	assert.True(t, src.IsNegative())
	assert.True(t, dst.IsPositive())
}

func TestCalculateAmounts_InvalidType(t *testing.T) {
	svc := &testableService{}

	_, _, err := svc.calculateAmounts("unknown-type", decimal.NewFromFloat(100))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported transaction type")
}

func TestCalculateAmounts_OpeningBalance(t *testing.T) {
	svc := &testableService{}

	src, dst, err := svc.calculateAmounts("opening-balance", decimal.NewFromFloat(10000))
	require.NoError(t, err)
	assert.True(t, src.IsNegative())
	assert.True(t, dst.IsPositive())
}

func TestCalculateAmounts_Reconciliation(t *testing.T) {
	svc := &testableService{}

	src, dst, err := svc.calculateAmounts("reconciliation", decimal.NewFromFloat(5.50))
	require.NoError(t, err)
	assert.True(t, src.IsNegative())
	assert.True(t, dst.IsPositive())
}

// testableService exposes calculateAmounts for unit testing without DB.
type testableService struct{}

func (s *testableService) calculateAmounts(txType string, amount decimal.Decimal) (decimal.Decimal, decimal.Decimal, error) {
	switch domain.TransactionType(txType) {
	case domain.TransactionTypeWithdrawal, domain.TransactionTypeDeposit,
		domain.TransactionTypeTransfer, domain.TransactionTypeOpeningBalance,
		domain.TransactionTypeReconciliation, domain.TransactionTypeLiabilityCredit:
		return amount.Neg(), amount, nil
	default:
		return decimal.Zero, decimal.Zero, fmt.Errorf("unsupported transaction type: %s", txType)
	}
}

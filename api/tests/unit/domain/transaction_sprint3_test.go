package domain_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/azfirazka/gofin-full/api/internal/domain"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestTransactionType_AllTypes(t *testing.T) {
	types := []domain.TransactionType{
		domain.TransactionTypeWithdrawal,
		domain.TransactionTypeDeposit,
		domain.TransactionTypeTransfer,
		domain.TransactionTypeOpeningBalance,
		domain.TransactionTypeReconciliation,
		domain.TransactionTypeLiabilityCredit,
		domain.TransactionTypeInvalid,
	}
	assert.Len(t, types, 7)
	for _, tt := range types {
		assert.NotEmpty(t, string(tt))
	}
}

func TestTransactionGroup_WithUserGroupID(t *testing.T) {
	group := &domain.TransactionGroup{
		ID:          1,
		UserID:      42,
		UserGroupID: 5,
		GroupTitle:  "Test",
	}

	data, err := json.Marshal(group)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"user_group_id":5`)
	assert.Contains(t, string(data), `"group_title":"Test"`)
}

func TestTransactionJournal_CurrencyIDAsString(t *testing.T) {
	j := &domain.TransactionJournal{
		CurrencyID: "EUR",
	}
	data, err := json.Marshal(j)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"currency_id":"EUR"`)
}

func TestTransactionJournal_ForeignCurrencyNil(t *testing.T) {
	j := &domain.TransactionJournal{
		CurrencyID: "USD",
	}
	data, err := json.Marshal(j)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"currency_id":"USD"`)
	assert.NotContains(t, string(data), `"foreign_currency_id"`)
}

func TestTransactionJournal_ForeignCurrencySet(t *testing.T) {
	fc := "GBP"
	j := &domain.TransactionJournal{
		CurrencyID:        "USD",
		ForeignCurrencyID: &fc,
	}
	data, err := json.Marshal(j)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"foreign_currency_id":"GBP"`)
}

func TestTransactionJournal_OptionalDateFields(t *testing.T) {
	date := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	j := &domain.TransactionJournal{
		Description: "Test",
		InterestDate: &date,
		PaymentDate:  &date,
	}

	data, err := json.Marshal(j)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"interest_date"`)
	assert.Contains(t, string(data), `"payment_date"`)
}

func TestTransaction_NegativeAmount(t *testing.T) {
	tx := &domain.Transaction{
		ID:                   1,
		TransactionJournalID: 100,
		AccountID:            5,
		Amount:               decimal.NewFromFloat(-100.50),
		NativeAmount:         decimal.NewFromFloat(-100.50),
	}

	data, err := json.Marshal(tx)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"amount":"-100.5"`)
}

func TestTransaction_ForeignAmountNil(t *testing.T) {
	tx := &domain.Transaction{
		Amount:       decimal.NewFromFloat(50),
		NativeAmount: decimal.NewFromFloat(50),
	}

	data, err := json.Marshal(tx)
	assert.NoError(t, err)
	assert.NotContains(t, string(data), `"foreign_amount"`)
}

func TestTransaction_ForeignAmountSet(t *testing.T) {
	fa := decimal.NewFromFloat(55)
	tx := &domain.Transaction{
		Amount:            decimal.NewFromFloat(50),
		NativeAmount:      decimal.NewFromFloat(50),
		ForeignAmount:     &fa,
	}

	data, err := json.Marshal(tx)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"foreign_amount":"55"`)
}

func TestTransactionJournalMeta_JSON(t *testing.T) {
	meta := &domain.TransactionJournalMeta{
		ID:                    1,
		TransactionJournalID: 100,
		Name:                  "import_hash",
		Value:                 "abc123",
	}

	data, err := json.Marshal(meta)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"name":"import_hash"`)
	assert.Contains(t, string(data), `"value":"abc123"`)
}

func TestTransactionJournalLink_JSON(t *testing.T) {
	link := &domain.TransactionJournalLink{
		ID:            1,
		LinkTypeID:    1,
		SourceID:      10,
		DestinationID: 20,
	}

	data, err := json.Marshal(link)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"source_id":10`)
	assert.Contains(t, string(data), `"destination_id":20`)
}

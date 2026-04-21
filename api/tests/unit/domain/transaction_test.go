package domain_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ajianaz/gofin-full/api/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestTransactionGroup_JSONSerialization(t *testing.T) {
	group := &domain.TransactionGroup{
		ID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		UserID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		GroupTitle: "Groceries",
		Journals: []domain.TransactionJournal{
			{
				ID:                uuid.MustParse("00000000-0000-0000-0000-000000000064"),
				TransactionGroupID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				UserID:            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				UserGroupID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Type:              domain.TransactionTypeWithdrawal,
				Date:              time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
				Description:       "Weekly groceries",
				CurrencyID:        "EUR",
				Reconciled:        false,
				Tags:              []domain.Tag{{ID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), Tag: "shopping"}},
			},
		},
	}

	data, err := json.Marshal(group)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"group_title":"Groceries"`)
	assert.Contains(t, string(data), `"transactions"`)
	assert.Contains(t, string(data), `"type":"withdrawal"`)
	assert.Contains(t, string(data), `"description":"Weekly groceries"`)
}

func TestTransactionGroup_NoJournals(t *testing.T) {
	group := &domain.TransactionGroup{
		ID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		UserID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		GroupTitle: "Empty group",
	}

	data, err := json.Marshal(group)
	assert.NoError(t, err)
	// Journals should be omitted when empty due to omitempty
	assert.NotContains(t, string(data), `"transactions"`)
}

func TestTransactionJournal_MetaFields(t *testing.T) {
	journal := &domain.TransactionJournal{
		ID:          uuid.MustParse("00000000-0000-0000-0000-000000000064"),
		Description: "SEPA transfer",
		SepaCC:      ptrStr("0000"),
		SepaCTID:    ptrStr("ID123"),
		ExternalID:  ptrStr("EXT-001"),
	}

	data, err := json.Marshal(journal)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"sepa_cc":"0000"`)
	assert.Contains(t, string(data), `"sepa_ct_id":"ID123"`)
	assert.Contains(t, string(data), `"external_id":"EXT-001"`)
}

func TestTransaction_AmountSerialization(t *testing.T) {
	tx := &domain.Transaction{
		ID:                   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		TransactionJournalID: uuid.MustParse("00000000-0000-0000-0000-000000000064"),
		AccountID:            uuid.MustParse("00000000-0000-0000-0000-000000000005"),
		Amount:               decimal.NewFromFloat(42.50),
		NativeAmount:         decimal.NewFromFloat(42.50),
	}

	data, err := json.Marshal(tx)
	assert.NoError(t, err)
	// shopspring/decimal strips trailing zeros: 42.50 -> "42.5"
	assert.Contains(t, string(data), `"amount":"42.5"`)
	assert.Contains(t, string(data), `"native_amount":"42.5"`)
}

func TestTransactionType_Values(t *testing.T) {
	assert.Equal(t, domain.TransactionType("withdrawal"), domain.TransactionTypeWithdrawal)
	assert.Equal(t, domain.TransactionType("deposit"), domain.TransactionTypeDeposit)
	assert.Equal(t, domain.TransactionType("transfer"), domain.TransactionTypeTransfer)
	assert.Equal(t, domain.TransactionType("opening-balance"), domain.TransactionTypeOpeningBalance)
	assert.Equal(t, domain.TransactionType("reconciliation"), domain.TransactionTypeReconciliation)
	assert.Equal(t, domain.TransactionType("liability-credit"), domain.TransactionTypeLiabilityCredit)
	assert.Equal(t, domain.TransactionType("invalid"), domain.TransactionTypeInvalid)
}

func ptrStr(v string) *string { return &v }

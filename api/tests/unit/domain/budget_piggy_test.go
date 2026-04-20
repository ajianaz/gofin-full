package domain_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/azfirazka/gofin-full/api/internal/domain"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestBudget_JSONSerialization(t *testing.T) {
	b := &domain.Budget{
		ID:         1,
		UserID:     42,
		UserGroupID: 5,
		Name:       "Groceries",
		Active:     true,
		Order:      1,
	}

	data, err := json.Marshal(b)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"id":1`)
	assert.Contains(t, string(data), `"name":"Groceries"`)
	assert.Contains(t, string(data), `"active":true`)
	assert.NotContains(t, string(data), `"deleted_at"`)
}

func TestBudget_Limits(t *testing.T) {
	b := &domain.Budget{
		ID: 1, Name: "Food",
		Limits: []domain.BudgetLimit{
			{
				ID:       10,
				BudgetID: 1,
				Start:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				End:      time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
				Amount:   decimal.NewFromFloat(500),
			},
		},
	}

	data, err := json.Marshal(b)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"limits"`)
	assert.Contains(t, string(data), `"amount":"500"`)
}

func TestBudget_NoLimits(t *testing.T) {
	b := &domain.Budget{ID: 1, Name: "Empty"}
	data, err := json.Marshal(b)
	assert.NoError(t, err)
	assert.NotContains(t, string(data), `"limits"`)
}

func TestAutoBudgetType_Values(t *testing.T) {
	assert.Equal(t, domain.AutoBudgetType("none"), domain.AutoBudgetTypeNone)
	assert.Equal(t, domain.AutoBudgetType("reset"), domain.AutoBudgetTypeReset)
	assert.Equal(t, domain.AutoBudgetType("rollover"), domain.AutoBudgetTypeRollover)
	assert.Equal(t, domain.AutoBudgetType("adjusted"), domain.AutoBudgetTypeAdjusted)
}

func TestPiggyBank_JSONSerialization(t *testing.T) {
	pb := &domain.PiggyBank{
		ID:           1,
		AccountID:    5,
		Name:         "Vacation Fund",
		TargetAmount: decimal.NewFromFloat(5000),
		CurrentAmount: decimal.NewFromFloat(1200),
		LeftToTarget:  decimal.NewFromFloat(3800),
		Percentage:    24.0,
		Order:        1,
	}

	data, err := json.Marshal(pb)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"name":"Vacation Fund"`)
	assert.Contains(t, string(data), `"target_amount":"5000"`)
	assert.Contains(t, string(data), `"current_amount":"1200"`)
	assert.Contains(t, string(data), `"left_to_target":"3800"`)
	assert.Contains(t, string(data), `"percentage":24`)
}

func TestPiggyBank_OptionalFields(t *testing.T) {
	pb := &domain.PiggyBank{
		ID: 1, Name: "Test",
	}
	data, err := json.Marshal(pb)
	assert.NoError(t, err)
	assert.NotContains(t, string(data), `"start_date"`)
	assert.NotContains(t, string(data), `"target_date"`)
	assert.NotContains(t, string(data), `"notes"`)
}

func TestPiggyBank_WithDates(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	target := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	notes := "Save for holiday"

	pb := &domain.PiggyBank{
		ID: 1, Name: "Holiday",
		StartDate: &start, TargetDate: &target, Notes: &notes,
	}

	data, err := json.Marshal(pb)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"start_date"`)
	assert.Contains(t, string(data), `"target_date"`)
	assert.Contains(t, string(data), `"notes":"Save for holiday"`)
}

func TestPiggyBankEvent_JSON(t *testing.T) {
	evt := &domain.PiggyBankEvent{
		ID:          1,
		PiggyBankID: 5,
		Amount:      decimal.NewFromFloat(250),
	}

	data, err := json.Marshal(evt)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"piggy_bank_id":5`)
	assert.Contains(t, string(data), `"amount":"250"`)
}

func TestPiggyBankRepetition_JSON(t *testing.T) {
	rep := &domain.PiggyBankRepetition{
		ID:            1,
		PiggyBankID:   5,
		TargetAmount:  decimal.NewFromFloat(1000),
		CurrentAmount: decimal.NewFromFloat(500),
	}

	data, err := json.Marshal(rep)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"target_amount":"1000"`)
	assert.Contains(t, string(data), `"current_amount":"500"`)
}

func TestAvailableBudget_JSON(t *testing.T) {
	ab := &domain.AvailableBudget{
		ID:         1,
		BudgetID:   10,
		Amount:     decimal.NewFromFloat(300),
	}

	data, err := json.Marshal(ab)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"amount":"300"`)
}

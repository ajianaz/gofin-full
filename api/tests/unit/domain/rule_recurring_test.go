package domain_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/azfirazka/gofin-full/api/internal/domain"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestRuleGroup_JSON(t *testing.T) {
	rg := &domain.RuleGroup{
		ID: 1, Title: "Auto-categorize", Active: true, Order: 1,
	}
	data, err := json.Marshal(rg)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"title":"Auto-categorize"`)
	assert.Contains(t, string(data), `"active":true`)
	assert.NotContains(t, string(data), `"deleted_at"`)
}

func TestRule_WithTriggersAndActions(t *testing.T) {
	rule := &domain.Rule{
		ID: 1, Title: "Set category for groceries",
		Active: true, Strict: true, StopProcessing: false,
		Triggers: []domain.RuleTrigger{
			{ID: 1, TriggerType: "description_contains", TriggerValue: "supermarket", StopProcessing: false},
		},
		Actions: []domain.RuleAction{
			{ID: 1, ActionType: "set_category", ActionValue: "Groceries", Order: 1},
		},
	}

	data, err := json.Marshal(rule)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"triggers"`)
	assert.Contains(t, string(data), `"actions"`)
	assert.Contains(t, string(data), `"description_contains"`)
	assert.Contains(t, string(data), `"set_category"`)
}

func TestRule_NoTriggersOrActions(t *testing.T) {
	rule := &domain.Rule{ID: 1, Title: "Empty"}
	data, err := json.Marshal(rule)
	assert.NoError(t, err)
	assert.NotContains(t, string(data), `"triggers"`)
	assert.NotContains(t, string(data), `"actions"`)
}

func TestRuleTrigger_JSON(t *testing.T) {
	tr := &domain.RuleTrigger{
		ID: 1, TriggerType: "amount_is", TriggerValue: ">100", StopProcessing: true,
	}
	data, err := json.Marshal(tr)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"trigger_type":"amount_is"`)
	assert.Contains(t, string(data), `"stop_processing":true`)
}

func TestRuleAction_JSON(t *testing.T) {
	act := &domain.RuleAction{
		ID: 1, ActionType: "convert_amount", ActionValue: "EUR", Order: 2,
	}
	data, err := json.Marshal(act)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"action_type":"convert_amount"`)
	assert.Contains(t, string(data), `"order":2`)
}

func TestRecurrence_JSON(t *testing.T) {
	rec := &domain.Recurrence{
		ID: 1, Title: "Monthly salary",
		FirstDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		RepeatFreq: "monthly",
		Active: true,
	}
	data, err := json.Marshal(rec)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"title":"Monthly salary"`)
	assert.Contains(t, string(data), `"repeat_freq":"monthly"`)
	assert.NotContains(t, string(data), `"deleted_at"`)
}

func TestRecurrence_WithTransactions(t *testing.T) {
	rec := &domain.Recurrence{
		ID: 1, Title: "Rent",
		Transactions: []domain.RecurringTransaction{
			{
				ID: 1, Type: "withdrawal", Description: "Monthly rent",
				Amount: decimal.NewFromFloat(1500), SourceID: 1, DestinationID: 2, Order: 1,
			},
		},
	}
	data, err := json.Marshal(rec)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"transactions"`)
	assert.Contains(t, string(data), `"amount":"1500"`)
}

func TestRecurrence_OptionalFields(t *testing.T) {
	rec := &domain.Recurrence{ID: 1, Title: "Test"}
	data, err := json.Marshal(rec)
	assert.NoError(t, err)
	assert.NotContains(t, string(data), `"description"`)
	assert.NotContains(t, string(data), `"repeat_until"`)
	assert.NotContains(t, string(data), `"transactions"`)
}

func TestRecurringTransaction_JSON(t *testing.T) {
	tx := &domain.RecurringTransaction{
		ID: 1, Type: "deposit", Description: "Salary",
		Amount: decimal.NewFromFloat(3000), SourceID: 10, DestinationID: 5,
		BudgetID: ptrInt64(3), Order: 1,
	}
	data, err := json.Marshal(tx)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"budget_id":3`)
	assert.NotContains(t, string(data), `"category_id"`)
}

func TestRecurringRepetition_JSON(t *testing.T) {
	rep := &domain.RecurringRepetition{
		ID: 1, RecurrenceID: 10,
		RelevantDate: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
	}
	data, err := json.Marshal(rep)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"relevant_date"`)
}

func TestRecurrenceMeta_JSON(t *testing.T) {
	meta := &domain.RecurrenceMeta{
		ID: 1, Name: "notes", Value: "remember to check",
	}
	data, err := json.Marshal(meta)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"name":"notes"`)
	assert.Contains(t, string(data), `"value":"remember to check"`)
}

func ptrInt64(v int64) *int64 { return &v }

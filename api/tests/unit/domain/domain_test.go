package domain_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ajianaz/gofin-full/api/internal/domain"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestUser_JSONTags(t *testing.T) {
	groupID := int64(1)
	user := &domain.User{
		ID:          1,
		Email:       "test@example.com",
		Blocked:     false,
		UserGroupID: &groupID,
	}

	data, err := json.Marshal(user)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"id":1`)
	assert.Contains(t, string(data), `"email":"test@example.com"`)
	// Password should not be serialized
	assert.NotContains(t, string(data), "password")
}

func TestUser_PasswordNotInJSON(t *testing.T) {
	user := &domain.User{
		ID:       1,
		Email:    "test@example.com",
		Password: "super-secret",
	}

	data, err := json.Marshal(user)
	assert.NoError(t, err)
	assert.NotContains(t, string(data), "super-secret")
	assert.NotContains(t, string(data), "password")
}

func TestUser_DeletedAtNotInJSON(t *testing.T) {
	now := time.Now()
	user := &domain.User{
		ID:        1,
		DeletedAt: &now,
	}

	data, err := json.Marshal(user)
	assert.NoError(t, err)
	assert.NotContains(t, string(data), "deleted_at")
}

func TestRole_JSONTags(t *testing.T) {
	role := &domain.Role{
		ID:    1,
		Title: "owner",
	}

	data, err := json.Marshal(role)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"id":1`)
	assert.Contains(t, string(data), `"title":"owner"`)
}

func TestUserGroup_JSONTags(t *testing.T) {
	group := &domain.UserGroup{
		ID:    1,
		Title: "Personal",
	}

	data, err := json.Marshal(group)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"title":"Personal"`)
}

func TestCurrency_JSONTags(t *testing.T) {
	currency := &domain.Currency{
		ID:            1,
		Code:          "USD",
		Name:          "US Dollar",
		Symbol:        "$",
		DecimalPlaces: 2,
		Enabled:       true,
	}

	data, err := json.Marshal(currency)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"code":"USD"`)
	assert.Contains(t, string(data), `"symbol":"$"`)
	assert.Contains(t, string(data), `"decimal_places":2`)
}

func TestBudget_JSONTags(t *testing.T) {
	budget := &domain.Budget{
		ID:     1,
		UserID: 1,
		Name:   "Monthly Budget",
		Active: true,
	}

	data, err := json.Marshal(budget)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"name":"Monthly Budget"`)
	assert.Contains(t, string(data), `"active":true`)
}

func TestBudgetLimit_JSONTags(t *testing.T) {
	limit := &domain.BudgetLimit{
		ID:       1,
		BudgetID: 1,
		Amount:   decimal.NewFromFloat(500.00),
	}

	data, err := json.Marshal(limit)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"amount":"500"`)
}

func TestBill_JSONTags(t *testing.T) {
	bill := &domain.Bill{
		ID:       1,
		Name:     "Electricity",
		AmountMin: decimal.NewFromFloat(50.00),
		AmountMax: decimal.NewFromFloat(100.00),
		Active:   true,
	}

	data, err := json.Marshal(bill)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"name":"Electricity"`)
	assert.Contains(t, string(data), `"amount_min":"50"`)
	assert.Contains(t, string(data), `"amount_max":"100"`)
}

func TestPiggyBank_JSONTags(t *testing.T) {
	pb := &domain.PiggyBank{
		ID:           1,
		AccountID:    1,
		Name:         "Vacation Fund",
		TargetAmount: decimal.NewFromFloat(5000.00),
	}

	data, err := json.Marshal(pb)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"name":"Vacation Fund"`)
	assert.Contains(t, string(data), `"target_amount":"5000"`)
}

func TestCategory_JSONTags(t *testing.T) {
	cat := &domain.Category{
		ID:    1,
		Name:  "Groceries",
	}

	data, err := json.Marshal(cat)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"name":"Groceries"`)
}

func TestTag_JSONTags(t *testing.T) {
	tag := &domain.Tag{
		ID:  1,
		Tag: "shopping",
	}

	data, err := json.Marshal(tag)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"tag":"shopping"`)
}

func TestRule_JSONTags(t *testing.T) {
	rule := &domain.Rule{
		ID:              1,
		Title:           "Auto-categorize groceries",
		Active:          true,
		Strict:          false,
		StopProcessing:  false,
		Triggers: []domain.RuleTrigger{
			{ID: 1, TriggerType: "description_contains", TriggerValue: "walmart"},
		},
		Actions: []domain.RuleAction{
			{ID: 1, ActionType: "set_category", ActionValue: "3"},
		},
	}

	data, err := json.Marshal(rule)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"title":"Auto-categorize groceries"`)
	assert.Contains(t, string(data), `"triggers"`)
	assert.Contains(t, string(data), `"actions"`)
	assert.Contains(t, string(data), `"trigger_type":"description_contains"`)
	assert.Contains(t, string(data), `"action_type":"set_category"`)
}

func TestWebhook_JSONTags(t *testing.T) {
	wh := &domain.Webhook{
		ID:     1,
		Title:  "Slack notification",
		URL:    "https://hooks.slack.com/services/xxx",
		Active: true,
	}

	data, err := json.Marshal(wh)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"url":"https://hooks.slack.com/services/xxx"`)
}

func TestAttachment_JSONTags(t *testing.T) {
	att := &domain.Attachment{
		ID:             1,
		Filename:       "receipt.pdf",
		MimeType:       "application/pdf",
		Size:           1024,
		AttachableType: "TransactionJournal",
		AttachableID:   100,
	}

	data, err := json.Marshal(att)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"filename":"receipt.pdf"`)
	assert.Contains(t, string(data), `"mime_type":"application/pdf"`)
}

func TestRecurringTransaction_JSONTags(t *testing.T) {
	rt := &domain.Recurrence{
		ID:          1,
		Title:       "Monthly rent",
		RepeatFreq:  "monthly",
		Active:      true,
		ApplyRules:  true,
	}

	data, err := json.Marshal(rt)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"title":"Monthly rent"`)
	assert.Contains(t, string(data), `"repeat_freq":"monthly"`)
}

func TestExchangeRate_JSONTags(t *testing.T) {
	rate := &domain.ExchangeRate{
		ID:             1,
		FromCurrencyID: "USD",
		ToCurrencyID:   "EUR",
		Rate:           decimal.NewFromFloat(1.13),
	}

	data, err := json.Marshal(rate)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"rate":"1.13"`)
}

func TestNotification_JSONTags(t *testing.T) {
	n := &domain.Notification{
		ID:      1,
		Channel: "email",
		Type:    "bill_reminder",
		Title:   "Bill Due",
		Message: "Your electricity bill is due tomorrow",
		Read:    false,
	}

	data, err := json.Marshal(n)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"type":"bill_reminder"`)
	assert.Contains(t, string(data), `"read":false`)
}

func TestPreference_JSONTags(t *testing.T) {
	p := &domain.Preference{
		ID:    1,
		Name:  "currencyPreference",
		Data:  "EUR",
	}

	data, err := json.Marshal(p)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"name":"currencyPreference"`)
}

package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/ajianaz/gofin-full/api/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestWalletType_SourceValidation(t *testing.T) {
	tests := []struct {
		wt       domain.WalletType
		expected bool
	}{
		{domain.WalletTypeAsset, true},
		{domain.WalletTypeDefault, true},
		{domain.WalletTypeCash, true},
		{domain.WalletTypeDebt, true},
		{domain.WalletTypeInitialBalance, true},
		{domain.WalletTypeLoan, true},
		{domain.WalletTypeMortgage, true},
		{domain.WalletTypeReconciliation, true},
		{domain.WalletTypeExpense, false},
		{domain.WalletTypeRevenue, false},
		{domain.WalletTypeBeneficiary, false},
		{domain.WalletTypeCreditCard, false},
		{domain.WalletTypeImport, false},
		{domain.WalletTypeLiabilityCredit, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.wt), func(t *testing.T) {
			assert.Equal(t, tt.expected, domain.IsSourceValid(tt.wt))
		})
	}
}

func TestWalletType_DestinationValidation(t *testing.T) {
	tests := []struct {
		wt       domain.WalletType
		expected bool
	}{
		{domain.WalletTypeAsset, true},
		{domain.WalletTypeDefault, true},
		{domain.WalletTypeCash, true},
		{domain.WalletTypeDebt, true},
		{domain.WalletTypeInitialBalance, true},
		{domain.WalletTypeLoan, true},
		{domain.WalletTypeMortgage, true},
		{domain.WalletTypeExpense, true},
		{domain.WalletTypeRevenue, false},
		{domain.WalletTypeBeneficiary, false},
		{domain.WalletTypeCreditCard, false},
		{domain.WalletTypeImport, false},
		{domain.WalletTypeLiabilityCredit, false},
		{domain.WalletTypeReconciliation, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.wt), func(t *testing.T) {
			assert.Equal(t, tt.expected, domain.IsDestinationValid(tt.wt))
		})
	}
}

func TestWalletType_CanHoldPiggyBanks(t *testing.T) {
	tests := []struct {
		wt       domain.WalletType
		expected bool
	}{
		{domain.WalletTypeAsset, true},
		{domain.WalletTypeDefault, true},
		{domain.WalletTypeLoan, true},
		{domain.WalletTypeMortgage, true},
		{domain.WalletTypeDebt, true},
		{domain.WalletTypeExpense, false},
		{domain.WalletTypeRevenue, false},
		{domain.WalletTypeCash, false},
		{domain.WalletTypeCreditCard, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.wt), func(t *testing.T) {
			assert.Equal(t, tt.expected, domain.CanHoldPiggyBanks(tt.wt))
		})
	}
}

func TestWalletType_CanHaveOpeningBalance(t *testing.T) {
	tests := []struct {
		wt       domain.WalletType
		expected bool
	}{
		{domain.WalletTypeAsset, true},
		{domain.WalletTypeLoan, true},
		{domain.WalletTypeMortgage, true},
		{domain.WalletTypeDebt, true},
		{domain.WalletTypeExpense, false},
		{domain.WalletTypeRevenue, false},
		{domain.WalletTypeCash, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.wt), func(t *testing.T) {
			assert.Equal(t, tt.expected, domain.CanHaveOpeningBalance(tt.wt))
		})
	}
}

func TestWalletType_CanHaveCurrency(t *testing.T) {
	tests := []struct {
		wt       domain.WalletType
		expected bool
	}{
		{domain.WalletTypeAsset, true},
		{domain.WalletTypeCash, true},
		{domain.WalletTypeCreditCard, true},
		{domain.WalletTypeExpense, false},
		{domain.WalletTypeRevenue, false},
		{domain.WalletTypeBeneficiary, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.wt), func(t *testing.T) {
			assert.Equal(t, tt.expected, domain.CanHaveCurrency(tt.wt))
		})
	}
}

func TestWallet_JSONTags(t *testing.T) {
	w := &domain.Wallet{
		ID:             uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		UserID:         uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		UserGroupID:    uuid.MustParse("00000000-0000-0000-0000-000000000003"),
		Name:           "Main Account",
		Active:         true,
		VirtualBalance: decimal.NewFromFloat(1000.50),
		IncludeNetWorth: true,
	}

	data, err := json.Marshal(w)
	assert.NoError(t, err)

	// Verify key JSON fields are present
	assert.Contains(t, string(data), `"name":"Main Account"`)
	assert.Contains(t, string(data), `"active":true`)
	// shopspring/decimal strips trailing zeros: 1000.50 -> "1000.5"
	assert.Contains(t, string(data), `"virtual_balance":"1000.5"`)
	assert.Contains(t, string(data), `"include_net_worth":true`)

	// Verify password-like fields are not in JSON
	assert.NotContains(t, string(data), "deleted_at")
}

func TestWallet_OptionalFields(t *testing.T) {
	iban := "NL91ABNA0417164300"
	bic := "ABNANL2A"
	w := &domain.Wallet{
		IBAN:      &iban,
		BIC:       &bic,
		Latitude:  ptrFloat64(52.3676),
		Longitude: ptrFloat64(4.9041),
	}

	data, err := json.Marshal(w)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"iban":"NL91ABNA0417164300"`)
	assert.Contains(t, string(data), `"bic":"ABNANL2A"`)
}

func TestWallet_LiabilityFields(t *testing.T) {
	liabType := "loan"
	liabDir := "credit"
	interest := decimal.NewFromFloat(3.5)
	w := &domain.Wallet{
		LiabilityType:      &liabType,
		LiabilityDirection: &liabDir,
		InterestRate:       &interest,
		InterestPeriod:     ptrStr("monthly"),
	}

	data, err := json.Marshal(w)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"liability_type":"loan"`)
	assert.Contains(t, string(data), `"liability_direction":"credit"`)
	assert.Contains(t, string(data), `"interest":"3.5"`)
}

func TestWalletMember_RolePermissions(t *testing.T) {
	tests := []struct {
		role     string
		isOwner  bool
		canWrite bool
		canShare bool
	}{
		{"owner", true, true, true},
		{"editor", false, true, false},
		{"viewer", false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.role, func(t *testing.T) {
			m := &domain.WalletMember{Role: tt.role}
			assert.Equal(t, tt.isOwner, m.IsOwner())
			assert.Equal(t, tt.canWrite, m.CanWrite())
			assert.Equal(t, tt.canShare, m.CanShare())
		})
	}
}

func ptrFloat64(v float64) *float64 { return &v }

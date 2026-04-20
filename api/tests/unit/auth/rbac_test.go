package auth_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/azfirazka/gofin-full/api/internal/auth"
)

func TestRBAC_HasPermission_OwnerHasAll(t *testing.T) {
	roles := auth.AllGroupRoles()
	for _, required := range roles {
		assert.True(t, auth.HasPermission(auth.RoleOwner, required),
			"owner should have permission: %s", required)
	}
}

func TestRBAC_HasPermission_FullHasAll(t *testing.T) {
	roles := auth.AllGroupRoles()
	for _, required := range roles {
		if required == auth.RoleOwner {
			continue // full doesn't have owner
		}
		assert.True(t, auth.HasPermission(auth.RoleFull, required),
			"full should have permission: %s", required)
	}
}

func TestRBAC_HasPermission_ReadOnlyMinimal(t *testing.T) {
	assert.True(t, auth.HasPermission(auth.RoleReadOnly, auth.RoleReadOnly))
	assert.False(t, auth.HasPermission(auth.RoleReadOnly, auth.RoleManageTransactions))
}

func TestRBAC_HasPermission_ManageTransactions(t *testing.T) {
	assert.True(t, auth.HasPermission(auth.RoleManageTransactions, auth.RoleReadOnly))
	assert.True(t, auth.HasPermission(auth.RoleManageTransactions, auth.RoleManageTransactions))
	assert.False(t, auth.HasPermission(auth.RoleManageTransactions, auth.RoleManageMeta))
}

func TestRBAC_HasPermission_Hierarchy(t *testing.T) {
	// read_only < manage_transactions < manage_meta < read_budgets < ...
	hierarchy := []auth.GroupRole{
		auth.RoleReadOnly,
		auth.RoleManageTransactions,
		auth.RoleManageMeta,
		auth.RoleReadBudgets,
		auth.RoleManageBudgets,
		auth.RoleReadPiggyBanks,
		auth.RoleManagePiggyBanks,
		auth.RoleReadSubscriptions,
		auth.RoleManageSubscriptions,
		auth.RoleReadRules,
		auth.RoleManageRules,
		auth.RoleReadRecurring,
		auth.RoleManageRecurring,
		auth.RoleReadWebhooks,
		auth.RoleManageWebhooks,
		auth.RoleReadCurrencies,
		auth.RoleManageCurrencies,
		auth.RoleViewReports,
		auth.RoleViewMemberships,
		auth.RoleFull,
		auth.RoleOwner,
	}

	for i, role := range hierarchy {
		for j := 0; j <= i; j++ {
			assert.True(t, auth.HasPermission(role, hierarchy[j]),
				"%s should have %s", role, hierarchy[j])
		}
	}
}

func TestRBAC_IsOwnerOrFull(t *testing.T) {
	assert.True(t, auth.IsOwnerOrFull(auth.RoleOwner))
	assert.True(t, auth.IsOwnerOrFull(auth.RoleFull))
	assert.False(t, auth.IsOwnerOrFull(auth.RoleReadOnly))
	assert.False(t, auth.IsOwnerOrFull(auth.RoleManageTransactions))
}

func TestRBAC_AllGroupRoles_Count(t *testing.T) {
	roles := auth.AllGroupRoles()
	assert.Len(t, roles, 21) // 21 group-level roles (demo is global-only)
}

func TestRBAC_AllGroupRoles_ContainsAll(t *testing.T) {
	roles := auth.AllGroupRoles()
	// Verify key roles exist
	assert.Contains(t, roles, auth.RoleReadOnly)
	assert.Contains(t, roles, auth.RoleOwner)
	assert.Contains(t, roles, auth.RoleFull)
	assert.Contains(t, roles, auth.RoleManageTransactions)
	assert.Contains(t, roles, auth.RoleManageBudgets)
	assert.Contains(t, roles, auth.RoleViewReports)
	assert.Contains(t, roles, auth.RoleViewMemberships)
}

package auth

import (
	"github.com/gofiber/fiber/v2"

	apperrors "github.com/azfirazka/gofin-full/api/pkg/errors"
)

// GroupRole represents the 22 group-level permission values from Firefly III.
type GroupRole string

const (
	RoleReadOnly            GroupRole = "read_only"
	RoleManageTransactions  GroupRole = "manage_transactions"
	RoleManageMeta          GroupRole = "manage_meta"
	RoleReadBudgets         GroupRole = "read_budgets"
	RoleManageBudgets       GroupRole = "manage_budgets"
	RoleReadPiggyBanks      GroupRole = "read_piggy_banks"
	RoleManagePiggyBanks    GroupRole = "manage_piggy_banks"
	RoleReadSubscriptions   GroupRole = "read_subscriptions"
	RoleManageSubscriptions GroupRole = "manage_subscriptions"
	RoleReadRules           GroupRole = "read_rules"
	RoleManageRules         GroupRole = "manage_rules"
	RoleReadRecurring       GroupRole = "read_recurring"
	RoleManageRecurring     GroupRole = "manage_recurring"
	RoleReadWebhooks        GroupRole = "read_webhooks"
	RoleManageWebhooks      GroupRole = "manage_webhooks"
	RoleReadCurrencies      GroupRole = "read_currencies"
	RoleManageCurrencies    GroupRole = "manage_currencies"
	RoleViewReports         GroupRole = "view_reports"
	RoleViewMemberships     GroupRole = "view_memberships"
	RoleFull                GroupRole = "full"
	RoleOwner               GroupRole = "owner"
)

// roleHierarchy defines the permission hierarchy.
// Each role implicitly includes all roles below it.
var roleHierarchy = map[GroupRole]int{
	RoleReadOnly:            1,
	RoleManageTransactions:  2,
	RoleManageMeta:          3,
	RoleReadBudgets:         4,
	RoleManageBudgets:       5,
	RoleReadPiggyBanks:      6,
	RoleManagePiggyBanks:    7,
	RoleReadSubscriptions:   8,
	RoleManageSubscriptions: 9,
	RoleReadRules:           10,
	RoleManageRules:         11,
	RoleReadRecurring:       12,
	RoleManageRecurring:     13,
	RoleReadWebhooks:        14,
	RoleManageWebhooks:      15,
	RoleReadCurrencies:      16,
	RoleManageCurrencies:    17,
	RoleViewReports:         18,
	RoleViewMemberships:     19,
	RoleFull:                20,
	RoleOwner:               21,
}

// HasPermission checks if a user's role satisfies the required role.
// FULL and OWNER cascade down to all roles.
func HasPermission(userRole, requiredRole GroupRole) bool {
	userLevel := roleHierarchy[userRole]
	requiredLevel := roleHierarchy[requiredRole]
	return userLevel >= requiredLevel
}

// IsOwnerOrFull checks if the role is owner or full (highest privileges).
func IsOwnerOrFull(role GroupRole) bool {
	return role == RoleOwner || role == RoleFull
}

// RBACMiddleware creates a middleware that checks group-level permissions.
// It reads the user's role in the active group from the context (set by previous middleware).
func RBACMiddleware(requiredRole GroupRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRoleVal := c.Locals("user_group_role")
		if userRoleVal == nil {
			return apperrors.ErrForbidden
		}

		userRole, ok := userRoleVal.(GroupRole)
		if !ok {
			return apperrors.ErrForbidden
		}

		if !HasPermission(userRole, requiredRole) {
			return apperrors.ErrForbidden
		}

		return c.Next()
	}
}

// AdminMiddleware checks for global owner role.
func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		isAdminVal := c.Locals("is_admin")
		if isAdminVal == nil {
			return apperrors.ErrForbidden
		}

		isAdmin, ok := isAdminVal.(bool)
		if !ok || !isAdmin {
			return apperrors.ErrForbidden
		}

		return c.Next()
	}
}

// AllGroupRoles returns all 22 group-level roles.
func AllGroupRoles() []GroupRole {
	return []GroupRole{
		RoleReadOnly, RoleManageTransactions, RoleManageMeta,
		RoleReadBudgets, RoleManageBudgets,
		RoleReadPiggyBanks, RoleManagePiggyBanks,
		RoleReadSubscriptions, RoleManageSubscriptions,
		RoleReadRules, RoleManageRules,
		RoleReadRecurring, RoleManageRecurring,
		RoleReadWebhooks, RoleManageWebhooks,
		RoleReadCurrencies, RoleManageCurrencies,
		RoleViewReports, RoleViewMemberships,
		RoleFull, RoleOwner,
	}
}

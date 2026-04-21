package auth

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RoleLookup is implemented by repositories that can resolve a user's group role.
type RoleLookup interface {
	GetUserRoleInGroup(ctx context.Context, userID, groupID uuid.UUID) (GroupRole, error)
	HasGlobalRole(ctx context.Context, userID uuid.UUID, roleTitle string) (bool, error)
}

// GroupRoleMiddleware looks up the authenticated user's role in their active group
// and sets c.Locals("user_group_role") so RBACMiddleware can enforce permissions.
// It also checks global admin role and sets c.Locals("is_admin").
func GroupRoleMiddleware(roleLookup RoleLookup) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := GetUser(c)
		if user == nil {
			return c.Next()
		}

		groupID := GetActiveGroupID(c)
		if groupID == nil || *groupID == uuid.Nil {
			return c.Next()
		}

		role, err := roleLookup.GetUserRoleInGroup(c.Context(), user.ID, *groupID)
		if err != nil {
			// User has no membership in this group — leave role unset,
			// RBACMiddleware will reject with 403.
			return c.Next()
		}

		c.Locals("user_group_role", role)

		// Check global admin
		isAdmin, _ := roleLookup.HasGlobalRole(c.Context(), user.ID, "owner")
		c.Locals("is_admin", isAdmin)

		return c.Next()
	}
}

// SetGroupRoleForTest sets the user's group role directly in context (for testing).
func SetGroupRoleForTest(c *fiber.Ctx, role GroupRole) {
	c.Locals("user_group_role", role)
}

// SetIsAdminForTest sets the admin flag directly in context (for testing).
func SetIsAdminForTest(c *fiber.Ctx, isAdmin bool) {
	c.Locals("is_admin", isAdmin)
}

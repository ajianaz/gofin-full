package auth

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// authLog is the package-level zerolog logger for the auth package.
var authLog = zerolog.New(os.Stderr).With().Timestamp().Logger()

// RoleLookup is implemented by repositories that can resolve a user's group role.
type RoleLookup interface {
	GetUserRoleInGroup(ctx context.Context, userID, groupID uuid.UUID) (GroupRole, error)
	HasGlobalRole(ctx context.Context, userID uuid.UUID, roleTitle string) (bool, error)
	HasAnyGlobalRole(ctx context.Context, userID uuid.UUID, roleTitles ...string) (bool, error)
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
			// DB error — fail closed to prevent unauthorised access.
			authLog.Error().Err(err).Str("user_id", user.ID.String()).Str("group_id", groupID.String()).Msg("GroupRoleMiddleware: failed to lookup role")
			return c.Status(401).JSON(fiber.Map{
				"status": 401,
				"title":  "Unable to verify permissions",
			})
		}

		c.Locals("user_group_role", role)

		// Check global admin (owner or admin)
		isAdmin, _ := roleLookup.HasAnyGlobalRole(c.Context(), user.ID, "owner", "admin")
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

// InvalidateTokenCache invalidates the token version cache for a user if
// an invalidator is available in context (set by router middleware).
func InvalidateTokenCache(c *fiber.Ctx, userID uuid.UUID) {
	if inv, ok := c.Locals("token_invalidator").(TokenInvalidator); ok {
		inv.InvalidateTokenVersion(c.Context(), userID)
	}
}

package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// contextKey is an unexported type for Fiber context keys.
type contextKey string

const (
	userKey    contextKey = "user"
	claimsKey  contextKey = "claims"
	groupIDKey contextKey = "group_id"
)

// SetUser stores the user identity in the Fiber context.
func SetUser(c *fiber.Ctx, user *UserIdentity) {
	c.Locals(string(userKey), user)
}

// GetUser retrieves the user identity from the Fiber context.
func GetUser(c *fiber.Ctx) *UserIdentity {
	val := c.Locals(string(userKey))
	if val == nil {
		return nil
	}
	user, ok := val.(*UserIdentity)
	if !ok {
		return nil
	}
	return user
}

// SetClaims stores the JWT claims in the Fiber context.
func SetClaims(c *fiber.Ctx, claims *Claims) {
	c.Locals(string(claimsKey), claims)
}

// GetClaims retrieves the JWT claims from the Fiber context.
func GetClaims(c *fiber.Ctx) *Claims {
	val := c.Locals(string(claimsKey))
	if val == nil {
		return nil
	}
	claims, ok := val.(*Claims)
	if !ok {
		return nil
	}
	return claims
}

// SetActiveGroupID stores the active group ID in the Fiber context.
func SetActiveGroupID(c *fiber.Ctx, groupID uuid.UUID) {
	c.Locals(string(groupIDKey), groupID)
}

// GetActiveGroupID retrieves the active group ID from the Fiber context.
func GetActiveGroupID(c *fiber.Ctx) *uuid.UUID {
	val := c.Locals(string(groupIDKey))
	if val == nil {
		return nil
	}
	id, ok := val.(uuid.UUID)
	if !ok {
		return nil
	}
	return &id
}

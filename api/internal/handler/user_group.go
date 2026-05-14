package handler

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/repository"
	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors"
)

// UserGroupHandler handles user group endpoints.
type UserGroupHandler struct {
	groupRepo *repository.UserGroupRepository
	userRepo  *repository.UserRepository
	db        *pgxpool.Pool
	jwtMgr    *auth.JWTManager
}

// NewUserGroupHandler creates a new user group handler.
func NewUserGroupHandler(groupRepo *repository.UserGroupRepository, userRepo *repository.UserRepository, db *pgxpool.Pool, jwtMgr *auth.JWTManager) *UserGroupHandler {
	return &UserGroupHandler{groupRepo: groupRepo, userRepo: userRepo, db: db, jwtMgr: jwtMgr}
}

// Index handles GET /api/v1/groups.
func (h *UserGroupHandler) Index(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	if user == nil {
		return apperrors.ErrUnauthorized
	}

	groups, err := h.groupRepo.List(c.Context(), user.ID)
	if err != nil {
		return apperrors.ErrInternal
	}

	var data []fiber.Map
	for _, g := range groups {
		data = append(data, fiber.Map{
			"type": "user_groups",
			"id":   g.ID,
			"attributes": fiber.Map{
				"title": g.Title,
			},
		})
	}

	return c.JSON(fiber.Map{"data": data})
}

// Show handles GET /api/v1/groups/:id.
func (h *UserGroupHandler) Show(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	if user == nil {
		return apperrors.ErrUnauthorized
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.ErrBadRequest
	}

	// Check membership
	isMember, err := h.groupRepo.IsMember(c.Context(), user.ID, id)
	if err != nil || !isMember {
		return apperrors.ErrNotFound
	}

	group, err := h.groupRepo.FindByID(c.Context(), id)
	if err != nil {
		return apperrors.NotFoundResource("user_group", id)
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"type": "user_groups",
			"id":   group.ID,
			"attributes": fiber.Map{
				"title": group.Title,
			},
		},
	})
}

// Store handles POST /api/v1/groups.
func (h *UserGroupHandler) Store(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	if user == nil {
		return apperrors.ErrUnauthorized
	}

	var req struct {
		Title string `json:"title"`
	}
	if err := c.BodyParser(&req); err != nil || req.Title == "" {
		return apperrors.NewValidationError(map[string][]string{
			"title": {"A title is required."},
		})
	}

	group, err := h.groupRepo.Create(c.Context(), req.Title)
	if err != nil {
		return apperrors.ErrInternal
	}

	// Add the creator as owner member
	if err := h.groupRepo.AddMember(c.Context(), user.ID, group.ID); err != nil {
		return apperrors.ErrInternal
	}

	return c.Status(201).JSON(fiber.Map{
		"data": fiber.Map{
			"type": "user_groups",
			"id":   group.ID,
			"attributes": fiber.Map{
				"title": group.Title,
			},
		},
	})
}

// Update handles PUT /api/v1/groups/:id.
func (h *UserGroupHandler) Update(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	if user == nil {
		return apperrors.ErrUnauthorized
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.ErrBadRequest
	}

	// Check membership — only members can update
	isMember, err := h.groupRepo.IsMember(c.Context(), user.ID, id)
	if err != nil || !isMember {
		return apperrors.ErrNotFound
	}

	var req struct {
		Title string `json:"title"`
	}
	if err := c.BodyParser(&req); err != nil || req.Title == "" {
		return apperrors.NewValidationError(map[string][]string{
			"title": {"A title is required."},
		})
	}

	if err := h.groupRepo.Update(c.Context(), id, req.Title); err != nil {
		return apperrors.ErrInternal
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"type": "user_groups",
			"id":   id,
			"attributes": fiber.Map{
				"title": req.Title,
			},
		},
	})
}

// Delete handles DELETE /api/v1/groups/:id.
func (h *UserGroupHandler) Delete(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	if user == nil {
		return apperrors.ErrUnauthorized
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.ErrBadRequest
	}

	// Check membership in the TARGET group (not active group)
	isMember, err := h.groupRepo.IsMember(c.Context(), user.ID, id)
	if err != nil || !isMember {
		return apperrors.ErrNotFound
	}

	if err := h.groupRepo.Delete(c.Context(), id); err != nil {
		return apperrors.ErrInternal
	}

	// Invalidate tokens for all group members since the group no longer exists
	go func() {
		rows, err := h.db.Query(context.Background(), `SELECT user_id FROM group_memberships WHERE user_group_id = $1`, id)
		if err != nil {
			return
		}
		defer rows.Close()
		for rows.Next() {
			var uid uuid.UUID
			if rows.Scan(&uid) == nil {
				_ = h.userRepo.IncrementTokenVersion(context.Background(), uid)
			}
		}
	}()

	return c.Status(204).Send(nil)
}

// Switch handles POST /api/v1/groups/switch.
func (h *UserGroupHandler) Switch(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	if user == nil {
		return apperrors.ErrUnauthorized
	}

	var req struct {
		UserGroupID uuid.UUID `json:"user_group_id"`
	}
	if err := c.BodyParser(&req); err != nil || req.UserGroupID == uuid.Nil {
		return apperrors.NewValidationError(map[string][]string{
			"user_group_id": {"A valid group ID is required."},
		})
	}

	if err := h.userRepo.SetActiveGroup(c.Context(), user.ID, req.UserGroupID); err != nil {
		log.Printf("group switch failed: %v", err)
		return apperrors.New(400, "Failed to switch group. You may not be a member of the target group.")
	}

	// Re-issue JWT with updated group claim
	claims := auth.GetClaims(c)
	if claims != nil {
		identity := &auth.UserIdentity{
			ID:           claims.UserID,
			Email:        claims.Email,
			DemoUser:     claims.DemoUser,
			TokenVersion: claims.TokenVersion,
		}
		tokens, err := h.jwtMgr.GenerateTokenPair(identity, &req.UserGroupID)
		if err == nil {
			return c.JSON(fiber.Map{
				"data": fiber.Map{
					"type": "user_groups",
					"id":   req.UserGroupID,
				},
				"meta": fiber.Map{
					"message": "Active group switched successfully.",
				},
				"tokens": tokens,
			})
		}
	}

	// Fallback: return success without new tokens if JWT generation fails
	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"type": "user_groups",
			"id":   req.UserGroupID,
		},
		"meta": fiber.Map{
			"message": "Active group switched successfully.",
		},
	})
}

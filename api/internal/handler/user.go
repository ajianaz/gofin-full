package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/repository"
	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors"
)

// UserHandler handles user endpoints.
type UserHandler struct {
	repo *repository.UserRepository
}

// NewUserHandler creates a new user handler.
func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// ChangePassword handles POST /api/v1/users/me/password.
func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	if user == nil {
		return apperrors.ErrUnauthorized
	}

	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{
			"body": {"Invalid request body."},
		})
	}

	if req.CurrentPassword == "" || req.NewPassword == "" {
		return apperrors.NewValidationError(map[string][]string{
			"current_password": {"Current password is required."},
			"new_password":     {"New password is required."},
		})
	}

	if len(req.NewPassword) < 8 {
		return apperrors.NewValidationError(map[string][]string{
			"new_password": {"New password must be at least 8 characters."},
		})
	}

	// Fetch current user to verify old password
	u, err := h.repo.FindByID(c.Context(), user.ID)
	if err != nil {
		return apperrors.ErrInternal
	}

	if !auth.CheckPassword(u.Password, req.CurrentPassword) {
		return apperrors.NewWithDetail(401, "Unauthenticated", "Current password is incorrect.")
	}

	hash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		return apperrors.ErrInternal
	}

	if err := h.repo.Update(c.Context(), user.ID, "", hash); err != nil {
		return apperrors.ErrInternal
	}

	// Invalidate all existing tokens after password change
	_ = h.repo.IncrementTokenVersion(c.Context(), user.ID)

	return c.JSON(fiber.Map{"data": fiber.Map{"type": "users", "id": user.ID}})
}

// Show handles GET /api/v1/users/me.
func (h *UserHandler) Show(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	if user == nil {
		return apperrors.ErrUnauthorized
	}

	u, err := h.repo.FindByID(c.Context(), user.ID)
	if err != nil {
		return apperrors.NotFoundResource("user", user.ID)
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"type": "users",
			"id":   u.ID,
			"attributes": fiber.Map{
				"email":     u.Email,
				"blocked":   u.Blocked,
				"demo_user": u.DemoUser,
			},
		},
	})
}

// Update handles PUT /api/v1/users/me.
func (h *UserHandler) Update(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	if user == nil {
		return apperrors.ErrUnauthorized
	}

	var req struct {
		Email string `json:"email"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{
			"body": {"Invalid request body."},
		})
	}

	// Validate email if provided
	if req.Email != "" {
		if !isValidEmail(req.Email) {
			return apperrors.NewValidationError(map[string][]string{
				"email": {"Invalid email format."},
			})
		}

		// Check for duplicate email
		existing, err := h.repo.FindByEmail(c.Context(), req.Email)
		if err == nil && existing != nil && existing.ID != user.ID {
			return c.Status(409).JSON(fiber.Map{
				"message": "A user with this email already exists.",
			})
		}
	}

	if err := h.repo.Update(c.Context(), user.ID, req.Email, ""); err != nil {
		return apperrors.ErrInternal
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"type": "users",
			"id":   user.ID,
			"attributes": fiber.Map{
				"email": req.Email,
			},
		},
	})
}

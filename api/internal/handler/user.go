package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/dto/response"
	"github.com/ajianaz/gofin-full/api/internal/repository"
	"github.com/ajianaz/gofin-full/api/internal/validation"
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
		return validation.FieldErrors{"body": {"Invalid request body."}}.ToAppError()
	}

	// Validate input using validation framework
	errs := make(validation.FieldErrors)
	validation.Required("current_password", req.CurrentPassword, errs)
	validation.Required("new_password", req.NewPassword, errs)
	validation.PasswordStrength("new_password", req.NewPassword, errs)
	if errs.Has() {
		return errs.ToAppError()
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

	if err := h.repo.UpdatePassword(c.Context(), user.ID, hash); err != nil {
		return apperrors.ErrInternal
	}

	// Invalidate all existing tokens after password change
	_ = h.repo.IncrementTokenVersion(c.Context(), user.ID)
	auth.InvalidateTokenCache(c, user.ID)

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
				"name":      u.Name,
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
		Name  string `json:"name"`
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
			return response.SendError(c, 409, "A user with this email already exists.")
		}
	}

	if err := h.repo.Update(c.Context(), user.ID, req.Email, req.Name); err != nil {
		return apperrors.ErrInternal
	}

	// Return current user data (fetch to get actual values after update)
	u, err := h.repo.FindByID(c.Context(), user.ID)
	if err != nil {
		return apperrors.ErrInternal
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"type": "users",
			"id":   user.ID,
			"attributes": fiber.Map{
				"email": u.Email,
				"name":  u.Name,
			},
		},
	})
}

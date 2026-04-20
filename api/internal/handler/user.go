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
			"type":       "users",
			"id":         u.ID,
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

	if err := h.repo.Update(c.Context(), user.ID, req.Email, ""); err != nil {
		return apperrors.ErrInternal
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"type":       "users",
			"id":         user.ID,
			"attributes": fiber.Map{
				"email": req.Email,
			},
		},
	})
}

package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/repository"
	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors"
)

type PreferenceHandler struct {
	repo *repository.PreferenceRepository
}

func NewPreferenceHandler(repo *repository.PreferenceRepository) *PreferenceHandler {
	return &PreferenceHandler{repo: repo}
}

func (h *PreferenceHandler) Index(c *fiber.Ctx) error {
	user := auth.GetUser(c)

	prefs, err := h.repo.List(c.Context(), user.ID)
	if err != nil {
		return apperrors.NewWithDetail(500, "failed to list preferences", err.Error())
	}

	var data []fiber.Map
	for _, p := range prefs {
		data = append(data, fiber.Map{
			"type": "preferences",
			"id":   p.ID,
			"attributes": fiber.Map{
				"name": p.Name,
				"data": p.Data,
			},
		})
	}
	return c.JSON(fiber.Map{"data": data})
}

func (h *PreferenceHandler) Show(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	name := c.Params("name")
	if name == "" {
		return apperrors.NewValidationError(map[string][]string{"name": {"name is required"}})
	}

	p, err := h.repo.Get(c.Context(), user.ID, name)
	if err != nil {
		return apperrors.NotFoundResource("preference", 0)
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "preferences",
		"id":   p.ID,
		"attributes": fiber.Map{
			"name": p.Name,
			"data": p.Data,
		},
	}})
}

func (h *PreferenceHandler) Set(c *fiber.Ctx) error {
	user := auth.GetUser(c)

	var req struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}
	if req.Name == "" {
		return apperrors.NewValidationError(map[string][]string{"name": {"name is required"}})
	}

	p, err := h.repo.Set(c.Context(), user.ID, req.Name, req.Data)
	if err != nil {
		return apperrors.NewWithDetail(500, "failed to set preference", err.Error())
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "preferences",
		"id":   p.ID,
		"attributes": fiber.Map{
			"name": p.Name,
			"data": p.Data,
		},
	}})
}

func (h *PreferenceHandler) Delete(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	name := c.Params("name")
	if name == "" {
		return apperrors.NewValidationError(map[string][]string{"name": {"name is required"}})
	}

	if err := h.repo.Delete(c.Context(), user.ID, name); err != nil {
		return apperrors.NotFoundResource("preference", 0)
	}

	return c.Status(204).Send(nil)
}

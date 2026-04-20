package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/azfirazka/gofin-full/api/internal/auth"
	"github.com/azfirazka/gofin-full/api/internal/repository"
	apperrors "github.com/azfirazka/gofin-full/api/pkg/errors"
)

type ConfigurationHandler struct {
	repo *repository.ConfigurationRepository
}

func NewConfigurationHandler(repo *repository.ConfigurationRepository) *ConfigurationHandler {
	return &ConfigurationHandler{repo: repo}
}

func (h *ConfigurationHandler) Index(c *fiber.Ctx) error {
	_ = auth.GetUser(c)

	configs, err := h.repo.List(c.Context())
	if err != nil {
		return apperrors.NewWithDetail(500, "failed to list configurations", err.Error())
	}

	var data []fiber.Map
	for _, cfg := range configs {
		data = append(data, fiber.Map{
			"type": "configurations",
			"id":   cfg.ID,
			"attributes": fiber.Map{
				"name":  cfg.Name,
				"value": cfg.Value,
			},
		})
	}
	return c.JSON(fiber.Map{"data": data})
}

func (h *ConfigurationHandler) Show(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	name := c.Params("name")
	if name == "" {
		return apperrors.NewValidationError(map[string][]string{"name": {"name is required"}})
	}

	cfg, err := h.repo.Get(c.Context(), name)
	if err != nil {
		return apperrors.NotFoundResource("configuration", 0)
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "configurations",
		"id":   cfg.ID,
		"attributes": fiber.Map{
			"name":  cfg.Name,
			"value": cfg.Value,
		},
	}})
}

func (h *ConfigurationHandler) Set(c *fiber.Ctx) error {
	_ = auth.GetUser(c)

	var req struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}
	if req.Name == "" {
		return apperrors.NewValidationError(map[string][]string{"name": {"name is required"}})
	}

	cfg, err := h.repo.Set(c.Context(), req.Name, req.Value)
	if err != nil {
		return apperrors.NewWithDetail(500, "failed to set configuration", err.Error())
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "configurations",
		"id":   cfg.ID,
		"attributes": fiber.Map{
			"name":  cfg.Name,
			"value": cfg.Value,
		},
	}})
}

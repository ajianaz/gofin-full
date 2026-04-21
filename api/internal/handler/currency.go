package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/repository"
	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors"
)

type CurrencyHandler struct {
	repo *repository.CurrencyRepository
}

func NewCurrencyHandler(repo *repository.CurrencyRepository) *CurrencyHandler {
	return &CurrencyHandler{repo: repo}
}

func (h *CurrencyHandler) Index(c *fiber.Ctx) error {
	_ = auth.GetUser(c)

	currencies, err := h.repo.List(c.Context())
	if err != nil {
		return apperrors.NewWithDetail(500, "failed to list currencies", err.Error())
	}

	var data []fiber.Map
	for _, cur := range currencies {
		data = append(data, fiber.Map{
			"type": "currencies",
			"id":   cur.Code,
			"attributes": fiber.Map{
				"name":           cur.Name,
				"symbol":         cur.Symbol,
				"decimal_places": cur.DecimalPlaces,
				"enabled":        cur.Enabled,
			},
		})
	}
	return c.JSON(fiber.Map{"data": data})
}

func (h *CurrencyHandler) Show(c *fiber.Ctx) error {
	code := c.Params("code")
	if code == "" {
		return apperrors.NewValidationError(map[string][]string{"code": {"code is required"}})
	}

	cur, err := h.repo.FindByCode(c.Context(), code)
	if err != nil {
		return apperrors.NotFoundResource("currency", uuid.Nil)
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "currencies",
		"id":   cur.Code,
		"attributes": fiber.Map{
			"name":           cur.Name,
			"symbol":         cur.Symbol,
			"decimal_places": cur.DecimalPlaces,
			"enabled":        cur.Enabled,
		},
	}})
}

package handler

import (
"time"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/repository"
	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors")

type ExchangeRateHandler struct {
	repo *repository.ExchangeRateRepository
}

func NewExchangeRateHandler(repo *repository.ExchangeRateRepository) *ExchangeRateHandler {
	return &ExchangeRateHandler{repo: repo}
}

func (h *ExchangeRateHandler) Index(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	rates, err := h.repo.List(c.Context(), *groupID)
	if err != nil {
		log.Printf("handler/Index: failed to list exchange rates: %v", err)
		return apperrors.ErrInternal
	}

	var data []fiber.Map
	for _, r := range rates {
		data = append(data, fiber.Map{
			"type": "exchange_rates",
			"id":   r.ID,
			"attributes": fiber.Map{
				"from_currency_id": r.FromCurrencyID,
				"to_currency_id":   r.ToCurrencyID,
				"rate":             r.Rate.StringFixed(6),
				"date":             r.Date.Format(time.RFC3339),
			},
		})
	}
	return c.JSON(fiber.Map{"data": data})
}

func (h *ExchangeRateHandler) Store(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	var req struct {
		FromCurrencyID string `json:"from_currency_id"`
		ToCurrencyID   string `json:"to_currency_id"`
		Rate           string `json:"rate"`
		Date           string `json:"date"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}
	if req.FromCurrencyID == "" || req.ToCurrencyID == "" {
		return apperrors.NewValidationError(map[string][]string{"from_currency_id": {"from and to currency are required"}})
	}

	rate, err := decimal.NewFromString(req.Rate)
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"rate": {"invalid rate value"}})
	}

	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		date = time.Now().UTC()
	}

	er, err := h.repo.Create(c.Context(), user.ID, *groupID, req.FromCurrencyID, req.ToCurrencyID, rate, date)
	if err != nil {
		log.Printf("handler/Index: failed to create exchange rate: %v", err)
		return apperrors.ErrInternal
	}

	return c.Status(201).JSON(fiber.Map{"data": fiber.Map{
		"type": "exchange_rates",
		"id":   er.ID,
		"attributes": fiber.Map{
			"from_currency_id": er.FromCurrencyID,
			"to_currency_id":   er.ToCurrencyID,
			"rate":             er.Rate.StringFixed(6),
			"date":             er.Date.Format(time.RFC3339),
		},
	}})
}

func (h *ExchangeRateHandler) Show(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	from := c.Query("from")
	to := c.Query("to")
	dateStr := c.Query("date")
	if from == "" || to == "" || dateStr == "" {
		return apperrors.NewValidationError(map[string][]string{
			"query": {"from, to, and date query params are required"},
		})
	}

	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"date": {"invalid date format, use RFC3339"}})
	}

	rate, err := h.repo.FindRate(c.Context(), *groupID, from, to, date)
	if err != nil {
		log.Printf("handler/Index: failed to find exchange rate: %v", err)
		return apperrors.ErrInternal
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "exchange_rates",
		"attributes": fiber.Map{
			"from_currency_id": from,
			"to_currency_id":   to,
			"rate":             rate.StringFixed(6),
			"date":             dateStr,
		},
	}})
}

func (h *ExchangeRateHandler) Delete(c *fiber.Ctx) error {
	_ = auth.GetUser(c)

	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(403, "active group not found")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id format"}})
	}

	if err := h.repo.Delete(c.Context(), *groupID, id); err != nil {
		return apperrors.NotFoundResource("exchange_rate", id)
	}

	return c.Status(204).Send(nil)
}

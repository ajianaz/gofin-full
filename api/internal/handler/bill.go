package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"

	"github.com/azfirazka/gofin-full/api/internal/auth"
	"github.com/azfirazka/gofin-full/api/internal/repository"
	apperrors "github.com/azfirazka/gofin-full/api/pkg/errors"
)

type BillHandler struct {
	repo *repository.BillRepository
}

func NewBillHandler(repo *repository.BillRepository) *BillHandler {
	return &BillHandler{repo: repo}
}

func (h *BillHandler) Index(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	bills, err := h.repo.List(c.Context(), *groupID)
	if err != nil {
		return apperrors.NewWithDetail(500, "failed to list bills", err.Error())
	}

	var data []fiber.Map
	for _, b := range bills {
		data = append(data, fiber.Map{
			"type": "bills",
			"id":   b.ID,
			"attributes": fiber.Map{
				"name":        b.Name,
				"amount_min":  b.AmountMin.StringFixed(2),
				"amount_max":  b.AmountMax.StringFixed(2),
				"date":        b.Date.Format(time.RFC3339),
				"end_date":    fmtTime(b.EndDate),
				"repeat_freq": b.RepeatFreq,
				"active":      b.Active,
				"currency_id": b.CurrencyID,
			},
		})
	}
	return c.JSON(fiber.Map{"data": data})
}

func (h *BillHandler) Show(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id"}})
	}

	b, err := h.repo.FindByID(c.Context(), int64(id), *groupID)
	if err != nil {
		return apperrors.NotFoundResource("bill", int64(id))
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "bills",
		"id":   b.ID,
		"attributes": fiber.Map{
			"name":        b.Name,
			"amount_min":  b.AmountMin.StringFixed(2),
			"amount_max":  b.AmountMax.StringFixed(2),
			"date":        b.Date.Format(time.RFC3339),
			"end_date":    fmtTime(b.EndDate),
			"repeat_freq": b.RepeatFreq,
			"skip":        b.Skip,
			"active":      b.Active,
			"order":       b.Order,
			"notes":       b.Notes,
			"currency_id": b.CurrencyID,
		},
	}})
}

func (h *BillHandler) Store(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	var req struct {
		Name       string  `json:"name"`
		AmountMin  string  `json:"amount_min"`
		AmountMax  string  `json:"amount_max"`
		Date       string  `json:"date"`
		RepeatFreq string  `json:"repeat_freq"`
		CurrencyID string  `json:"currency_id"`
		Order      int     `json:"order"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}
	if req.Name == "" {
		return apperrors.NewValidationError(map[string][]string{"name": {"name is required"}})
	}

	amountMin, err := decimal.NewFromString(req.AmountMin)
	if err != nil {
		amountMin = decimal.Zero
	}
	amountMax, err := decimal.NewFromString(req.AmountMax)
	if err != nil {
		amountMax = decimal.Zero
	}

	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		date = time.Now().UTC()
	}

	if req.RepeatFreq == "" {
		req.RepeatFreq = "monthly"
	}

	b, err := h.repo.Create(c.Context(), user.ID, *groupID, req.Name, amountMin, amountMax, date, req.RepeatFreq, req.CurrencyID, req.Order)
	if err != nil {
		return apperrors.NewWithDetail(500, "failed to create bill", err.Error())
	}

	return c.Status(201).JSON(fiber.Map{"data": fiber.Map{
		"type": "bills",
		"id":   b.ID,
		"attributes": fiber.Map{
			"name":        b.Name,
			"amount_min":  b.AmountMin.StringFixed(2),
			"amount_max":  b.AmountMax.StringFixed(2),
			"date":        b.Date.Format(time.RFC3339),
			"repeat_freq": b.RepeatFreq,
			"currency_id": b.CurrencyID,
		},
	}})
}

func (h *BillHandler) Update(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id"}})
	}

	var req struct {
		Name   string  `json:"name"`
		Active *bool   `json:"active"`
		Notes  *string `json:"notes"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}

	if err := h.repo.Update(c.Context(), int64(id), *groupID, req.Name, req.Active, req.Notes); err != nil {
		return apperrors.NotFoundResource("bill", int64(id))
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "bills", "id": id,
		"attributes": fiber.Map{"name": req.Name},
	}})
}

func (h *BillHandler) Delete(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id"}})
	}

	if err := h.repo.Delete(c.Context(), int64(id), *groupID); err != nil {
		return apperrors.NotFoundResource("bill", int64(id))
	}

	return c.Status(204).Send(nil)
}

// fmtTime is a helper for optional time pointers.
func fmtTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

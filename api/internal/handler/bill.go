package handler

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/domain"
	"github.com/ajianaz/gofin-full/api/internal/repository"
	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors"
)

type BillHandler struct {
	repo  *repository.BillRepository
	curry *CurrencyResolver
}

func NewBillHandler(repo *repository.BillRepository, curry *CurrencyResolver) *BillHandler {
	return &BillHandler{repo: repo, curry: curry}
}

func (h *BillHandler) Index(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	bills, err := h.repo.List(c.Context(), *groupID)
	if err != nil {
		log.Printf("handler/Index: failed to list bills: %v", err)
		return apperrors.ErrInternal
	}

	// Batch resolve currency
	cIDs := make([]string, 0, len(bills))
	for _, b := range bills {
		if b.CurrencyID != "" {
			cIDs = append(cIDs, b.CurrencyID)
		}
	}
	cMap := h.curry.ResolveMany(c.Context(), cIDs)

	var data []fiber.Map
	for _, b := range bills {
		data = append(data, billToMap(b, cMap))
	}
	return c.JSON(fiber.Map{"data": data})
}

func (h *BillHandler) Show(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id format"}})
	}

	b, err := h.repo.FindByID(c.Context(), id, *groupID)
	if err != nil {
		return apperrors.NotFoundResource("bill", id)
	}

	cMap := h.curry.ResolveMany(c.Context(), []string{b.CurrencyID})

	return c.JSON(fiber.Map{"data": billToMapFull(b, cMap)})
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
		log.Printf("handler/Index: failed to create bill: %v", err)
		return apperrors.ErrInternal
	}

	cMap := h.curry.ResolveMany(c.Context(), []string{b.CurrencyID})

	return c.Status(201).JSON(fiber.Map{"data": billToMap(*b, cMap)})
}

func (h *BillHandler) Update(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id format"}})
	}

	var req struct {
		Name   string  `json:"name"`
		Active *bool   `json:"active"`
		Notes  *string `json:"notes"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}

	if err := h.repo.Update(c.Context(), id, *groupID, req.Name, req.Active, req.Notes); err != nil {
		return apperrors.NotFoundResource("bill", id)
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

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id format"}})
	}

	if err := h.repo.Delete(c.Context(), id, *groupID); err != nil {
		return apperrors.NotFoundResource("bill", id)
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

func billToMap(b domain.Bill, cMap map[string]CurrencyInfo) fiber.Map {
	attrs := fiber.Map{
		"name":        b.Name,
		"amount_min":  b.AmountMin.StringFixed(2),
		"amount_max":  b.AmountMax.StringFixed(2),
		"date":        b.Date.Format(time.RFC3339),
		"end_date":    fmtTime(b.EndDate),
		"repeat_freq": b.RepeatFreq,
		"active":      b.Active,
		"currency_id": b.CurrencyID,
	}
	if b.CurrencyID != "" {
		if ci, ok := cMap[b.CurrencyID]; ok {
			attrs["currency_code"] = ci.Code
			attrs["currency_symbol"] = ci.Symbol
			attrs["currency_decimal_places"] = ci.DecimalPlaces
		}
	}
	return fiber.Map{"type": "bills", "id": b.ID, "attributes": attrs}
}

func billToMapFull(b *domain.Bill, cMap map[string]CurrencyInfo) fiber.Map {
	attrs := fiber.Map{
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
	}
	if b.CurrencyID != "" {
		if ci, ok := cMap[b.CurrencyID]; ok {
			attrs["currency_code"] = ci.Code
			attrs["currency_symbol"] = ci.Symbol
			attrs["currency_decimal_places"] = ci.DecimalPlaces
		}
	}
	return fiber.Map{"type": "bills", "id": b.ID, "attributes": attrs}
}

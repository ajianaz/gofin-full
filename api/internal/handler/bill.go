package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/domain"
	"github.com/ajianaz/gofin-full/api/internal/repository"
	"github.com/ajianaz/gofin-full/api/internal/validation"
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

	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 20)
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	} else if perPage > 100 {
		perPage = 100
	}

	bills, total, err := h.repo.ListPaginated(c.Context(), *groupID, page, perPage)
	if err != nil {
		log.Error().Err(err).Msg("handler/Index: failed to list bills")
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

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	return c.JSON(fiber.Map{
		"data": data,
		"meta": fiber.Map{
			"pagination": fiber.Map{
				"total":        total,
				"count":        len(data),
				"per_page":     perPage,
				"current_page": page,
				"total_pages":  totalPages,
			},
		},
	})
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

	errs := make(validation.FieldErrors)
	validation.NonNegative("amount_min", req.AmountMin, errs)
	validation.NonNegative("amount_max", req.AmountMax, errs)
	validation.MinGreaterThanMax("amount_min", "amount_max", req.AmountMin, req.AmountMax, errs)
	if req.RepeatFreq != "" {
		validation.OneOf("repeat_freq", req.RepeatFreq, []string{"weekly", "monthly", "quarterly", "half-yearly", "yearly"}, errs)
	}
	if req.Date != "" {
		validation.OptionalDateString("date", req.Date, errs)
	}
	if errs.Has() {
		return errs.ToAppError()
	}

	amountMin, _ := decimal.NewFromString(req.AmountMin)
	if amountMin.IsZero() && req.AmountMin == "" {
		amountMin = decimal.Zero
	}
	amountMax, _ := decimal.NewFromString(req.AmountMax)
	if amountMax.IsZero() && req.AmountMax == "" {
		amountMax = decimal.Zero
	}
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		dateErrs := make(validation.FieldErrors)
		validation.OptionalDateString("date", req.Date, dateErrs)
		return dateErrs.ToAppError()
	}

	b, err := h.repo.Create(c.Context(), user.ID, *groupID, sanitizeStr(req.Name), amountMin, amountMax, date, req.RepeatFreq, req.CurrencyID, req.Order)
	if err != nil {
		log.Error().Err(err).Msg("handler/Index: failed to create bill")
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

	if err := h.repo.Update(c.Context(), id, *groupID, sanitizeStr(req.Name), req.Active, sanitizePtr(req.Notes)); err != nil {
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
	dp := 2
	if b.CurrencyID != "" {
		if ci, ok := cMap[b.CurrencyID]; ok {
			dp = ci.DecimalPlaces
		}
	}
	attrs := fiber.Map{
		"name":        b.Name,
		"amount_min":  b.AmountMin.StringFixed(int32(dp)),
		"amount_max":  b.AmountMax.StringFixed(int32(dp)),
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
	dp := 2
	if b.CurrencyID != "" {
		if ci, ok := cMap[b.CurrencyID]; ok {
			dp = ci.DecimalPlaces
		}
	}
	attrs := fiber.Map{
		"name":        b.Name,
		"amount_min":  b.AmountMin.StringFixed(int32(dp)),
		"amount_max":  b.AmountMax.StringFixed(int32(dp)),
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

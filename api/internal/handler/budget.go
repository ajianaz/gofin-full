package handler

import (
	"github.com/rs/zerolog/log"
"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/repository"
	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors")

type BudgetHandler struct {
	repo  *repository.BudgetRepository
	curry *CurrencyResolver
}

func NewBudgetHandler(repo *repository.BudgetRepository, curry *CurrencyResolver) *BudgetHandler {
	return &BudgetHandler{repo: repo, curry: curry}
}

func (h *BudgetHandler) Index(c *fiber.Ctx) error {
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

	budgets, total, err := h.repo.ListPaginated(c.Context(), *groupID, page, perPage)
	if err != nil {
		log.Error().Err(err).Msg("handler/Index: failed to list budgets")
		return apperrors.ErrInternal
	}

	var data []fiber.Map
	for _, b := range budgets {
		data = append(data, fiber.Map{
			"type":       "budgets",
			"id":         b.ID,
			"attributes": fiber.Map{"name": b.Name, "active": b.Active, "order": b.Order},
		})
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

func (h *BudgetHandler) Show(c *fiber.Ctx) error {
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
		return apperrors.NotFoundResource("budget", id)
	}

	var limits []fiber.Map
	for _, l := range b.Limits {
		limits = append(limits, fiber.Map{
			"id": l.ID, "start": l.Start, "end": l.End,
			"amount": l.Amount.StringFixed(2),
		})
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "budgets", "id": b.ID,
		"attributes": fiber.Map{"name": b.Name, "active": b.Active, "order": b.Order, "limits": limits},
	}})
}

func (h *BudgetHandler) Store(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	var req struct {
		Name  string `json:"name"`
		Order int    `json:"order"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}
	if req.Name == "" {
		return apperrors.NewValidationError(map[string][]string{"name": {"name is required"}})
	}

	b, err := h.repo.Create(c.Context(), user.ID, *groupID, sanitizeStr(req.Name), req.Order)
	if err != nil {
		log.Error().Err(err).Msg("handler/Index: failed to create budget")
		return apperrors.ErrInternal
	}

	return c.Status(201).JSON(fiber.Map{"data": fiber.Map{
		"type": "budgets", "id": b.ID,
		"attributes": fiber.Map{"name": b.Name, "active": b.Active, "order": b.Order},
	}})
}

func (h *BudgetHandler) Update(c *fiber.Ctx) error {
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
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}

	// Verify budget exists before update
	if _, err := h.repo.FindByID(c.Context(), id, *groupID); err != nil {
		return apperrors.NotFoundResource("budget", id)
	}
	if err := h.repo.Update(c.Context(), id, *groupID, sanitizeStr(req.Name), req.Active); err != nil {
		return apperrors.NotFoundResource("budget", id)
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "budgets", "id": id,
		"attributes": fiber.Map{"name": req.Name},
	}})
}

func (h *BudgetHandler) Delete(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id format"}})
	}

	// Verify budget exists before delete
	if _, err := h.repo.FindByID(c.Context(), id, *groupID); err != nil {
		return apperrors.NotFoundResource("budget", id)
	}
	if err := h.repo.Delete(c.Context(), id, *groupID); err != nil {
		return apperrors.NotFoundResource("budget", id)
	}

	return c.Status(204).Send(nil)
}

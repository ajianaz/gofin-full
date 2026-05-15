package handler

import (
"time"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/domain"
	"github.com/ajianaz/gofin-full/api/internal/repository"
	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors")

type RecurrenceHandler struct {
	repo *repository.RecurrenceRepository
}

func NewRecurrenceHandler(repo *repository.RecurrenceRepository) *RecurrenceHandler {
	return &RecurrenceHandler{repo: repo}
}

func (h *RecurrenceHandler) Index(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	recurrences, err := h.repo.List(c.Context(), *groupID)
	if err != nil {
		log.Printf("handler/Index: failed to list recurrences: %v", err)
		return apperrors.ErrInternal
	}

	var data []fiber.Map
	for _, rec := range recurrences {
		data = append(data, fiber.Map{
			"type": "recurrences", "id": rec.ID,
			"attributes": fiber.Map{
				"title": rec.Title, "first_date": rec.FirstDate,
				"repeat_freq": rec.RepeatFreq, "active": rec.Active,
			},
		})
	}
	return c.JSON(fiber.Map{"data": data})
}

func (h *RecurrenceHandler) Show(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id format"}})
	}

	rec, err := h.repo.FindByID(c.Context(), id, *groupID)
	if err != nil {
		return apperrors.NotFoundResource("recurrence", id)
	}

	return c.JSON(fiber.Map{"data": recurrenceToMap(rec)})
}

func (h *RecurrenceHandler) Store(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	var req struct {
		Title        string                      `json:"title"`
		Description  *string                     `json:"description"`
		FirstDate    time.Time                   `json:"first_date"`
		RepeatFreq   string                      `json:"repeat_freq"`
		RepeatUntil  *time.Time                  `json:"repeat_until"`
		Transactions []domain.RecurringTransaction `json:"transactions"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}

	fieldErrors := make(map[string][]string)
	if req.Title == "" {
		fieldErrors["title"] = append(fieldErrors["title"], "title is required")
	}
	if req.FirstDate.IsZero() {
		fieldErrors["first_date"] = append(fieldErrors["first_date"], "first_date is required")
	}
	if len(fieldErrors) > 0 {
		return apperrors.NewValidationError(fieldErrors)
	}

	if req.RepeatFreq == "" {
		req.RepeatFreq = "monthly"
	}

	rec, err := h.repo.Create(c.Context(), user.ID, *groupID, req.Title, req.FirstDate, req.RepeatFreq)
	if err != nil {
		log.Printf("handler/Index: failed to create recurrence: %v", err)
		return apperrors.ErrInternal
	}

	// Add transaction templates
	for _, txn := range req.Transactions {
		h.repo.AddTransaction(c.Context(), rec.ID, &txn)
	}

	return c.Status(201).JSON(fiber.Map{"data": fiber.Map{
		"type": "recurrences", "id": rec.ID,
		"attributes": fiber.Map{"title": rec.Title, "repeat_freq": rec.RepeatFreq},
	}})
}

func (h *RecurrenceHandler) Update(c *fiber.Ctx) error {
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
		Title       string     `json:"title"`
		RepeatFreq  string     `json:"repeat_freq"`
		Active      *bool      `json:"active"`
		Description *string    `json:"description"`
		RepeatUntil *time.Time `json:"repeat_until"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}

	if err := h.repo.Update(c.Context(), id, *groupID, req.Title, req.RepeatFreq, req.Active, req.Description, req.RepeatUntil); err != nil {
		return apperrors.NotFoundResource("recurrence", id)
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "recurrences", "id": id,
		"attributes": fiber.Map{"title": req.Title},
	}})
}

func (h *RecurrenceHandler) Delete(c *fiber.Ctx) error {
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
		return apperrors.NotFoundResource("recurrence", id)
	}

	return c.Status(204).Send(nil)
}

func recurrenceToMap(r *domain.Recurrence) fiber.Map {
	var txns []fiber.Map
	for _, t := range r.Transactions {
		txn := fiber.Map{
			"id": t.ID, "type": t.Type, "description": t.Description,
			"amount": t.Amount.StringFixed(2),
		}
		if t.SourceID != uuid.Nil {
			txn["source_id"] = t.SourceID
		}
		if t.DestinationID != uuid.Nil {
			txn["destination_id"] = t.DestinationID
		}
		txns = append(txns, txn)
	}

	return fiber.Map{
		"type": "recurrences", "id": r.ID,
		"attributes": fiber.Map{
			"title": r.Title, "description": r.Description,
			"first_date": r.FirstDate, "repeat_freq": r.RepeatFreq,
			"repeat_until": r.RepeatUntil, "active": r.Active,
			"transactions": txns,
		},
	}
}

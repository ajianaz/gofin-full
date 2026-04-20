package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/azfirazka/gofin-full/api/internal/auth"
	"github.com/azfirazka/gofin-full/api/internal/repository"
	apperrors "github.com/azfirazka/gofin-full/api/pkg/errors"
)

type NoteHandler struct {
	repo *repository.NoteRepository
}

func NewNoteHandler(repo *repository.NoteRepository) *NoteHandler {
	return &NoteHandler{repo: repo}
}

func (h *NoteHandler) Index(c *fiber.Ctx) error {
	_ = auth.GetUser(c)

	noteableType := c.Query("noteable_type")
	noteableID := c.QueryInt("noteable_id")
	if noteableType == "" || noteableID == 0 {
		return apperrors.NewValidationError(map[string][]string{
			"query": {"noteable_type and noteable_id query params are required"},
		})
	}

	notes, err := h.repo.ListByEntity(c.Context(), noteableType, int64(noteableID))
	if err != nil {
		return apperrors.NewWithDetail(500, "failed to list notes", err.Error())
	}

	var data []fiber.Map
	for _, n := range notes {
		data = append(data, fiber.Map{
			"type": "notes",
			"id":   n.ID,
			"attributes": fiber.Map{
				"note":           n.Note,
				"noteable_type":  n.NoteableType,
				"noteable_id":    n.NoteableID,
			},
		})
	}
	return c.JSON(fiber.Map{"data": data})
}

func (h *NoteHandler) Store(c *fiber.Ctx) error {
	_ = auth.GetUser(c)

	var req struct {
		NoteableType string `json:"noteable_type"`
		NoteableID   int64  `json:"noteable_id"`
		Note         string `json:"note"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}
	if req.NoteableType == "" {
		return apperrors.NewValidationError(map[string][]string{"noteable_type": {"noteable_type is required"}})
	}

	n, err := h.repo.Create(c.Context(), req.NoteableType, req.NoteableID, req.Note)
	if err != nil {
		return apperrors.NewWithDetail(500, "failed to create note", err.Error())
	}

	return c.Status(201).JSON(fiber.Map{"data": fiber.Map{
		"type": "notes",
		"id":   n.ID,
		"attributes": fiber.Map{
			"note":           n.Note,
			"noteable_type":  n.NoteableType,
			"noteable_id":    n.NoteableID,
		},
	}})
}

func (h *NoteHandler) Update(c *fiber.Ctx) error {
	_ = auth.GetUser(c)

	id, err := c.ParamsInt("id")
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id"}})
	}

	var req struct {
		Note string `json:"note"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}

	if err := h.repo.Update(c.Context(), int64(id), req.Note); err != nil {
		return apperrors.NotFoundResource("note", int64(id))
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type":       "notes", "id": id,
		"attributes": fiber.Map{"note": req.Note},
	}})
}

func (h *NoteHandler) Delete(c *fiber.Ctx) error {
	_ = auth.GetUser(c)

	id, err := c.ParamsInt("id")
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id"}})
	}

	if err := h.repo.Delete(c.Context(), int64(id)); err != nil {
		return apperrors.NotFoundResource("note", int64(id))
	}

	return c.Status(204).Send(nil)
}

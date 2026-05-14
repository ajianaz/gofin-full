package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/domain"
	"github.com/ajianaz/gofin-full/api/internal/repository"
	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors"
)

type NoteHandler struct {
	repo *repository.NoteRepository
}

func NewNoteHandler(repo *repository.NoteRepository) *NoteHandler {
	return &NoteHandler{repo: repo}
}

func (h *NoteHandler) Index(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	if user == nil {
		return apperrors.ErrUnauthorized
	}
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	noteableType := c.Query("noteable_type")
	noteableIDStr := c.Query("noteable_id")
	if noteableType == "" || noteableIDStr == "" {
		return apperrors.NewValidationError(map[string][]string{
			"query": {"noteable_type and noteable_id query params are required"},
		})
	}

	noteableID, err := uuid.Parse(noteableIDStr)
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{
			"noteable_id": {"invalid noteable_id format"},
		})
	}

	notes, err := h.repo.ListByEntity(c.Context(), noteableType, noteableID)
	if err != nil {
		log.Printf("handler: failed to list notes: %v", err)
		return apperrors.ErrInternal
	}

	// Filter notes to only include those owned by the current user
	var filtered []domain.Note
	for _, n := range notes {
		if n.UserID == user.ID {
			filtered = append(filtered, n)
		}
	}

	var data []fiber.Map
	for _, n := range filtered {
		data = append(data, fiber.Map{
			"type": "notes",
			"id":   n.ID,
			"attributes": fiber.Map{
				"note":          n.Note,
				"noteable_type": n.NoteableType,
				"noteable_id":   n.NoteableID,
			},
		})
	}
	return c.JSON(fiber.Map{"data": data})
}

func (h *NoteHandler) Store(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	if user == nil {
		return apperrors.ErrUnauthorized
	}
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	var req struct {
		NoteableType string    `json:"noteable_type"`
		NoteableID   uuid.UUID `json:"noteable_id"`
		Note         string    `json:"note"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}
	if req.NoteableType == "" {
		return apperrors.NewValidationError(map[string][]string{"noteable_type": {"noteable_type is required"}})
	}

	n, err := h.repo.Create(c.Context(), user.ID, *groupID, req.NoteableType, req.NoteableID, req.Note)
	if err != nil {
		log.Printf("handler: failed to create note: %v", err)
		return apperrors.ErrInternal
	}

	return c.Status(201).JSON(fiber.Map{"data": fiber.Map{
		"type": "notes",
		"id":   n.ID,
		"attributes": fiber.Map{
			"note":          n.Note,
			"noteable_type": n.NoteableType,
			"noteable_id":   n.NoteableID,
		},
	}})
}

func (h *NoteHandler) Update(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	if user == nil {
		return apperrors.ErrUnauthorized
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id format"}})
	}

	n, err := h.repo.FindByID(c.Context(), id, user.ID)
	if err != nil || n == nil {
		return apperrors.NotFoundResource("note", id)
	}

	var req struct {
		Note string `json:"note"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}

	if err := h.repo.Update(c.Context(), id, user.ID, req.Note); err != nil {
		return apperrors.NotFoundResource("note", id)
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "notes", "id": id,
		"attributes": fiber.Map{"note": req.Note},
	}})
}

func (h *NoteHandler) Delete(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	if user == nil {
		return apperrors.ErrUnauthorized
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id format"}})
	}

	n, err := h.repo.FindByID(c.Context(), id, user.ID)
	if err != nil || n == nil {
		return apperrors.NotFoundResource("note", id)
	}

	if err := h.repo.Delete(c.Context(), id, user.ID); err != nil {
		return apperrors.NotFoundResource("note", id)
	}

	return c.Status(204).Send(nil)
}

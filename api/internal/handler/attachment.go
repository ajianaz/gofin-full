package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/repository"
	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors"
)

type AttachmentHandler struct {
	repo *repository.AttachmentRepository
}

func NewAttachmentHandler(repo *repository.AttachmentRepository) *AttachmentHandler {
	return &AttachmentHandler{repo: repo}
}

func (h *AttachmentHandler) Index(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	if user == nil {
		return apperrors.ErrUnauthorized
	}

	attachableType := c.Query("attachable_type")
	attachableIDStr := c.Query("attachable_id")
	if attachableType == "" || attachableIDStr == "" {
		return apperrors.NewValidationError(map[string][]string{
			"query": {"attachable_type and attachable_id query params are required"},
		})
	}

	attachableID, err := uuid.Parse(attachableIDStr)
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{
			"attachable_id": {"invalid attachable_id format"},
		})
	}

	attachments, err := h.repo.ListByEntityAndUser(c.Context(), attachableType, attachableID, user.ID)
	if err != nil {
		log.Printf("handler: failed to list attachments: %v", err)
		return apperrors.ErrInternal
	}

	var data []fiber.Map
	for _, a := range attachments {
		data = append(data, fiber.Map{
			"type": "attachments",
			"id":   a.ID,
			"attributes": fiber.Map{
				"filename":  a.Filename,
				"mime_type": a.MimeType,
				"size":      a.Size,
				"uploaded":  a.Uploaded,
			},
		})
	}
	return c.JSON(fiber.Map{"data": data})
}

func (h *AttachmentHandler) Show(c *fiber.Ctx) error {
	user := auth.GetUser(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id format"}})
	}

	a, err := h.repo.FindByID(c.Context(), id)
	if err != nil {
		return apperrors.NotFoundResource("attachment", id)
	}

	if user == nil || a.UserID != user.ID {
		return apperrors.ErrNotFound
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "attachments",
		"id":   a.ID,
		"attributes": fiber.Map{
			"filename":        a.Filename,
			"mime_type":       a.MimeType,
			"size":            a.Size,
			"uploaded":        a.Uploaded,
			"attachable_type": a.AttachableType,
			"attachable_id":   a.AttachableID,
		},
	}})
}

func (h *AttachmentHandler) Store(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	if user == nil {
		return apperrors.ErrUnauthorized
	}

	var req struct {
		AttachableType string    `json:"attachable_type"`
		AttachableID   uuid.UUID `json:"attachable_id"`
		Filename       string    `json:"filename"`
		MimeType       string    `json:"mime_type"`
		Size           int64     `json:"size"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}
	if req.Filename == "" {
		return apperrors.NewValidationError(map[string][]string{"filename": {"filename is required"}})
	}

	allowedTypes := map[string]bool{
		"Transaction": true, "Journal": true, "Bill": true,
		"PiggyBank": true, "Recurring": true, "Budget": true,
	}
	if req.AttachableType != "" && !allowedTypes[req.AttachableType] {
		return apperrors.NewValidationError(map[string][]string{
			"attachable_type": {"invalid attachable type"},
		})
	}

	a, err := h.repo.Create(c.Context(), user.ID, req.AttachableType, req.AttachableID, req.Filename, req.MimeType, req.Size)
	if err != nil {
		log.Printf("handler: failed to create attachment: %v", err)
		return apperrors.ErrInternal
	}

	return c.Status(201).JSON(fiber.Map{"data": fiber.Map{
		"type": "attachments",
		"id":   a.ID,
		"attributes": fiber.Map{
			"filename":        a.Filename,
			"mime_type":       a.MimeType,
			"size":            a.Size,
			"uploaded":        a.Uploaded,
			"attachable_type": a.AttachableType,
			"attachable_id":   a.AttachableID,
		},
	}})
}

func (h *AttachmentHandler) Delete(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	if user == nil {
		return apperrors.ErrUnauthorized
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id format"}})
	}

	a, err := h.repo.FindByID(c.Context(), id)
	if err != nil {
		return apperrors.NotFoundResource("attachment", id)
	}
	if a.UserID != user.ID {
		return apperrors.ErrNotFound
	}

	if err := h.repo.Delete(c.Context(), id); err != nil {
		return apperrors.NotFoundResource("attachment", id)
	}

	return c.Status(204).Send(nil)
}

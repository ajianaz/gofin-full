package handler

import (
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
	_ = auth.GetUser(c)

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

	attachments, err := h.repo.ListByEntity(c.Context(), attachableType, attachableID)
	if err != nil {
		return apperrors.NewWithDetail(500, "failed to list attachments", err.Error())
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
	_ = auth.GetUser(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id format"}})
	}

	a, err := h.repo.FindByID(c.Context(), id)
	if err != nil {
		return apperrors.NotFoundResource("attachment", id)
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "attachments",
		"id":   a.ID,
		"attributes": fiber.Map{
			"filename":       a.Filename,
			"mime_type":      a.MimeType,
			"size":           a.Size,
			"uploaded":       a.Uploaded,
			"attachable_type": a.AttachableType,
			"attachable_id":  a.AttachableID,
		},
	}})
}

func (h *AttachmentHandler) Store(c *fiber.Ctx) error {
	user := auth.GetUser(c)

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

	a, err := h.repo.Create(c.Context(), user.ID, req.AttachableType, req.AttachableID, req.Filename, req.MimeType, req.Size)
	if err != nil {
		return apperrors.NewWithDetail(500, "failed to create attachment", err.Error())
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
	_ = auth.GetUser(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id format"}})
	}

	if err := h.repo.Delete(c.Context(), id); err != nil {
		return apperrors.NotFoundResource("attachment", id)
	}

	return c.Status(204).Send(nil)
}

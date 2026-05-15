package handler

import (
"time"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/repository"
	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors")

type TagHandler struct {
	repo *repository.TagRepository
}

func NewTagHandler(repo *repository.TagRepository) *TagHandler {
	return &TagHandler{repo: repo}
}

func (h *TagHandler) Index(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	tags, err := h.repo.List(c.Context(), *groupID)
	if err != nil {
		log.Printf("handler/Index: failed to list tags: %v", err)
		return apperrors.ErrInternal
	}

	var data []fiber.Map
	for _, t := range tags {
		data = append(data, fiber.Map{
			"type":       "tags",
			"id":         t.ID,
			"attributes": fiber.Map{"tag": t.Tag, "date": t.Date, "created_at": t.CreatedAt, "updated_at": t.UpdatedAt},
		})
	}
	return c.JSON(fiber.Map{"data": data})
}

func (h *TagHandler) Show(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id format"}})
	}

	t, err := h.repo.FindByID(c.Context(), id, *groupID)
	if err != nil {
		return apperrors.NotFoundResource("tag", id)
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type":       "tags",
		"id":         t.ID,
		"attributes": fiber.Map{"tag": t.Tag, "date": t.Date, "created_at": t.CreatedAt, "updated_at": t.UpdatedAt},
	}})
}

func (h *TagHandler) Store(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	var req struct {
		Tag  string     `json:"tag"`
		Date *time.Time `json:"date"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}
	if req.Tag == "" {
		return apperrors.NewValidationError(map[string][]string{"tag": {"tag is required"}})
	}

	t, err := h.repo.Create(c.Context(), user.ID, *groupID, req.Tag, req.Date)
	if err != nil {
		log.Printf("handler/Index: failed to create tag: %v", err)
		return apperrors.ErrInternal
	}

	return c.Status(201).JSON(fiber.Map{"data": fiber.Map{
		"type":       "tags",
		"id":         t.ID,
		"attributes": fiber.Map{"tag": t.Tag, "date": t.Date, "created_at": t.CreatedAt, "updated_at": t.UpdatedAt},
	}})
}

func (h *TagHandler) Update(c *fiber.Ctx) error {
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
		Tag  string     `json:"tag"`
		Date *time.Time `json:"date"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}
	if req.Tag == "" {
		return apperrors.NewValidationError(map[string][]string{"tag": {"tag is required"}})
	}

	if err := h.repo.Update(c.Context(), id, *groupID, req.Tag, req.Date); err != nil {
		return apperrors.NotFoundResource("tag", id)
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "tags", "id": id,
		"attributes": fiber.Map{"tag": req.Tag, "date": req.Date},
	}})
}

func (h *TagHandler) Delete(c *fiber.Ctx) error {
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
		return apperrors.NotFoundResource("tag", id)
	}

	return c.Status(204).Send(nil)
}

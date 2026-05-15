package handler

import (
	"log"
"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/repository"
	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors")

type ObjectGroupHandler struct {
	repo *repository.ObjectGroupRepository
}

func NewObjectGroupHandler(repo *repository.ObjectGroupRepository) *ObjectGroupHandler {
	return &ObjectGroupHandler{repo: repo}
}

func (h *ObjectGroupHandler) Index(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	groups, err := h.repo.List(c.Context(), *groupID)
	if err != nil {
		log.Printf("handler/Index: failed to list object groups: %v", err)
		return apperrors.ErrInternal
	}

	var data []fiber.Map
	for _, g := range groups {
		data = append(data, fiber.Map{
			"type": "object_groups",
			"id":   g.ID,
			"attributes": fiber.Map{
				"title": g.Title,
				"order": g.Order,
			},
		})
	}
	return c.JSON(fiber.Map{"data": data})
}

func (h *ObjectGroupHandler) Show(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id format"}})
	}

	g, err := h.repo.FindByID(c.Context(), id, *groupID)
	if err != nil {
		return apperrors.NotFoundResource("object_group", id)
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "object_groups",
		"id":   g.ID,
		"attributes": fiber.Map{
			"title": g.Title,
			"order": g.Order,
		},
	}})
}

func (h *ObjectGroupHandler) Store(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	var req struct {
		Title string `json:"title"`
		Order int    `json:"order"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}
	if req.Title == "" {
		return apperrors.NewValidationError(map[string][]string{"title": {"title is required"}})
	}

	g, err := h.repo.Create(c.Context(), user.ID, *groupID, req.Title, req.Order)
	if err != nil {
		log.Printf("handler/Index: failed to create object group: %v", err)
		return apperrors.ErrInternal
	}

	return c.Status(201).JSON(fiber.Map{"data": fiber.Map{
		"type": "object_groups",
		"id":   g.ID,
		"attributes": fiber.Map{
			"title": g.Title,
			"order": g.Order,
		},
	}})
}

func (h *ObjectGroupHandler) Update(c *fiber.Ctx) error {
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
		Title string `json:"title"`
		Order *int   `json:"order"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}

	if err := h.repo.Update(c.Context(), id, *groupID, req.Title, req.Order); err != nil {
		return apperrors.NotFoundResource("object_group", id)
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "object_groups", "id": id,
		"attributes": fiber.Map{"title": req.Title},
	}})
}

func (h *ObjectGroupHandler) Delete(c *fiber.Ctx) error {
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
		return apperrors.NotFoundResource("object_group", id)
	}

	return c.Status(204).Send(nil)
}

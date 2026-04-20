package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/azfirazka/gofin-full/api/internal/auth"
	"github.com/azfirazka/gofin-full/api/internal/repository"
	apperrors "github.com/azfirazka/gofin-full/api/pkg/errors"
)

type CategoryHandler struct {
	repo *repository.CategoryRepository
}

func NewCategoryHandler(repo *repository.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{repo: repo}
}

func (h *CategoryHandler) Index(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	categories, err := h.repo.List(c.Context(), *groupID)
	if err != nil {
		return apperrors.NewWithDetail(500, "failed to list categories", err.Error())
	}

	var data []fiber.Map
	for _, cat := range categories {
		data = append(data, fiber.Map{
			"type":       "categories",
			"id":         cat.ID,
			"attributes": fiber.Map{"name": cat.Name, "created_at": cat.CreatedAt, "updated_at": cat.UpdatedAt},
		})
	}
	return c.JSON(fiber.Map{"data": data})
}

func (h *CategoryHandler) Show(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id"}})
	}

	cat, err := h.repo.FindByID(c.Context(), int64(id), *groupID)
	if err != nil {
		return apperrors.NotFoundResource("category", int64(id))
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type":       "categories",
		"id":         cat.ID,
		"attributes": fiber.Map{"name": cat.Name, "created_at": cat.CreatedAt, "updated_at": cat.UpdatedAt},
	}})
}

func (h *CategoryHandler) Store(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}
	if req.Name == "" {
		return apperrors.NewValidationError(map[string][]string{"name": {"name is required"}})
	}

	cat, err := h.repo.Create(c.Context(), user.ID, *groupID, req.Name)
	if err != nil {
		return apperrors.NewWithDetail(500, "failed to create category", err.Error())
	}

	return c.Status(201).JSON(fiber.Map{"data": fiber.Map{
		"type":       "categories",
		"id":         cat.ID,
		"attributes": fiber.Map{"name": cat.Name, "created_at": cat.CreatedAt, "updated_at": cat.UpdatedAt},
	}})
}

func (h *CategoryHandler) Update(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id"}})
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}
	if req.Name == "" {
		return apperrors.NewValidationError(map[string][]string{"name": {"name is required"}})
	}

	if err := h.repo.Update(c.Context(), int64(id), *groupID, req.Name); err != nil {
		return apperrors.NotFoundResource("category", int64(id))
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "categories", "id": id,
		"attributes": fiber.Map{"name": req.Name},
	}})
}

func (h *CategoryHandler) Delete(c *fiber.Ctx) error {
	_ = auth.GetUser(c)
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id"}})
	}

	if err := h.repo.Delete(c.Context(), int64(id), *groupID); err != nil {
		return apperrors.NotFoundResource("category", int64(id))
	}

	return c.Status(204).Send(nil)
}

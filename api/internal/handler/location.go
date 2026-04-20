package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/azfirazka/gofin-full/api/internal/auth"
	"github.com/azfirazka/gofin-full/api/internal/repository"
	apperrors "github.com/azfirazka/gofin-full/api/pkg/errors"
)

type LocationHandler struct {
	repo *repository.LocationRepository
}

func NewLocationHandler(repo *repository.LocationRepository) *LocationHandler {
	return &LocationHandler{repo: repo}
}

func (h *LocationHandler) Show(c *fiber.Ctx) error {
	_ = auth.GetUser(c)

	locatableType := c.Query("locatable_type")
	locatableID := c.QueryInt("locatable_id")
	if locatableType == "" || locatableID == 0 {
		return apperrors.NewValidationError(map[string][]string{
			"query": {"locatable_type and locatable_id query params are required"},
		})
	}

	loc, err := h.repo.GetByEntity(c.Context(), locatableType, int64(locatableID))
	if err != nil {
		return apperrors.NotFoundResource("location", 0)
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"type": "locations",
		"id":   loc.ID,
		"attributes": fiber.Map{
			"latitude":       loc.Latitude,
			"longitude":      loc.Longitude,
			"zoom_level":     loc.ZoomLevel,
			"locatable_type": loc.LocatableType,
			"locatable_id":   loc.LocatableID,
		},
	}})
}

func (h *LocationHandler) Store(c *fiber.Ctx) error {
	_ = auth.GetUser(c)

	var req struct {
		LocatableType string   `json:"locatable_type"`
		LocatableID   int64    `json:"locatable_id"`
		Latitude      *float64 `json:"latitude"`
		Longitude     *float64 `json:"longitude"`
		ZoomLevel     int      `json:"zoom_level"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewValidationError(map[string][]string{"body": {"invalid JSON"}})
	}
	if req.LocatableType == "" {
		return apperrors.NewValidationError(map[string][]string{"locatable_type": {"locatable_type is required"}})
	}

	loc, err := h.repo.Set(c.Context(), req.LocatableType, req.LocatableID, req.Latitude, req.Longitude, req.ZoomLevel)
	if err != nil {
		return apperrors.NewWithDetail(500, "failed to set location", err.Error())
	}

	return c.Status(201).JSON(fiber.Map{"data": fiber.Map{
		"type": "locations",
		"id":   loc.ID,
		"attributes": fiber.Map{
			"latitude":       loc.Latitude,
			"longitude":      loc.Longitude,
			"zoom_level":     loc.ZoomLevel,
			"locatable_type": loc.LocatableType,
			"locatable_id":   loc.LocatableID,
		},
	}})
}

func (h *LocationHandler) Delete(c *fiber.Ctx) error {
	_ = auth.GetUser(c)

	id, err := c.ParamsInt("id")
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id"}})
	}

	if err := h.repo.Delete(c.Context(), int64(id)); err != nil {
		return apperrors.NotFoundResource("location", int64(id))
	}

	return c.Status(204).Send(nil)
}

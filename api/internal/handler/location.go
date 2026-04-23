package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/repository"
	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors"
)

type LocationHandler struct {
	repo *repository.LocationRepository
}

func NewLocationHandler(repo *repository.LocationRepository) *LocationHandler {
	return &LocationHandler{repo: repo}
}

func (h *LocationHandler) Show(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	if user == nil {
		return apperrors.ErrUnauthorized
	}
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	locatableType := c.Query("locatable_type")
	locatableIDStr := c.Query("locatable_id")
	if locatableType == "" || locatableIDStr == "" {
		return apperrors.NewValidationError(map[string][]string{
			"query": {"locatable_type and locatable_id query params are required"},
		})
	}

	locatableID, err := uuid.Parse(locatableIDStr)
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{
			"locatable_id": {"invalid locatable_id format"},
		})
	}

	loc, err := h.repo.GetByEntity(c.Context(), locatableType, locatableID)
	if err != nil {
		return apperrors.NotFoundResource("location", uuid.Nil)
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
	user := auth.GetUser(c)
	if user == nil {
		return apperrors.ErrUnauthorized
	}
	groupID := auth.GetActiveGroupID(c)
	if groupID == nil {
		return apperrors.New(400, "no active group")
	}

	var req struct {
		LocatableType string   `json:"locatable_type"`
		LocatableID   uuid.UUID `json:"locatable_id"`
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

	loc, err := h.repo.Set(c.Context(), user.ID, *groupID, req.LocatableType, req.LocatableID, req.Latitude, req.Longitude, req.ZoomLevel)
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
	user := auth.GetUser(c)
	if user == nil {
		return apperrors.ErrUnauthorized
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.NewValidationError(map[string][]string{"id": {"invalid id format"}})
	}

	loc, err := h.repo.FindByID(c.Context(), id)
	if err != nil || loc == nil {
		return apperrors.NotFoundResource("location", id)
	}
	if loc.UserID != user.ID {
		return apperrors.ErrNotFound
	}

	if err := h.repo.Delete(c.Context(), id); err != nil {
		return apperrors.NotFoundResource("location", id)
	}

	return c.Status(204).Send(nil)
}

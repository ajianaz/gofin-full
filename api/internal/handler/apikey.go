package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/repository"
)

// APIKeyHandler handles API key management endpoints.
type APIKeyHandler struct {
	keyRepo *repository.APIKeyRepository
}

// NewAPIKeyHandler creates a new API key handler.
func NewAPIKeyHandler(keyRepo *repository.APIKeyRepository) *APIKeyHandler {
	return &APIKeyHandler{keyRepo: keyRepo}
}

// Create handles POST /api/v1/api-keys.
// Generates a new API key and returns the raw key (shown only once).
func (h *APIKeyHandler) Create(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	// Debug: verify user ID
	_ = user.ID // ensure user is used

	var req struct {
		Name string `json:"name"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(422).JSON(fiber.Map{
			"message": "Invalid request body.",
		})
	}

	if req.Name == "" {
		return c.Status(422).JSON(fiber.Map{
			"message": "Name is required.",
		})
	}

	apiKey, rawKey, err := h.keyRepo.Create(c.Context(), user.ID, req.Name)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to create API key.",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data": fiber.Map{
			"id":         apiKey.ID,
			"name":       apiKey.Name,
			"key_prefix": apiKey.KeyPrefix,
			"key":        rawKey,
			"created_at": apiKey.CreatedAt,
		},
		"message": "Save this key now. It will not be shown again.",
	})
}

// List handles GET /api/v1/api-keys.
// Returns all API keys for the current user (masked — prefix only).
func (h *APIKeyHandler) List(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	keys, err := h.keyRepo.ListByUser(c.Context(), user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to list API keys.",
		})
	}

	var data []fiber.Map
	for _, k := range keys {
		lastUsed := ""
		if k.LastUsedAt != nil {
			lastUsed = k.LastUsedAt.UTC().Format("2006-01-02T15:04:05Z")
		}
		data = append(data, fiber.Map{
			"id":         k.ID,
			"name":       k.Name,
			"key_prefix": k.KeyPrefix,
			"last_used":  lastUsed,
			"created_at": k.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
		})
	}

	return c.JSON(fiber.Map{
		"data": data,
	})
}

// Delete handles DELETE /api/v1/api-keys/:id.
// Soft-deletes the API key.
func (h *APIKeyHandler) Delete(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(422).JSON(fiber.Map{
			"message": "Invalid key ID.",
		})
	}

	if err := h.keyRepo.Delete(c.Context(), id, user.ID); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to delete API key.",
		})
	}

	return c.JSON(fiber.Map{
		"message": "API key deleted.",
	})
}

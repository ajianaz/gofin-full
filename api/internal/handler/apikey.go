package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/dto/response"
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
// JWT auth required — API key auth is not allowed to create new keys (privilege escalation prevention).
func (h *APIKeyHandler) Create(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	// Reject API key auth — only JWT can manage API keys
	if c.Locals("auth_method") == "api_key" {
		return response.SendError(c, 403, "API key management requires full authentication.")
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.SendError(c, 422, "Invalid request body.")
	}

	if req.Name == "" {
		return response.SendError(c, 422, "Name is required.")
	}

	apiKey, rawKey, err := h.keyRepo.Create(c.Context(), user.ID, req.Name)
	if err != nil {
		return response.SendError(c, 500, "Failed to create API key.")
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
		return response.SendError(c, 500, "Failed to list API keys.")
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
// JWT auth required — API key auth cannot delete keys.
func (h *APIKeyHandler) Delete(c *fiber.Ctx) error {
	user := auth.GetUser(c)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	// Reject API key auth — only JWT can manage API keys
	if c.Locals("auth_method") == "api_key" {
		return response.SendError(c, 403, "API key management requires full authentication.")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.SendError(c, 422, "Invalid key ID.")
	}

	if err := h.keyRepo.Delete(c.Context(), id, user.ID); err != nil {
		return response.SendError(c, 500, "Failed to delete API key.")
	}

	return c.JSON(fiber.Map{
		"message": "API key deleted.",
	})
}

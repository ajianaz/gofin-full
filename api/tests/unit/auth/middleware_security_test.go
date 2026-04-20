package auth_test

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ajianaz/gofin-full/api/internal/auth"
)

func newTestAppWithMiddleware(jwtMgr *auth.JWTManager, mw fiber.Handler) *fiber.App {
	app := fiber.New()
	if mw != nil {
		app.Use(mw)
	}
	app.Get("/me", func(c *fiber.Ctx) error {
		groupID := auth.GetActiveGroupID(c)
		if groupID == nil {
			return c.JSON(fiber.Map{"group_id": nil})
		}
		return c.JSON(fiber.Map{"group_id": *groupID})
	})
	return app
}

func generateTestToken(t *testing.T, jwtMgr *auth.JWTManager, userID int64, email string, groupID *int64) string {
	t.Helper()
	identity := &auth.UserIdentity{ID: userID, Email: email}
	pair, err := jwtMgr.GenerateTokenPair(identity, groupID)
	require.NoError(t, err)
	return pair.AccessToken
}

func TestAuthMiddleware_GroupOverrideQueryParamIgnored(t *testing.T) {
	jwtMgr := auth.NewJWTManager("test-secret-32-chars-minimum!!", 60, 30)

	// Token with NO group ID in claims
	token := generateTestToken(t, jwtMgr, 1, "user@example.com", nil)

	app := newTestAppWithMiddleware(jwtMgr, auth.AuthMiddleware(jwtMgr))

	// Try to override group via query param — should be ignored
	req := httptest.NewRequest("GET", "/me?user_group_id=999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var body map[string]interface{}
	bodyBytes, _ := io.ReadAll(resp.Body)
	json.Unmarshal(bodyBytes, &body)

	// group_id should be nil — query param override must be ignored
	assert.Nil(t, body["group_id"], "group override via query param must be ignored")
}

func TestAuthMiddleware_GroupIDFromClaims(t *testing.T) {
	jwtMgr := auth.NewJWTManager("test-secret-32-chars-minimum!!", 60, 30)

	groupID := int64(42)
	token := generateTestToken(t, jwtMgr, 1, "user@example.com", &groupID)

	app := newTestAppWithMiddleware(jwtMgr, auth.AuthMiddleware(jwtMgr))

	req := httptest.NewRequest("GET", "/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var body map[string]interface{}
	bodyBytes, _ := io.ReadAll(resp.Body)
	json.Unmarshal(bodyBytes, &body)

	assert.Equal(t, float64(42), body["group_id"], "group_id should come from JWT claims")
}

func TestAuthMiddleware_QueryParamDoesNotOverrideClaims(t *testing.T) {
	jwtMgr := auth.NewJWTManager("test-secret-32-chars-minimum!!", 60, 30)

	groupID := int64(42)
	token := generateTestToken(t, jwtMgr, 1, "user@example.com", &groupID)

	app := newTestAppWithMiddleware(jwtMgr, auth.AuthMiddleware(jwtMgr))

	// Even with both claims group_id AND query param, query param must not override
	req := httptest.NewRequest("GET", "/me?user_group_id=999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	var body map[string]interface{}
	bodyBytes, _ := io.ReadAll(resp.Body)
	json.Unmarshal(bodyBytes, &body)

	// Should still be 42 from claims, not 999 from query
	assert.Equal(t, float64(42), body["group_id"], "query param must not override JWT claims")
}

func TestOptionalAuthMiddleware_GroupOverrideIgnored(t *testing.T) {
	jwtMgr := auth.NewJWTManager("test-secret-32-chars-minimum!!", 60, 30)

	token := generateTestToken(t, jwtMgr, 1, "user@example.com", nil)

	app := newTestAppWithMiddleware(jwtMgr, auth.OptionalAuthMiddleware(jwtMgr))

	req := httptest.NewRequest("GET", "/me?user_group_id=888", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	var body map[string]interface{}
	bodyBytes, _ := io.ReadAll(resp.Body)
	json.Unmarshal(bodyBytes, &body)

	assert.Nil(t, body["group_id"], "optional auth must also ignore group override query param")
}

func TestOptionalAuthMiddleware_NoTokenPassesThrough(t *testing.T) {
	jwtMgr := auth.NewJWTManager("test-secret-32-chars-minimum!!", 60, 30)

	app := newTestAppWithMiddleware(jwtMgr, auth.OptionalAuthMiddleware(jwtMgr))

	req := httptest.NewRequest("GET", "/me", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var body map[string]interface{}
	bodyBytes, _ := io.ReadAll(resp.Body)
	json.Unmarshal(bodyBytes, &body)

	assert.Nil(t, body["group_id"], "no token means no group set")
}

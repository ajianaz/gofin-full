package integration

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/azfirazka/gofin-full/api/tests/integration/testhelpers"
)

// --- Refresh Token Rotation ---

func TestRefreshToken_Rotation(t *testing.T) {
	app := testApp.App

	// Step 1: Login to get fresh token pair
	// Seed password is "password123" (see testhelpers/database.go line 168)
	loginBody := fmt.Sprintf(`{"email":"%s","password":"password123"}`, testApp.Seed.OwnerEmail)
	loginResp := testhelpers.MakeRequest(t, app, "POST", "/api/v1/auth/login", loginBody, "")
	require.Equal(t, http.StatusOK, loginResp.StatusCode, "login should succeed")

	loginData := testhelpers.ParseResponse(t, loginResp)
	accessToken, ok1 := loginData["access_token"].(string)
	refreshToken, ok2 := loginData["refresh_token"].(string)
	require.True(t, ok1, "login response should have access_token")
	require.True(t, ok2, "login response should have refresh_token")
	require.NotEmpty(t, refreshToken)

	// Step 2: Use refresh token to get a new pair
	refreshBody := fmt.Sprintf(`{"refresh_token":"%s"}`, refreshToken)
	refreshResp := testhelpers.MakeRequest(t, app, "POST", "/api/v1/auth/refresh", refreshBody, "")
	require.Equal(t, http.StatusOK, refreshResp.StatusCode, "refresh should succeed")

	refreshData := testhelpers.ParseResponse(t, refreshResp)
	newAccessToken, ok3 := refreshData["access_token"].(string)
	newRefreshToken, ok4 := refreshData["refresh_token"].(string)
	require.True(t, ok3, "refresh response should have access_token")
	require.True(t, ok4, "refresh response should have refresh_token")
	require.NotEmpty(t, newRefreshToken)

	// New tokens should be different from old ones
	assert.NotEqual(t, accessToken, newAccessToken, "new access token should differ")
	assert.NotEqual(t, refreshToken, newRefreshToken, "new refresh token should differ")

	// Step 3: Old refresh token should be revoked — using it again should fail
	secondRefreshResp := testhelpers.MakeRequest(t, app, "POST", "/api/v1/auth/refresh", refreshBody, "")
	assert.Equal(t, http.StatusUnauthorized, secondRefreshResp.StatusCode,
		"old refresh token should be revoked after rotation")
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	app := testApp.App

	body := `{"refresh_token":"invalid-token-here"}`
	resp := testhelpers.MakeRequest(t, app, "POST", "/api/v1/auth/refresh", body, "")
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode,
		"invalid refresh token should return 401")
}

func TestRefreshToken_MissingField(t *testing.T) {
	app := testApp.App

	resp := testhelpers.MakeRequest(t, app, "POST", "/api/v1/auth/refresh", `{}`, "")
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode,
		"missing refresh_token field should return 422")
}

func TestRefreshToken_EmptyField(t *testing.T) {
	app := testApp.App

	resp := testhelpers.MakeRequest(t, app, "POST", "/api/v1/auth/refresh", `{"refresh_token":""}`, "")
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode,
		"empty refresh_token should return 422")
}

// --- Logout Revocation ---

func TestLogout_RevokesRefreshToken(t *testing.T) {
	app := testApp.App

	// Login to get tokens
	loginBody := fmt.Sprintf(`{"email":"%s","password":"%s"}`, testApp.Seed.OwnerEmail, "password123")
	loginResp := testhelpers.MakeRequest(t, app, "POST", "/api/v1/auth/login", loginBody, "")
	require.Equal(t, http.StatusOK, loginResp.StatusCode)

	loginData := testhelpers.ParseResponse(t, loginResp)
	refreshToken := loginData["refresh_token"].(string)
	require.NotEmpty(t, refreshToken)

	// Logout with the refresh token
	logoutBody := fmt.Sprintf(`{"refresh_token":"%s"}`, refreshToken)
	logoutResp := testhelpers.MakeRequest(t, app, "POST", "/api/v1/auth/logout", logoutBody, "")
	assert.Equal(t, http.StatusOK, logoutResp.StatusCode, "logout should succeed")

	// Try to refresh after logout — should fail
	refreshBody := fmt.Sprintf(`{"refresh_token":"%s"}`, refreshToken)
	refreshResp := testhelpers.MakeRequest(t, app, "POST", "/api/v1/auth/refresh", refreshBody, "")
	assert.Equal(t, http.StatusUnauthorized, refreshResp.StatusCode,
		"refresh token should be revoked after logout")
}

func TestLogout_WithoutToken(t *testing.T) {
	app := testApp.App

	resp := testhelpers.MakeRequest(t, app, "POST", "/api/v1/auth/logout", `{}`, "")
	assert.Equal(t, http.StatusOK, resp.StatusCode,
		"logout without token should still return 200")
}

// --- Refresh Token Rotation Chain ---

func TestRefreshToken_MultipleRotations(t *testing.T) {
	app := testApp.App

	// Login
	loginBody := fmt.Sprintf(`{"email":"%s","password":"%s"}`, testApp.Seed.OwnerEmail, "password123")
	loginResp := testhelpers.MakeRequest(t, app, "POST", "/api/v1/auth/login", loginBody, "")
	require.Equal(t, http.StatusOK, loginResp.StatusCode)

	loginData := testhelpers.ParseResponse(t, loginResp)
	currentRefresh := loginData["refresh_token"].(string)

	// Rotate 3 times — each old token should be invalidated
	for i := 0; i < 3; i++ {
		refreshBody := fmt.Sprintf(`{"refresh_token":"%s"}`, currentRefresh)
		resp := testhelpers.MakeRequest(t, app, "POST", "/api/v1/auth/refresh", refreshBody, "")
		require.Equal(t, http.StatusOK, resp.StatusCode,
			"rotation %d should succeed", i+1)

		data := testhelpers.ParseResponse(t, resp)
		newRefresh, ok := data["refresh_token"].(string)
		require.True(t, ok)
		require.NotEmpty(t, newRefresh)
		require.NotEqual(t, currentRefresh, newRefresh,
			"rotation %d should produce a new token", i+1)

		// Verify old token is dead
		oldResp := testhelpers.MakeRequest(t, app, "POST", "/api/v1/auth/refresh", refreshBody, "")
		assert.Equal(t, http.StatusUnauthorized, oldResp.StatusCode,
			"old token from rotation %d should be rejected", i+1)

		currentRefresh = newRefresh
	}
}

// --- First-User-Only Global Owner ---

func TestFirstUserOnlyOwner(t *testing.T) {
	app := testApp.App

	// First user (owner) should be able to access admin endpoints
	t.Run("first_user_is_admin", func(t *testing.T) {
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "GET", "/api/v1/admin/users", "", testApp.Seed.OwnerToken)
		assert.Equal(t, http.StatusOK, resp.StatusCode,
			"first user (owner) should access admin endpoints")
	})

	// Create a second user via admin endpoint
	t.Run("admin_creates_second_user", func(t *testing.T) {
		email := fmt.Sprintf("second-user-%d@gofin.io", testApp.Seed.GroupID)
		body := fmt.Sprintf(`{"email":"%s","password":"password1234"}`, email)
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "POST", "/api/v1/admin/users", body, testApp.Seed.OwnerToken)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	// Login as second user
	t.Run("second_user_can_login", func(t *testing.T) {
		email := fmt.Sprintf("second-user-%d@gofin.io", testApp.Seed.GroupID)
		body := fmt.Sprintf(`{"email":"%s","password":"password1234"}`, email)
		resp := testhelpers.MakeRequest(t, app, "POST", "/api/v1/auth/login", body, "")
		require.Equal(t, http.StatusOK, resp.StatusCode)

		data := testhelpers.ParseResponse(t, resp)
		secondToken, ok := data["access_token"].(string)
		require.True(t, ok)
		require.NotEmpty(t, secondToken)

		// Second user should NOT be able to access admin endpoints
		adminResp := testhelpers.MakeAuthenticatedRequest(t, app, "GET", "/api/v1/admin/users", "", secondToken)
		assert.Equal(t, http.StatusForbidden, adminResp.StatusCode,
			"second user should NOT have admin access (not a global owner)")
	})
}

// --- Admin Role Verification ---

func TestAdminEndpoints_RequireAdminRole(t *testing.T) {
	app := testApp.App

	t.Run("feature_flags_require_admin", func(t *testing.T) {
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "GET", "/api/v1/admin/feature-flags", "", testApp.Seed.OwnerToken)
		assert.Equal(t, http.StatusOK, resp.StatusCode,
			"owner should access feature flags")
	})
}

// Helper to decode response JSON
func decodeJSON(t *testing.T, resp *http.Response) map[string]interface{} {
	t.Helper()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)
	return result
}

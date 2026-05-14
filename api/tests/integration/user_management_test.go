package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ajianaz/gofin-full/api/tests/integration/testhelpers"
)

func TestAdminCreateUser(t *testing.T) {
	app := testApp.App
	token := testApp.Seed.OwnerToken

	t.Run("admin_can_create_user", func(t *testing.T) {
		email := fmt.Sprintf("newuser-%s@gofin.io", testApp.Seed.GroupID)
		body := fmt.Sprintf(`{"email":"%s","password": "SecurePass1!"}`, email)
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "POST", "/api/v1/admin/users", body, token)
		require.Equal(t, http.StatusCreated, resp.StatusCode,
			"admin should be able to create a user")
	})

	t.Run("created_user_can_login", func(t *testing.T) {
		email := fmt.Sprintf("login-test-%s@gofin.io", testApp.Seed.GroupID)
		body := fmt.Sprintf(`{"email":"%s","password": "SecurePass1!"}`, email)

		// Create user via admin
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "POST", "/api/v1/admin/users", body, token)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		// Login with the new user
		loginResp := testhelpers.MakeRequest(t, app, "POST", "/api/v1/auth/login", body, "")
		require.Equal(t, http.StatusOK, loginResp.StatusCode,
			"newly created user should be able to login")
	})

	t.Run("admin_create_duplicate_email_returns_409", func(t *testing.T) {
		body := fmt.Sprintf(`{"email":"%s","password": "SecurePass1!"}`, testApp.Seed.OwnerEmail)
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "POST", "/api/v1/admin/users", body, token)
		require.Equal(t, http.StatusConflict, resp.StatusCode,
			"duplicate email should return 409")
	})

	t.Run("admin_create_missing_email_returns_422", func(t *testing.T) {
		body := `{"password":"securepass123"}`
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "POST", "/api/v1/admin/users", body, token)
		require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode,
			"missing email should return 422")
	})

	t.Run("admin_create_short_password_returns_422", func(t *testing.T) {
		body := `{"email":"short@gofin.io","password":"abc"}`
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "POST", "/api/v1/admin/users", body, token)
		require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode,
			"short password should return 422")
	})

	t.Run("read_only_cannot_create_user", func(t *testing.T) {
		token := testApp.Seed.ReadOnlyToken
		body := `{"email":"blocked@gofin.io","password":"securepass123"}`
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "POST", "/api/v1/admin/users", body, token)
		require.Equal(t, http.StatusForbidden, resp.StatusCode,
			"read_only should not be able to create users")
	})
}

func TestSelfRegistration(t *testing.T) {
	app := testApp.App

	t.Run("self_register_enabled_creates_user_and_returns_tokens", func(t *testing.T) {
		email := fmt.Sprintf("register-%s@gofin.io", testApp.Seed.GroupID)
		body := fmt.Sprintf(`{"email":"%s","password": "SecurePass1!"}`, email)
		resp := testhelpers.MakeRequest(t, app, "POST", "/api/v1/auth/register", body, "")
		require.Equal(t, http.StatusCreated, resp.StatusCode,
			"self-registration should work when enabled")
	})

	t.Run("self_register_duplicate_email_returns_409", func(t *testing.T) {
		body := fmt.Sprintf(`{"email":"%s","password": "SecurePass1!"}`, testApp.Seed.OwnerEmail)
		resp := testhelpers.MakeRequest(t, app, "POST", "/api/v1/auth/register", body, "")
		require.Equal(t, http.StatusConflict, resp.StatusCode,
			"duplicate email on registration should return 409")
	})

	t.Run("self_register_missing_email_returns_422", func(t *testing.T) {
		body := `{"password":"registerpass1"}`
		resp := testhelpers.MakeRequest(t, app, "POST", "/api/v1/auth/register", body, "")
		require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode,
			"missing email should return 422")
	})

	t.Run("self_register_short_password_returns_422", func(t *testing.T) {
		body := `{"email":"short@gofin.io","password":"abc"}`
		resp := testhelpers.MakeRequest(t, app, "POST", "/api/v1/auth/register", body, "")
		require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode,
			"short password should return 422")
	})

	t.Run("self_register_invalid_body_returns_422", func(t *testing.T) {
		resp := testhelpers.MakeRequest(t, app, "POST", "/api/v1/auth/register", "not json", "")
		require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
}

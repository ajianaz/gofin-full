package testhelpers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	"github.com/azfirazka/gofin-full/api/internal/auth"
)

// MakeRequest sends a request through the Fiber app using httptest.
// If token is non-empty the Authorization header is set.
func MakeRequest(t *testing.T, app *fiber.App, method, path, body, token string) *http.Response {
	t.Helper()

	var bodyReader io.Reader
	if body != "" {
		bodyReader = bytes.NewBufferString(body)
	}

	req := httptest.NewRequest(method, path, bodyReader)
	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := app.Test(req, -1) // -1 = no timeout
	require.NoError(t, err, "app.Test should not error for %s %s", method, path)

	return resp
}

// MakeAuthenticatedRequest always sends the request with an Authorization header.
func MakeAuthenticatedRequest(t *testing.T, app *fiber.App, method, path, body, token string) *http.Response {
	t.Helper()
	return MakeRequest(t, app, method, path, body, token)
}

// ParseResponse reads the response body and unmarshals it into a map.
func ParseResponse(t *testing.T, resp *http.Response) map[string]interface{} {
	t.Helper()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "reading response body should not error")
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	require.NoError(t, err, "response body should be valid JSON: %s", string(body))

	return result
}

// ParseResponseBytes returns the raw response body bytes.
func ParseResponseBytes(t *testing.T, resp *http.Response) []byte {
	t.Helper()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "reading response body should not error")
	defer resp.Body.Close()

	return body
}

// GenerateTestToken creates a valid JWT access token for the given user and group.
func GenerateTestToken(jwtMgr *auth.JWTManager, userID int64, email string, groupID *int64) string {
	identity := &auth.UserIdentity{
		ID:    userID,
		Email: email,
	}
	pair, err := jwtMgr.GenerateTokenPair(identity, groupID)
	if err != nil {
		panic("failed to generate test token: " + err.Error())
	}
	return pair.AccessToken
}

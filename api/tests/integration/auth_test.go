package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ajianaz/gofin-full/api/tests/integration/testhelpers"
)

func TestUnauthenticatedAccess(t *testing.T) {
	app := testApp.App

	t.Run("public_endpoints_return_ok", func(t *testing.T) {
		publicEndpoints := []struct {
			method string
			path   string
		}{
			{"GET", "/health"},
			{"GET", "/api/v1"},
			{"GET", "/api/v1/auth/provider"},
		}

		for _, ep := range publicEndpoints {
			ep := ep
			t.Run(ep.method+" "+ep.path, func(t *testing.T) {
				resp := testhelpers.MakeRequest(t, app, ep.method, ep.path, "", "")
				if ep.path == "/health" {
					// Health may return 503 when Redis is unavailable in test env
					assert.Contains(t, []int{http.StatusOK, http.StatusServiceUnavailable}, resp.StatusCode)
				} else {
					require.Equal(t, http.StatusOK, resp.StatusCode,
						"public endpoint %s %s should return 200", ep.method, ep.path)
				}
			})
		}
	})

	t.Run("protected_endpoints_return_401", func(t *testing.T) {
		protectedEndpoints := []struct {
			method string
			path   string
			body   string
		}{
			{"GET", "/api/v1/users/me", ""},
			{"GET", "/api/v1/wallets", ""},
			{"POST", "/api/v1/transactions", `{"type":"withdrawal"}`},
			{"GET", "/api/v1/groups", ""},
			{"GET", "/api/v1/categories", ""},
			{"GET", "/api/v1/budgets", ""},
			{"GET", "/api/v1/tags", ""},
		}

		for _, ep := range protectedEndpoints {
			ep := ep
			t.Run(ep.method+" "+ep.path, func(t *testing.T) {
				resp := testhelpers.MakeRequest(t, app, ep.method, ep.path, ep.body, "")
				require.Equal(t, http.StatusUnauthorized, resp.StatusCode,
					"protected endpoint %s %s should return 401 without token", ep.method, ep.path)
			})
		}
	})

	t.Run("invalid_token_returns_401", func(t *testing.T) {
		resp := testhelpers.MakeRequest(t, app, "GET", "/api/v1/users/me", "", "invalid-jwt-token")
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

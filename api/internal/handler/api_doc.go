package handler

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

// openapiSpec holds the spec loaded from disk (YAML) at init time.
// If loading fails, it stays nil and OpenAPISpec falls back to inline JSON.
var openapiSpec []byte

func init() {
	data, err := os.ReadFile("docs/openapi.yaml")
	if err == nil {
		openapiSpec = data
	}
}

// APIDocHandler serves the OpenAPI spec and API documentation.
type APIDocHandler struct{}

func NewAPIDocHandler() *APIDocHandler {
	return &APIDocHandler{}
}

// OpenAPISpec returns the OpenAPI 3.0 spec.
// If docs/openapi.yaml was loaded at init time, it is served as YAML;
// otherwise the built-in inline JSON spec is used as a fallback.
func (h *APIDocHandler) OpenAPISpec(c *fiber.Ctx) error {
	if openapiSpec != nil {
		c.Set("Content-Type", "application/yaml")
		return c.Send(openapiSpec)
	}
	c.Set("Content-Type", "application/json")
	return c.Send(getFallbackOpenAPISpec())
}

// APIDocs returns a simple HTML documentation page.
func (h *APIDocHandler) APIDocs(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.Send(getAPIDocsHTML())
}

// getFallbackOpenAPISpec returns a minimal inline JSON spec used when
// docs/openapi.yaml is not available on disk.
func getFallbackOpenAPISpec() []byte {
	return []byte(`{
  "openapi": "3.0.3",
  "info": {
    "title": "Gofin API",
    "description": "Personal finance management API — a Go rewrite of Firefly III",
    "version": "1.0.0",
    "contact": { "name": "Gofin" },
    "license": { "name": "MIT" }
  },
  "servers": [
    { "url": "/api/v1", "description": "API v1" }
  ],
  "paths": {
    "/health": {
      "get": {
        "summary": "Health check",
        "description": "Returns service health status including database and Redis connectivity",
        "tags": ["System"],
        "responses": { "200": { "description": "OK" } }
      }
    },
    "/auth/login": {
      "post": {
        "summary": "Login",
        "description": "Authenticate user and return JWT tokens",
        "tags": ["Auth"],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "email": { "type": "string" },
                  "password": { "type": "string" }
                },
                "required": ["email", "password"]
              }
            }
          }
        },
        "responses": {
          "200": { "description": "OK with JWT tokens" },
          "401": { "description": "Invalid credentials" }
        }
      }
    },
    "/auth/refresh": {
      "post": {
        "summary": "Refresh token",
        "description": "Exchange a refresh token for new access token",
        "tags": ["Auth"],
        "responses": {
          "200": { "description": "OK with new JWT" },
          "401": { "description": "Invalid refresh token" }
        }
      }
    },
    "/auth/provider": {
      "get": {
        "summary": "Get auth provider",
        "description": "Returns the configured authentication provider",
        "tags": ["Auth"],
        "responses": { "200": { "description": "OK" } }
      }
    }
  },
  "components": {
    "securitySchemes": {
      "bearerAuth": {
        "type": "http",
        "scheme": "bearer",
        "bearerFormat": "JWT"
      }
    }
  }
}`)
}

func getAPIDocsHTML() []byte {
	return []byte(`<!DOCTYPE html>
<html>
<head>
  <title>Gofin API Documentation</title>
  <style>
    body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; max-width: 900px; margin: 0 auto; padding: 40px 20px; color: #1a1a1a; line-height: 1.6; }
    h1 { border-bottom: 2px solid #2563eb; padding-bottom: 10px; }
    h2 { color: #2563eb; margin-top: 30px; }
    code { background: #f1f5f9; padding: 2px 6px; border-radius: 4px; font-size: 0.9em; }
    .endpoint { background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 8px; padding: 16px; margin: 12px 0; }
    .method { font-weight: bold; color: #059669; font-family: monospace; }
    .path { font-family: monospace; font-weight: 600; }
    table { width: 100%; border-collapse: collapse; margin: 12px 0; }
    th, td { text-align: left; padding: 8px 12px; border-bottom: 1px solid #e2e8f0; }
    th { background: #f8fafc; }
  </style>
</head>
<body>
  <h1>Gofin API v1</h1>
  <p>Personal finance management API. All protected endpoints require a JWT bearer token in the Authorization header.</p>
  <p><a href="/api/v1/openapi.json">Open / OpenAPI Spec</a></p>

  <h2>Authentication</h2>
  <div class="endpoint">
    <p><span class="method">POST</span> <span class="path">/api/v1/auth/login</span></p>
    <p>Login with email/password to receive access and refresh tokens.</p>
    <table><tr><th>Field</th><th>Type</th><th>Description</th></tr>
    <tr><td>email</td><td>string</td><td>User email</td></tr>
    <tr><td>password</td><td>string</td><td>User password</td></tr></table>
  </div>

  <div class="endpoint">
    <p><span class="method">POST</span> <span class="path">/api/v1/auth/refresh</span></p>
    <p>Exchange refresh token for new access token.</p>
  </div>

  <h2>Users</h2>
  <div class="endpoint">
    <p><span class="method">GET</span> <span class="path">/api/v1/users/me</span></p>
    <p>Get current user profile.</p>
  </div>
  <div class="endpoint">
    <p><span class="method">PUT</span> <span class="path">/api/v1/users/me</span></p>
    <p>Update current user profile.</p>
  </div>

  <h2>Groups</h2>
  <div class="endpoint">
    <p><span class="method">GET</span> <span class="path">/api/v1/groups</span></p>
    <p>List all user groups.</p>
  </div>
  <div class="endpoint">
    <p><span class="method">POST</span> <span class="path">/api/v1/groups</span></p>
    <p>Create a new group.</p>
  </div>
  <div class="endpoint">
    <p><span class="method">POST</span> <span class="path">/api/v1/groups/switch</span></p>
    <p>Switch active group context.</p>
  </div>

  <h2>Accounts</h2>
  <div class="endpoint">
    <p><span class="method">GET</span> <span class="path">/api/v1/accounts</span></p>
    <p>List all financial accounts (wallets).</p>
  </div>
  <div class="endpoint">
    <p><span class="method">POST</span> <span class="path">/api/v1/accounts</span></p>
    <p>Create a new account. Requires manage_meta role.</p>
  </div>

  <h2>Transactions</h2>
  <div class="endpoint">
    <p><span class="method">GET</span> <span class="path">/api/v1/transactions?page=1&limit=50</span></p>
    <p>List transactions with pagination and filtering.</p>
  </div>
  <div class="endpoint">
    <p><span class="method">POST</span> <span class="path">/api/v1/transactions</span></p>
    <p>Create a double-entry transaction (withdrawal, deposit, or transfer).</p>
  </div>
  <div class="endpoint">
    <p><span class="method">POST</span> <span class="path">/api/v1/transactions/split</span></p>
    <p>Create a split transaction with multiple journals.</p>
  </div>

  <h2>Categories</h2>
  <div class="endpoint">
    <p><span class="method">GET</span> <span class="path">/api/v1/categories</span></p>
    <p>List all categories.</p>
  </div>
  <div class="endpoint">
    <p><span class="method">POST</span> <span class="path">/api/v1/categories</span></p>
    <p>Create a new category. Requires manage_meta role.</p>
  </div>

  <h2>Tags</h2>
  <div class="endpoint">
    <p><span class="method">GET</span> <span class="path">/api/v1/tags</span></p>
    <p>List all tags.</p>
  </div>
  <div class="endpoint">
    <p><span class="method">POST</span> <span class="path">/api/v1/tags</span></p>
    <p>Create a new tag. Requires manage_meta role.</p>
  </div>

  <h2>Budgets</h2>
  <div class="endpoint">
    <p><span class="method">GET</span> <span class="path">/api/v1/budgets</span></p>
    <p>List all budgets.</p>
  </div>
  <div class="endpoint">
    <p><span class="method">POST</span> <span class="path">/api/v1/budgets</span></p>
    <p>Create a new budget. Requires manage_budgets role.</p>
  </div>

  <h2>Webhooks</h2>
  <div class="endpoint">
    <p><span class="method">GET</span> <span class="path">/api/v1/webhooks</span></p>
    <p>List all webhooks.</p>
  </div>
  <div class="endpoint">
    <p><span class="method">POST</span> <span class="path">/api/v1/webhooks</span></p>
    <p>Create a new webhook. Requires manage_webhooks role.</p>
  </div>

  <h2>Analytics</h2>
  <div class="endpoint">
    <p><span class="method">GET</span> <span class="path">/api/v1/analytics/spending-by-category</span></p>
    <p>Get spending totals grouped by category for a date range.</p>
  </div>
  <div class="endpoint">
    <p><span class="method">GET</span> <span class="path">/api/v1/analytics/spending-by-period</span></p>
    <p>Get spending totals grouped by time period.</p>
  </div>
  <div class="endpoint">
    <p><span class="method">GET</span> <span class="path">/api/v1/analytics/net-worth</span></p>
    <p>Get net worth calculation (assets minus liabilities).</p>
  </div>

  <h2>Export</h2>
  <div class="endpoint">
    <p><span class="method">GET</span> <span class="path">/api/v1/export/csv</span></p>
    <p>Export transactions as CSV (Firefly III compatible).</p>
  </div>
  <div class="endpoint">
    <p><span class="method">GET</span> <span class="path">/api/v1/export/ofx</span></p>
    <p>Export transactions in OFX format.</p>
  </div>

  <h2>Audit Logs</h2>
  <div class="endpoint">
    <p><span class="method">GET</span> <span class="path">/api/v1/audit-logs</span></p>
    <p>List audit trail entries. Requires view_memberships role.</p>
  </div>

  <h2>Admin</h2>
  <div class="endpoint">
    <p><span class="method">GET</span> <span class="path">/api/v1/admin/users</span></p>
    <p>List all users in the system. Admin only.</p>
  </div>
  <div class="endpoint">
    <p><span class="method">GET</span> <span class="path">/api/v1/admin/feature-flags</span></p>
    <p>List all feature flags. Admin only.</p>
  </div>
  <div class="endpoint">
    <p><span class="method">POST</span> <span class="path">/api/v1/admin/feature-flags</span></p>
    <p>Create or update a feature flag. Admin only.</p>
  </div>
</body>
</html>`)
}

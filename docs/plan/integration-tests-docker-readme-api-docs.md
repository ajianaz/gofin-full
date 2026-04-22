# Plan: Integration Tests, Docker, README, API Docs

## Context
All 15 sprints are complete. The codebase has 80+ API endpoints, a two-tier RBAC system (21 group roles + 3 wallet roles), but the RBAC middleware is defined but **not wired into the router**. Integration test directories exist but are empty. Docker files exist but lack `AUTH_JWT_SECRET` and migration auto-run. No README or proper OpenAPI spec exists.

Branch: `feat/testing-docker-docs` from `feat/sprint-15`

---

## Step 1: Branch Setup
- Create `feat/testing-docker-docs` from `feat/sprint-15`

## Step 2: Docker Improvements
- **`.dockerignore`** — exclude tests, docs, .git, bin, etc.
- **`deployments/docker/Dockerfile`** — add migrate binary build, copy migrations, add entrypoint script
- **`deployments/docker/entrypoint.sh`** — run migrations then start server
- **`deployments/docker/docker-compose.yml`** — add `AUTH_JWT_SECRET`, `AUTH_PROVIDER=disabled`
- **`deployments/docker/docker-compose.test.yml`** — lightweight PG + Redis only for integration tests (ports 5433/6380)

## Step 3: Wire RBAC Middleware into Router
**Modify `internal/router/router.go`:**
- Split protected routes into read (no RBAC) and write (with RBAC) groups
- Write routes get `auth.RBACMiddleware(appropriateRole)`:
  - Accounts/Categories/Tags write → `RoleManageMeta`
  - Transactions write → `RoleManageTransactions`
  - Budgets write → `RoleManageBudgets`
  - Piggy banks write → `RoleManagePiggyBanks`
  - Rules write → `RoleManageRules`
  - Recurrences write → `RoleManageRecurring`
  - Bills write → `RoleManageTransactions`
  - Webhooks write → `RoleManageWebhooks`
  - Group delete → `RoleOwner`

## Step 4: Integration Test Infrastructure
**Create `tests/integration/testhelpers/`:**
- `config.go` — test config with safe defaults, env overrides
- `database.go` — seed users/groups/wallets/members, truncate tables
- `http.go` — authenticated/unauthenticated request helpers, token generation
- `app.go` — boot full Fiber app with real DB/Redis, run migrations in-process

**Create `tests/integration/main_test.go`:**
- `TestMain` that boots test app once (migrations run once)

## Step 5: Integration Test Scenarios
**`tests/integration/auth_test.go`:**
- Unauthenticated: all protected endpoints return 401, public endpoints return 200
- Authenticated with no group: /users/me works, group-scoped endpoints return error

**`tests/integration/rbac_test.go`:**
- `owner` role: full CRUD on all resources (accounts, transactions, budgets, categories, tags)
- `full` role: same as owner except group delete
- `read_only` role: GET succeeds, POST/PUT/DELETE returns 403
- `manage_transactions` role: can CRUD transactions, cannot create budgets
- `manage_budgets` role: can CRUD budgets + transactions (cascade), cannot manage rules
- Role cascade: `manage_budgets` (level 5) includes `manage_transactions` (level 2)

**`tests/integration/wallet_member_test.go`:**
- Owner can share wallet (POST /members)
- Editor can read + write transactions
- Viewer can only read
- Non-member gets 403/404

## Step 6: README.md
Create comprehensive English README with:
- Project overview (personal finance tracker)
- Features, tech stack
- Quick start with Docker (`docker compose up`)
- Local development setup
- API docs reference
- Project structure tree
- Testing commands
- License (MIT)

## Step 7: OpenAPI Specification
**Create `docs/openapi.yaml`:**
- OpenAPI 3.0.3 spec covering all 80+ endpoints
- Request/response schemas matching JSON:API envelope format
- Bearer JWT security scheme
- Standard error responses (400, 401, 403, 404, 422, 500)

**Update `internal/handler/api_doc.go`:**
- Embed and serve the external YAML spec instead of inline JSON

## Step 8: Update Makefile
- Add `test-integration-infra` / `test-integration-teardown` targets
- Update `test-integration` to use docker-compose.test.yml with proper env vars

## Step 9: Validate
- `go build ./...`
- `go vet ./...`
- `go test ./tests/unit/... -count=1`
- `docker compose -f deployments/docker/docker-compose.yml up -d --build` → verify health check
- `make test-integration` → verify all RBAC scenarios pass

---

## Files to Create (~14)
1. `.dockerignore`
2. `deployments/docker/entrypoint.sh`
3. `deployments/docker/docker-compose.test.yml`
4. `tests/integration/testhelpers/config.go`
5. `tests/integration/testhelpers/database.go`
6. `tests/integration/testhelpers/http.go`
7. `tests/integration/testhelpers/app.go`
8. `tests/integration/main_test.go`
9. `tests/integration/auth_test.go`
10. `tests/integration/rbac_test.go`
11. `tests/integration/wallet_member_test.go`
12. `README.md`
13. `docs/openapi.yaml`

## Files to Modify (~6)
1. `internal/router/router.go` — wire RBACMiddleware
2. `internal/handler/api_doc.go` — serve external spec
3. `deployments/docker/docker-compose.yml` — add auth env vars
4. `deployments/docker/Dockerfile` — add migrate + entrypoint
5. `Makefile` — add test targets
6. `.env.example` — add AUTH fields

---

## Key References

### RBAC System
- `internal/auth/rbac.go` — 21 GroupRole levels, `HasPermission()`, `RBACMiddleware()`
- `internal/auth/middleware.go` — `AuthMiddleware` sets user/group in context
- `internal/auth/context.go` — `SetUser/GetUser`, `SetClaims/GetClaims`, `SetActiveGroupID`
- `internal/auth/jwt.go` — `JWTManager.GenerateTokenPair(identity, groupID)`

### Group Role Hierarchy (level 1→21)
```
read_only(1) → manage_transactions(2) → manage_meta(3) → read_budgets(4) → manage_budgets(5)
→ read_piggy_banks(6) → manage_piggy_banks(7) → read_subscriptions(8) → manage_subscriptions(9)
→ read_rules(10) → manage_rules(11) → read_recurring(12) → manage_recurring(13)
→ read_webhooks(14) → manage_webhooks(15) → read_currencies(16) → manage_currencies(17)
→ view_reports(18) → view_memberships(19) → full(20) → owner(21)
```

### Wallet Member Roles
- `owner` — full CRUD + manage members + delete
- `editor` — create/edit/delete transactions
- `viewer` — read-only

### Existing Docker
- `deployments/docker/Dockerfile` — multi-stage build (golang:1.25-alpine + alpine:3.21)
- `deployments/docker/docker-compose.yml` — app + postgres:16 + redis:7 + keycloak:24

### Config
- `internal/config/config.go` — viper-based, env vars with defaults
- `.env.example` — all env vars documented

### Test Patterns
- `tests/unit/middleware/middleware_test.go` — `setupTestApp()` with `fiber.New()` + `httptest`
- `tests/unit/handler/router_test.go` — `newTestRouter()` pattern
- Framework: `github.com/stretchr/testify`

### Migration Format
- `migrations/postgres/000001` through `000010` — goose format (`+goose Up` / `+goose Down`)
- Runner: `cmd/migrate/main.go` — sorts files, tracks in `schema_migrations` table

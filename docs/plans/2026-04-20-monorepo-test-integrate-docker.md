# Monorepo Test, Integration, Docker Dev Environment Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Establish a fully tested, Docker-based development environment for the gofin-full monorepo with API tests passing, web app building with bun, API+web integration verified, and Playwright E2E UI tests running.

**Architecture:** Docker Compose dev stack runs PostgreSQL + Redis + API + Web together. Unit tests run locally or in Docker. Integration tests require the Docker test stack. Playwright E2E tests run against the live dev stack. Bun replaces npm as the web package manager. CLAUDE.md provides AI-assisted development guardrails.

**Tech Stack:** Go 1.25, Fiber, PostgreSQL 16, Redis 7, SvelteKit 5, Svelte 5, Tailwind CSS 4, Bun, Playwright, Docker Compose, Caddy

---

## Task 1: Create CLAUDE.md

**Files:**
- Create: `CLAUDE.md`

**Step 1: Write CLAUDE.md**

```markdown
# Gofin Full — Claude Code Guidelines

## Project Overview
Personal finance tracker monorepo. Go API backend + SvelteKit frontend.

## Structure
- `api/` — Go backend (module: github.com/ajianaz/gofin-full/api)
- `web/` — SvelteKit frontend (bun, not npm)
- `deployments/docker/` — All Docker configs
- `docs/` — OpenAPI spec, research, plans

## Rules
- **No `Co-Authored-By` in commits.** Just write the commit message.
- Use `bun` for web, not npm/node.
- API runs on port 8080, web dev on port 5173.
- Integration tests require Docker (postgres:5433, redis:6380).
- Unit tests run without Docker.
- Module path is `github.com/ajianaz/gofin-full/api` (NOT ajianaz, NOT azfirazka).

## Commands
- `make api-test-unit` — Unit tests only
- `make api-test-integration` — Integration tests (needs Docker)
- `make docker-up` — Start dev stack (postgres, redis, keycloak, api)
- `make web-dev` — Start SvelteKit dev server
- `make docker-test` — Run full test suite in Docker

## Commit Style
- `feat: description` — New feature
- `fix: description` — Bug fix
- `chore: description` — Maintenance
- `docs: description` — Documentation
- `refactor: description` — Code restructuring

## Testing
- Always run unit tests after changes: `cd api && go test ./tests/unit/... -count=1`
- Integration tests need: `make api-test-integration-infra` then `make api-test-integration`
- Web type check: `cd web && bun run check`
```

**Step 2: Commit**

```bash
git add CLAUDE.md
git commit -m "chore: add CLAUDE.md for AI-assisted development guidelines"
```

---

## Task 2: Fix .env.example (firefly → gofin)

**Files:**
- Modify: `.env.example:23-26,44-46`

**Step 1: Update database defaults**

Replace `DB_DATABASE=firefly` → `DB_DATABASE=gofin`, `DB_USERNAME=firefly` → `DB_USERNAME=gofin`, `DB_PASSWORD=firefly_secret` → `DB_PASSWORD=gofin_secret`.

Replace `KEYCLOAK_REALM=firefly` → `KEYCLOAK_REALM=gofin`, `KEYCLOAK_CLIENT_ID=firefly-go-api` → `KEYCLOAK_CLIENT_ID=gofin-api`.

**Step 2: Add DB_DSN and web env vars**

Add to the Database section:
```
DB_DSN=postgres://gofin:gofin_secret@localhost:5432/gofin?sslmode=prefer
```

Add a new Web section:
```
# Web (SvelteKit)
ORIGIN=http://localhost:5173
PUBLIC_API_BASE=/api/v1
```

**Step 3: Verify**

Run: `grep -n "firefly" .env.example`
Expected: 0 results

**Step 4: Commit**

```bash
git add .env.example
git commit -m "fix: correct .env.example defaults from firefly to gofin"
```

---

## Task 3: Run API Unit Tests

**Files:**
- Read: `api/tests/unit/**/*.go` (all 17 test files)

**Step 1: Run unit tests**

Run: `cd api && go test ./tests/unit/... -v -count=1 2>&1`

Expected: All PASS. If any FAIL, diagnose and fix the issue in the same task before committing.

**Step 2: If tests fail — fix and re-run**

Common issues to check:
- Import path mismatches (should all be `github.com/ajianaz/gofin-full/api/...`)
- Missing dependencies → `cd api && go mod tidy`

**Step 3: Commit (only if fixes were needed)**

```bash
git add api/
git commit -m "fix: resolve unit test failures after monorepo migration"
```

---

## Task 4: Run API Integration Tests via Docker

**Files:**
- Modify: `deployments/docker/docker-compose.test.yml` (if needed)
- Read: `api/tests/integration/main_test.go`

**Step 1: Start test infrastructure**

Run: `cd gofin-full && docker compose -f deployments/docker/docker-compose.test.yml up -d postgres redis`

Wait for healthy: `docker compose -f deployments/docker/docker-compose.test.yml ps`

**Step 2: Run integration tests**

Run: `cd api && AUTH_PROVIDER=disabled AUTH_JWT_SECRET=test-jwt-secret-for-integration-tests-32ch DB_HOST=localhost DB_PORT=5433 DB_DATABASE=gofin_test DB_USERNAME=gofin_test DB_PASSWORD=gofin_test REDIS_HOST=localhost REDIS_PORT=6380 go test ./tests/integration/... -v -count=1 2>&1`

Expected: All PASS.

**Step 3: If tests fail — diagnose and fix**

Check migration path in `api/tests/integration/testhelpers/database.go` — should be `filepath.Join("..", "..", "migrations", "postgres")` which resolves correctly within `api/tests/integration/testhelpers/`.

**Step 4: Teardown test infrastructure**

Run: `cd gofin-full && docker compose -f deployments/docker/docker-compose.test.yml down -v`

**Step 5: Commit (only if fixes were needed)**

```bash
git add api/ deployments/
git commit -m "fix: resolve integration test failures in monorepo"
```

---

## Task 5: Run Full API Test Suite in Docker

**Files:**
- Read: `deployments/docker/Dockerfile.api.test`
- Read: `deployments/docker/docker-compose.test.yml`

**Step 1: Run full test suite via Docker**

Run: `cd gofin-full && docker compose -f deployments/docker/docker-compose.test.yml up --build --abort-on-container-exit 2>&1`

Expected: Container runs unit tests then integration tests, all PASS, container exits 0.

**Step 2: If build fails — fix Dockerfile.api.test**

The Dockerfile.api.test builds from `context: ../../api` with `dockerfile: ../deployments/docker/Dockerfile.api.test`. Verify paths are correct.

**Step 3: Teardown**

Run: `cd gofin-full && docker compose -f deployments/docker/docker-compose.test.yml down -v`

**Step 4: Commit (only if fixes were needed)**

```bash
git add deployments/
git commit -m "fix: resolve Docker test suite issues"
```

---

## Task 6: Migrate Web to Bun

**Files:**
- Modify: `web/package.json` (add scripts)
- Delete: `web/package-lock.json`
- Create: `web/bun.lockb` (generated by `bun install`)
- Modify: `Makefile` (npm → bun)
- Modify: `deployments/docker/Dockerfile.web` (npm → bun)

**Step 1: Install bun globally (if not present)**

Run: `which bun || curl -fsSL https://bun.sh/install | bash`

**Step 2: Remove npm lockfile and install with bun**

```bash
cd gofin-full/web
rm -f package-lock.json
bun install
```

Expected: `bun.lockb` created, no errors.

**Step 3: Add missing scripts to package.json**

Add to `web/package.json` scripts:
```json
"lint": "svelte-kit sync && svelte-check --tsconfig ./tsconfig.json",
"test": "echo 'no unit tests yet' && exit 0",
"test:e2e": "playwright test"
```

**Step 4: Update root Makefile**

Replace all `npm` references with `bun`:
- `web-install`: `cd web && bun install`
- `web-dev`: `cd web && bun run dev`
- `web-build`: `cd web && bun run build`
- `web-lint`: `cd web && bun run lint`

**Step 5: Update Dockerfile.web**

Replace `node:20-alpine` with `oven/bun:1-alpine` in both stages. Replace npm commands with bun commands:

```dockerfile
# Build stage
FROM oven/bun:1-alpine AS builder

WORKDIR /app

COPY package.json bun.lockb ./
RUN bun install --frozen-lockfile

COPY . .
RUN bun run build

# Runtime stage
FROM oven/bun:1-alpine AS runner

WORKDIR /app

ENV NODE_ENV=production
ENV PORT=3000
ENV HOST=0.0.0.0

COPY --from=builder /app/build ./build
COPY --from=builder /app/package.json ./

EXPOSE 3000

CMD ["bun", "run", "build/index.js"]
```

Note: The runtime CMD may need adjustment based on SvelteKit adapter output format. Verify with `bun run build` output first.

**Step 6: Verify web builds with bun**

Run: `cd gofin-full/web && bun run build 2>&1`

Expected: Build succeeds, output in `build/` or `.svelte-kit/` directory.

**Step 7: Commit**

```bash
git add web/ Makefile deployments/docker/Dockerfile.web
git commit -m "chore: migrate web from npm to bun"
```

---

## Task 7: Fix Web Build Issues

**Files:**
- Modify: `web/svelte.config.js` (adapter for Docker)
- Modify: `web/package.json` (if deps missing)
- Read: `web/src/lib/services/client.ts` (verify API_BASE)

**Step 1: Switch adapter from auto to node**

For Docker compatibility, change `web/svelte.config.js`:
```javascript
import adapter from '@sveltejs/adapter-node';
// ... rest stays same
```

Install adapter: `cd web && bun add -d @sveltejs/adapter-node`

**Step 2: Add /health proxy to vite config**

Add to `web/vite.config.ts` server.proxy:
```typescript
'/health': {
  target: 'http://localhost:8080',
  changeOrigin: true
}
```

**Step 3: Verify build**

Run: `cd gofin-full/web && bun run build 2>&1`

Expected: Successful build, `build/` directory created with `index.js` entry point.

**Step 4: Update Dockerfile.web runtime CMD**

After verifying the build output structure, set the correct CMD:
```dockerfile
CMD ["node", "build/index.js"]
```

The adapter-node output uses Node.js runtime even with bun install, so keep `node:20-alpine` for the runtime stage or use `oven/bun:1-alpine` (bun can run Node.js files).

**Step 5: Commit**

```bash
git add web/ deployments/docker/Dockerfile.web
git commit -m "fix: configure web for Docker deployment with adapter-node"
```

---

## Task 8: Create Docker Dev Environment (API + Web)

**Files:**
- Create: `deployments/docker/docker-compose.dev.yml`
- Modify: `Makefile` (add docker-dev target)
- Modify: `web/vite.config.ts` (proxy for Docker network)

**Step 1: Create docker-compose.dev.yml**

This compose file starts PostgreSQL, Redis, and the Go API. The web app runs on the host (via `bun run dev`) for hot-reload DX, proxying API calls to the Docker network.

```yaml
# Development: Postgres + Redis + API in Docker, Web on host.
# Usage:
#   docker compose -f deployments/docker/docker-compose.dev.yml up -d
#   cd web && bun run dev
services:
  postgres:
    image: postgres:16-alpine
    ports:
      - "5432:5432"
    environment:
      POSTGRES_DB: gofin
      POSTGRES_USER: gofin
      POSTGRES_PASSWORD: gofin_secret
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U gofin -d gofin"]
      interval: 5s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 5s
      retries: 5

  api:
    build:
      context: ../../api
      dockerfile: ../deployments/docker/Dockerfile.api
    ports:
      - "8080:8080"
    environment:
      APP_ENV: local
      APP_DEBUG: "true"
      APP_URL: http://localhost:8080
      TZ: UTC
      DB_HOST: postgres
      DB_PORT: 5432
      DB_DATABASE: gofin
      DB_USERNAME: gofin
      DB_PASSWORD: gofin_secret
      DB_SSL_MODE: disable
      DB_DSN: postgres://gofin:gofin_secret@postgres:5432/gofin?sslmode=disable
      REDIS_HOST: redis
      REDIS_PORT: 6379
      AUTH_PROVIDER: disabled
      AUTH_JWT_SECRET: dev-jwt-secret-not-for-production-32chars!!
      AUTH_ALLOW_REGISTRATION: "true"
      RATE_LIMIT_MAX: 0
      CORS_ALLOWED_ORIGINS: http://localhost:5173
      LOG_LEVEL: debug
      LOG_FORMAT: console
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    volumes:
      - ../../api:/app  # Live reload via air or manual restart
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "-O", "/dev/null", "http://localhost:8080/health"]
      interval: 10s
      timeout: 5s
      retries: 3
      start_period: 15s

volumes:
  postgres_data:
  redis_data:
```

**Step 2: Add Makefile targets**

Add to root Makefile:
```makefile
docker-dev:
	docker compose -f $(COMPOSE_DIR)/docker-compose.dev.yml up -d

docker-dev-down:
	docker compose -f $(COMPOSE_DIR)/docker-compose.dev.yml down

docker-dev-logs:
	docker compose -f $(COMPOSE_DIR)/docker-compose.dev.yml logs -f api
```

**Step 3: Update docker-compose.yml (existing dev with Keycloak)**

Rename `docker-compose.yml` references in the existing Makefile `docker-up`/`docker-down` targets. Keep the existing Keycloak compose as `docker-compose.keycloak.yml` for OAuth development.

Actually — keep `docker-compose.yml` as-is (it's the full dev stack with Keycloak). The new `docker-compose.dev.yml` is the lightweight daily-driver without Keycloak.

**Step 4: Verify dev stack starts**

Run: `cd gofin-full && docker compose -f deployments/docker/docker-compose.dev.yml up -d --build 2>&1`

Expected: postgres, redis, api all healthy.

Run: `curl http://localhost:8080/health`
Expected: `{"status":"ok"}` or similar health response.

Run: `curl http://localhost:8080/api/v1/`
Expected: API info JSON.

**Step 5: Commit**

```bash
git add deployments/docker/docker-compose.dev.yml Makefile
git commit -m "feat: add docker-compose.dev.yml for API+DB+Redis dev stack"
```

---

## Task 9: Verify Web ↔ API Integration

**Files:**
- Read: `web/src/lib/services/client.ts`
- Read: `web/src/lib/services/auth.ts`
- Read: `web/src/lib/stores/auth.svelte.ts`
- Read: `web/src/routes/+page.svelte`

**Step 1: Start dev stack**

Run: `cd gofin-full && docker compose -f deployments/docker/docker-compose.dev.yml up -d`

**Step 2: Start web dev server**

Run: `cd gofin-full/web && bun run dev &`

Wait for: `Local: http://localhost:5173`

**Step 3: Test API proxy**

Run: `curl http://localhost:5173/api/v1/`
Expected: Same response as direct API call (proxy working).

Run: `curl http://localhost:5173/health`
Expected: Health response (if proxy configured).

**Step 4: Test auth flow**

1. Open browser to `http://localhost:5173`
2. Navigate to login page
3. Try registering a new user (AUTH_ALLOW_REGISTRATION=true in dev compose)
4. Verify login works
5. Verify dashboard loads after login

If auth fails, check:
- `web/src/lib/services/client.ts` — token handling
- `web/src/lib/stores/auth.svelte.ts` — auth state
- API CORS settings — `CORS_ALLOWED_ORIGINS=http://localhost:5173`
- API auth provider — `AUTH_PROVIDER=disabled` means no auth required

**Step 5: Document integration status**

Note any issues found. Common problems:
- Backend returns JSON:API format (`{type, id, attributes}`) but frontend expects flat objects
- Frontend uses mock data instead of real API calls
- Auth flow expects different response format

**Step 6: Commit (if fixes were needed)**

```bash
git add web/ deployments/docker/docker-compose.dev.yml
git commit -m "fix: resolve web-API integration issues"
```

**Step 7: Teardown**

Run: `cd gofin-full && docker compose -f deployments/docker/docker-compose.dev.yml down`

---

## Task 10: Install and Configure Playwright

**Files:**
- Modify: `web/package.json` (add @playwright/test)
- Create: `web/playwright.config.ts`
- Create: `web/tests/e2e/` directory structure
- Create: `web/tests/e2e/health.spec.ts` (basic smoke test)

**Step 1: Install Playwright**

```bash
cd gofin-full/web
bun add -d @playwright/test
bunx playwright install --with-deps chromium
```

**Step 2: Create playwright.config.ts**

```typescript
import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './tests/e2e',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: 'html',
  use: {
    baseURL: 'http://localhost:5173',
    trace: 'on-first-retry',
  },
  projects: [
    { name: 'chromium', use: { ...devices['Desktop Chrome'] } },
  ],
  webServer: {
    command: 'bun run dev',
    url: 'http://localhost:5173',
    reuseExistingServer: !process.env.CI,
  },
});
```

**Step 3: Create basic E2E test**

Create `web/tests/e2e/health.spec.ts`:
```typescript
import { test, expect } from '@playwright/test';

test('API health endpoint responds', async ({ request }) => {
  const response = await request.get('/health');
  expect(response.ok()).toBeTruthy();
});

test('home page loads', async ({ page }) => {
  await page.goto('/');
  await expect(page).toHaveTitle(/gofin/i);
});
```

**Step 4: Run Playwright (with dev stack running)**

Start dev stack first: `docker compose -f deployments/docker/docker-compose.dev.yml up -d`

Run: `cd gofin-full/web && bunx playwright test 2>&1`

Expected: Tests pass (health endpoint accessible, home page loads).

**Step 5: Add Playwright commands to Makefile**

```makefile
web-test-e2e:
	cd web && bunx playwright test

web-test-e2e-ui:
	cd web && bunx playwright test --ui
```

**Step 6: Commit**

```bash
git add web/ Makefile
git commit -m "feat: add Playwright E2E testing for web app"
```

---

## Task 11: Write E2E Tests for Auth Flow

**Files:**
- Create: `web/tests/e2e/auth.spec.ts`
- Read: `web/src/lib/services/auth.ts`
- Read: `web/src/routes/(auth)/login/+page.svelte`

**Step 1: Explore auth pages**

Read the login and register page components to understand form field names and selectors.

**Step 2: Write auth E2E tests**

Create `web/tests/e2e/auth.spec.ts`:
```typescript
import { test, expect } from '@playwright/test';

test.describe('Authentication', () => {
  test('login page renders', async ({ page }) => {
    await page.goto('/login');
    await expect(page.getByRole('heading', { name: /login/i })).toBeVisible();
  });

  test('register page renders', async ({ page }) => {
    await page.goto('/register');
    await expect(page.getByRole('heading', { name: /register|sign up/i })).toBeVisible();
  });

  test('redirects to dashboard after login', async ({ page }) => {
    // This test depends on AUTH_PROVIDER and registration being enabled
    await page.goto('/login');
    // Fill in credentials based on actual form fields
    // Submit and verify redirect
  });
});
```

Adjust selectors based on actual component markup after reading the auth pages.

**Step 3: Run tests**

Run: `cd gofin-full/web && bunx playwright test tests/e2e/auth.spec.ts 2>&1`

**Step 4: Commit**

```bash
git add web/tests/e2e/
git commit -m "feat: add E2E tests for authentication flow"
```

---

## Task 12: Write E2E Tests for Core Pages

**Files:**
- Create: `web/tests/e2e/dashboard.spec.ts`
- Create: `web/tests/e2e/wallets.spec.ts`
- Create: `web/tests/e2e/transactions.spec.ts`
- Read: `web/src/routes/(app)/dashboard/+page.svelte`
- Read: `web/src/routes/(app)/wallets/+page.svelte`
- Read: `web/src/routes/(app)/transactions/+page.svelte`

**Step 1: Explore page components**

Read the dashboard, wallets, and transactions pages to understand the UI structure and what to test.

**Step 2: Write smoke tests for each page**

These are smoke tests — verify the page loads and key elements are visible. They do NOT test full business logic (that's integration test territory).

Example `web/tests/e2e/dashboard.spec.ts`:
```typescript
import { test, expect } from '@playwright/test';

test.describe('Dashboard', () => {
  test.beforeEach(async ({ page }) => {
    // Login first if auth is required
    await page.goto('/login');
    // ... login steps
  });

  test('dashboard page loads with summary cards', async ({ page }) => {
    await page.goto('/dashboard');
    // Verify key elements based on actual markup
  });
});
```

**Step 3: Run all E2E tests**

Run: `cd gofin-full/web && bunx playwright test 2>&1`

**Step 4: Commit**

```bash
git add web/tests/e2e/
git commit -m "feat: add E2E smoke tests for dashboard, wallets, transactions"
```

---

## Task 13: Update Dockerfile.web for Playwright

**Files:**
- Modify: `deployments/docker/Dockerfile.web`
- Create: `deployments/docker/Dockerfile.web.test`

**Step 1: Create Dockerfile.web.test for E2E tests**

```dockerfile
FROM oven/bun:1-alpine AS builder

WORKDIR /app

COPY package.json bun.lockb ./
RUN bun install --frozen-lockfile

COPY . .
RUN bunx playwright install --with-deps chromium

# Run tests
CMD ["bunx", "playwright", "test"]
```

**Step 2: Verify Playwright works in Docker (optional)**

This is lower priority — Playwright in Docker requires additional setup (display, dependencies). Mark as future enhancement if it doesn't work immediately.

**Step 3: Commit**

```bash
git add deployments/docker/Dockerfile.web.test
git commit -m "feat: add Dockerfile.web.test for Playwright E2E tests"
```

---

## Task 14: Update Documentation

**Files:**
- Modify: `README.md`
- Modify: `docs/RUNBOOK.md`
- Modify: `docs/ui-pages.md`
- Modify: `docs/tests/TEST_REPORT.md`

**Step 1: Update README.md**

Update to reflect:
- Bun as package manager (not npm)
- Docker dev workflow (`make docker-dev` + `make web-dev`)
- Playwright E2E tests (`make web-test-e2e`)
- Updated directory structure with `docker-compose.dev.yml`

**Step 2: Update RUNBOOK.md**

Add sections for:
- Local development with Docker + Bun
- Running tests (unit, integration, E2E)
- Docker commands reference
- Troubleshooting common issues

**Step 3: Update TEST_REPORT.md**

Add results from:
- API unit test run
- API integration test run (Docker)
- Web build status
- Playwright E2E test results

**Step 4: Commit**

```bash
git add README.md docs/
git commit -m "docs: update documentation for monorepo dev workflow"
```

---

## Task 15: Final Verification & Cleanup

**Files:**
- All files (verification pass)

**Step 1: Full API test run**

Run: `cd gofin-full/api && go test ./tests/unit/... -count=1 2>&1`
Expected: All PASS

**Step 2: Full Docker test run**

Run: `cd gofin-full && docker compose -f deployments/docker/docker-compose.test.yml up --build --abort-on-container-exit 2>&1`
Expected: All PASS

**Step 3: Web build verification**

Run: `cd gofin-full/web && bun run build 2>&1`
Expected: Successful build

**Step 4: Web type check**

Run: `cd gofin-full/web && bun run check 2>&1`
Expected: No type errors

**Step 5: E2E test run (with dev stack)**

Run:
```bash
cd gofin-full && docker compose -f deployments/docker/docker-compose.dev.yml up -d
cd web && bunx playwright test
cd .. && docker compose -f deployments/docker/docker-compose.dev.yml down
```
Expected: All E2E tests pass

**Step 6: Docker compose config validation**

Run: `cd gofin-full && docker compose -f deployments/docker/docker-compose.dev.yml config > /dev/null && echo "OK"`
Run: `cd gofin-full && docker compose -f deployments/docker/docker-compose.selfhost.yml config > /dev/null && echo "OK"`
Expected: Both parse without errors

**Step 7: Grep for stale references**

Run: `grep -rn "firefly" gofin-full/ --include="*.yml" --include="*.yaml" --include="*.env*" --include="*.md" --include="Dockerfile*" | grep -v "docs/research/" | grep -v "docs/refactor/" | grep -v ".git/"`
Expected: 0 results (or only in research/refactor docs which are historical)

**Step 8: Commit any remaining fixes**

```bash
git add -A
git commit -m "chore: final cleanup and verification"
```

---

## Execution Notes

### Subagent Strategy

| Task | Subagent Type | Parallelizable |
|------|--------------|----------------|
| Task 1: CLAUDE.md | Direct (quick) | No — foundation |
| Task 2: Fix .env.example | Direct (quick) | Yes, with Task 1 |
| Task 3: Unit tests | Direct | Yes, after Task 1-2 |
| Task 4: Integration tests | Direct | No — needs Docker |
| Task 5: Docker test suite | Direct | No — needs Task 4 |
| Task 6: Migrate to bun | Direct | Yes, with Task 3-5 |
| Task 7: Fix web build | Direct | No — needs Task 6 |
| Task 8: Docker dev env | Direct | No — needs Task 7 |
| Task 9: Web↔API integration | browser agent | No — needs Task 8 |
| Task 10: Playwright setup | Direct | Yes, with Task 8-9 |
| Task 11: Auth E2E tests | Direct | No — needs Task 10 |
| Task 12: Core page E2E | Direct | No — needs Task 11 |
| Task 13: Dockerfile.web.test | Direct | Yes, with Task 12 |
| Task 14: Documentation | Direct | No — needs all above |
| Task 15: Final verification | Direct | No — needs all above |

### Critical Path

```
Task 1 → Task 2 → Task 3 (unit tests)
                    ↓
               Task 4 (integration tests) → Task 5 (Docker tests)
                                                    ↓
Task 6 (bun migration) → Task 7 (web build) → Task 8 (docker dev) → Task 9 (integration)
                                                                        ↓
                                                              Task 10 (Playwright) → Task 11 → Task 12 → Task 13
                                                                                          ↓
                                                                                    Task 14 (docs) → Task 15 (verify)
```

### Parallelization Opportunities

- **Wave 1** (sequential): Tasks 1-2 (foundation)
- **Wave 2** (parallel): Task 3 (unit tests) + Task 6 (bun migration)
- **Wave 3** (sequential): Task 4-5 (integration tests)
- **Wave 4** (sequential): Task 7-9 (web build, docker dev, integration)
- **Wave 5** (sequential): Task 10-12 (Playwright)
- **Wave 6** (parallel): Task 13 (Dockerfile) + Task 14 (docs)
- **Wave 7** (sequential): Task 15 (final verification)

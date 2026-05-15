# Development

Setup development environment and contribute to Gofin.

## Prerequisites

| Tool | Version | Notes |
|------|---------|-------|
| [Go](https://go.dev/dl/) | 1.25+ | Backend |
| [Node.js](https://nodejs.org/) | 20+ | Frontend tooling |
| [Docker](https://docs.docker.com/get-docker/) | 24+ | PostgreSQL + Redis |
| [Make](https://www.gnu.org/software/make/) | any | Build commands |

## Project Setup

### 1. Clone and install dependencies

```bash
git clone https://github.com/ajianaz/gofin-full.git
cd gofin-full

# Backend dependencies
cd api && go mod download && cd ..

# Frontend dependencies
cd web && npm install && cd ..
```

### 2. Start infrastructure

```bash
make docker-dev
```

This starts PostgreSQL and Redis. The API server also starts and runs migrations.

### 3. Start frontend dev server

```bash
make web-dev
```

Open `http://localhost:5173`. Vite proxies `/api/*` to `http://localhost:8080`.

## Make Commands

| Command | Description |
|---------|-------------|
| `make docker-dev` | Start PostgreSQL + Redis + API (dev mode) |
| `make docker-selfhost` | Start full production stack with Caddy |
| `make web-dev` | Start SvelteKit dev server on port 5173 |
| `make web-build` | Build SvelteKit for production |
| `make migrate` | Run database migrations |
| `make seed` | Seed database with sample data |
| `make test-unit` | Run backend unit tests |
| `make test-integration` | Run backend integration tests |

## Project Structure

```
gofin-full/
├── api/                    # Go backend
│   ├── cmd/server/         # API entry point
│   ├── internal/
│   │   ├── handler/        # HTTP handlers
│   │   ├── service/        # Business logic
│   │   ├── repository/     # Database queries
│   │   ├── middleware/      # Auth, RBAC, rate limit
│   │   ├── domain/         # Domain models
│   │   ├── dto/            # Request/response types
│   │   └── config/         # Configuration
│   ├── migrations/         # SQL migrations
│   └── tests/              # Unit + integration tests
│       ├── unit/
│       └── integration/
│           └── testhelpers/
├── web/                    # SvelteKit 5 frontend
│   ├── src/
│   │   ├── routes/         # File-based routing
│   │   ├── lib/
│   │   │   ├── components/ # UI components
│   │   │   ├── services/   # API client
│   │   │   └── stores/     # State management
│   │   └── app.html
│   └── tests/e2e/          # Playwright tests
├── docs/                   # VitePress documentation
└── deployments/docker/     # Docker Compose files
```

## Backend Development

### Code Structure

The backend follows a layered architecture:

```
Request → Router → Middleware → Handler → Service → Repository → Database
```

- **Handler:** HTTP request/response, input validation via DTOs
- **Service:** Business logic, orchestration between repositories
- **Repository:** Raw SQL queries, data mapping

### Adding a New Endpoint

1. **Define the route** in `internal/router/router.go`
2. **Create DTO** (request/response types) in `internal/dto/`
3. **Add handler** in `internal/handler/`
4. **Add service** logic in `internal/service/`
5. **Add repository** queries in `internal/repository/`
6. **Add migration** in `migrations/` if schema changes
7. **Add tests** in `tests/unit/` and `tests/integration/`

### Running Tests

```bash
# Unit tests (no database required)
cd api
go test -v ./tests/unit/...

# Integration tests (requires PostgreSQL)
DB_HOST=localhost DB_DATABASE=gofin DB_USERNAME=gofin DB_PASSWORD=gofin \
  go test -v ./tests/integration/...
```

### Coding Conventions

- **IDs:** UUID v7 (time-sortable, globally unique)
- **Money:** Always use `decimal.Decimal` — never `float64`
- **SQL:** Parameterized queries only (`$1`, `$2`, ...)
- **Errors:** Use `apperrors` package for typed HTTP errors
- **Context:** Always pass `c.Context()` from Fiber to service/repository layers

## Frontend Development

### Tech Stack

| Technology | Purpose |
|-----------|---------|
| SvelteKit 5 | Full-stack framework (used as static SPA) |
| Svelte 5 | UI components with runes ($state, $derived, $effect) |
| Tailwind CSS 4 | Utility-first styling |
| shadcn-svelte | UI component library (copy-paste, not dependency) |
| Chart.js | Analytics charts |

### Component Library

Two component libraries coexist:

| Path | Usage |
|------|-------|
| `src/lib/components/ui/` | Custom components (brand pages, dark mode) |
| `src/lib/components/ui/shadcn/` | shadcn-svelte components (admin pages) |

### Adding a Page

1. Create `+page.svelte` in the appropriate route directory
2. Create `+page.ts` for server-side data loading if needed
3. Add API service functions in `src/lib/services/`
4. Add navigation link in the sidebar/layout

### Running Frontend Tests

```bash
# Unit tests (Vitest)
cd web
npm run test

# Type checking
npm run check

# E2E tests (requires running API)
npx playwright test

# E2E tests with UI
npx playwright test --ui
```

### i18n

Translations are in `src/lib/i18n/`:

```
src/lib/i18n/
├── id.ts    # Bahasa Indonesia
└── en.ts    # English
```

Each file exports an object with key-value pairs. The active locale is stored in a Svelte store.

## Git Workflow

1. Create a branch from `main`: `git checkout -b feat/my-feature`
2. Make changes with atomic commits
3. Run tests locally before pushing
4. Create a Pull Request
5. Ensure CI passes (unit tests, integration tests, type check, build)
6. Merge after review

### Commit Convention

```
<type>: <short description>

Types: feat, fix, docs, refactor, test, chore
```

Examples:
```
feat: add wallet sharing with role-based access
fix: correct double-entry balance calculation
docs: update deployment guide for Caddy v2
test: add integration tests for budget endpoints
```

### Pre-commit Checklist

- [ ] Unit tests pass: `make test-unit`
- [ ] Integration tests pass: `make test-integration`
- [ ] Frontend type check passes: `cd web && npm run check`
- [ ] Frontend build succeeds: `cd web && npm run build`
- [ ] No hardcoded secrets in code
- [ ] Migration files are sequential (no gaps in numbering)

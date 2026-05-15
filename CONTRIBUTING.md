# Contributing to Gofin

Thank you for your interest in contributing! This guide covers the basics.

## Prerequisites

- [Go](https://go.dev/dl/) 1.25+
- [Node.js](https://nodejs.org/) 20+
- [Bun](https://bun.sh/) (frontend package manager)
- [Docker](https://docs.docker.com/get-docker/) + Docker Compose v2+

## Setup

1. Fork the repository
2. Clone your fork: `git clone https://github.com/<your-username>/gofin-full.git`
3. Create a branch from `develop`: `git checkout -b feat/your-feature develop`

## Development

```bash
# Start PostgreSQL + Redis + API (with hot-reload)
make docker-dev

# Install frontend dependencies
make web-install

# Start frontend dev server (port 5173)
make web-dev
```

## Code Style

### Backend (Go)
- Run `make api-lint` before committing
- Run `make api-tidy` after modifying dependencies
- Follow existing patterns: Handler → Service → Repository layers
- All DB queries must use parameterized SQL (`$1, $2, ...`)
- Use `shopspring/decimal` for monetary values — never `float64`

### Frontend (SvelteKit)
- Run `bun run lint` before committing
- Use shadcn-svelte components — check existing pages for patterns
- Use Svelte 5 runes syntax (`$state`, `$derived`, `$effect`, `$props`)
- Number props use `{N}` syntax, not `"N"` (e.g., `colspan={5}`)

## Testing

```bash
# Backend
make api-test-unit          # Unit tests
make api-test-integration   # Integration tests (needs running DB + Redis)

# Frontend
cd web && npx vitest run    # Unit tests
cd web && npx svelte-check  # Type checking

# Full stack in Docker
make docker-test
```

## Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add recurring transaction support
fix: resolve login redirect loop on expired tokens
docs: update deployment guide for Docker Compose v2
refactor: extract shared pagination component
test: add integration tests for budget CRUD
```

## Pull Requests

1. Branch from `develop` — never from `main`
2. Update `CHANGELOG.md` under `[Unreleased]`
3. Ensure all tests pass
4. Keep PRs focused — one concern per PR
5. Target the `develop` branch

## Project Structure

```
api/
├── cmd/          # Entry points
├── internal/     # Application code
│   ├── handler/  # HTTP handlers (request/response)
│   ├── service/  # Business logic
│   ├── repository/ # Database queries
│   ├── domain/   # Domain models
│   ├── dto/      # Request/response types
│   └── middleware/ # CORS, auth, RBAC, rate limit
├── migrations/   # PostgreSQL migrations
└── tests/        # Unit + integration tests

web/
├── src/
│   ├── routes/   # File-based routing
│   ├── lib/      # Components, services, stores
│   └── i18n/     # Translations (id, en)
└── tests/e2e/    # Playwright tests
```

## Questions?

Open a [Discussion](https://github.com/ajianaz/gofin-full/discussions) for questions, ideas, or feedback.

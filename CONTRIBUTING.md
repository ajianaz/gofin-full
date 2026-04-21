# Contributing to Gofin

Thanks for your interest in contributing! This guide covers the development workflow.

## Prerequisites

- **Go** 1.25+
- **bun** (not npm/yarn — used for the web frontend)
- **Docker** and Docker Compose
- **PostgreSQL** 17 (if running outside Docker)

## Setup

```bash
git clone https://github.com/ajianaz/gofin-full.git
cd gofin-full

# Copy environment config
cp api/.env.example api/.env

# Start development services (API + PostgreSQL + Redis)
make docker-dev

# Install frontend dependencies
make web-install
```

## Development Workflow

### Branch Naming

- `feat/<short-name>` — New features
- `fix/<short-name>` — Bug fixes
- `refactor/<short-name>` — Code restructuring
- `chore/<short-name>` — Maintenance
- `docs/<short-name>` — Documentation

Never commit directly to `main`. Create a branch, validate, then open a PR.

### Commit Style

- `feat: add budget limit warnings`
- `fix: correct wallet balance calculation`
- `docs: update API README`
- `chore: upgrade PostgreSQL to 17`

### Making Changes

**API (Go):**
```bash
cd api
go mod tidy          # Sync dependencies
make test-unit       # Run unit tests
make build           # Verify compilation
```

**Web (SvelteKit):**
```bash
cd web
bun install          # Install dependencies
bun run check        # Type check
bunx playwright test # E2E tests (requires running API)
```

## Running Tests

| Test Type | Command | Docker Required |
|-----------|---------|----------------|
| API unit | `make api-test-unit` | No |
| API integration | `make api-test-integration-infra` then `make api-test-integration` | Yes |
| Full API (Docker) | `make docker-test` | Yes |
| Web type check | `make web-lint` | No |
| E2E (Playwright) | `cd web && bunx playwright test` | Yes |

## Pull Request Process

1. Create a feature branch from `main`
2. Make changes with tests
3. Verify all tests pass locally
4. Push and open a PR
5. Address review feedback
6. Wait for approval before merge

## Code Style

- **Go**: Follow standard `go fmt` formatting. Use `golangci-lint` for linting.
- **TypeScript/Svelte**: Strict mode enabled. Use `bun run check` for type checking.
- **CSS**: Tailwind CSS utility classes only. No custom CSS files unless necessary.

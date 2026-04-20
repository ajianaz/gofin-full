# Gofin Full — Claude Code Guidelines

## Project Overview
Personal finance tracker monorepo. Go API backend + SvelteKit frontend.

## Structure
- `api/` — Go backend (module: github.com/ajianaz/gofin-full/api)
- `web/` — SvelteKit frontend (bun, not npm)
- `deployments/docker/` — All Docker configs
- `docs/` — OpenAPI spec, research, plans
- `mobile/` — Future Flutter app (placeholder)

## Rules
- **No `Co-Authored-By` in commits.** Just write the commit message.
- Use `bun` for web, not npm/node.
- API runs on port 8080, web dev on port 5173.
- Integration tests require Docker (postgres:5433, redis:6380).
- Unit tests run without Docker.
- Module path is `github.com/ajianaz/gofin-full/api` (NOT mis-puragroup, NOT azfirazka).

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

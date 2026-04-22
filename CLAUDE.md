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
- **Always create a new branch for each task/group of related work.** Never commit directly to `main`. Branch naming: `feat/<short-name>`, `fix/<short-name>`, `chore/<short-name>`.
- **Update CHANGELOG.md before every commit.** Add entries under `## [Unreleased]` section. Include Added/Changed/Fixed/Removed subsections as needed. Move to versioned section on release.
- Use `bun` for web, not npm/node.
- API runs on port 8080, web dev on port 5173.
- Integration tests require Docker (postgres:5433, redis:6380).
- Unit tests run without Docker.
- Module path is `github.com/ajianaz/gofin-full/api`.

## Commands
- `make docker-dev` — Start daily dev stack (API + Postgres + Redis)
- `make docker-up` — Start full dev stack (API + Postgres + Redis + Keycloak)
- `make docker-test` — Run full test suite in Docker (Postgres + Redis + Test runner)
- `make docker-selfhost` — Start self-hosted production stack
- `make api-test-unit` — Unit tests only
- `make api-test-integration` — Integration tests (needs Docker)
- `make api-test-integration-infra` — Start integration test infrastructure (Postgres:5433, Redis:6380)
- `make web-dev` — Start SvelteKit dev server
- `make web-lint` — Web type check

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

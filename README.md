# Gofin Full

Personal finance tracker monorepo — Go API backend + SvelteKit frontend.

## Structure

```
gofin-full/
├── api/                # Go backend (Fiber v2)
│   ├── cmd/server/     # API server entrypoint
│   ├── internal/       # auth, config, domain, handler, middleware, repository, router, service
│   ├── pkg/           # shared utilities (crypto, decimal, pagination, etc.)
│   ├── migrations/    # PostgreSQL migrations (goose format)
│   └── tests/         # unit, integration, benchmark
├── web/                # SvelteKit 5 frontend
│   ├── src/lib/        # components, services, stores, types, i18n
│   ├── src/routes/     # (auth)/ login/register, (app)/ dashboard, settings, etc.
│   └── tests/e2e/     # Playwright E2E tests
├── deployments/docker/ # Docker configs
│   ├── docker-compose.yml          # Full dev (API + Postgres + Redis + Keycloak)
│   ├── docker-compose.dev.yml     # Daily dev (API + Postgres + Redis)
│   ├── docker-compose.selfhost.yml # Production (Caddy + API + Web + Postgres + Redis)
│   └── docker-compose.test.yml    # Test runner (unit + integration)
├── docs/               # OpenAPI spec, runbook, research
├── scripts/            # Utility scripts
└── mobile/             # Future Flutter app (placeholder)
```

## Quick Start

### Docker Development (recommended)

```bash
# Start API + database + Redis
make docker-dev

# Start frontend dev server (in a separate terminal)
make web-dev
```

Open http://localhost:5173 — the web app proxies API requests to http://localhost:8080.

### Manual Setup

**Backend:**
```bash
cd api
cp .env.example .env
go run ./cmd/server
```

**Frontend:**
```bash
cd web
bun install
bun run dev
```

## Testing

| Test Type | Command | Requires Docker |
|-----------|---------|----------------|
| API unit tests | `make api-test-unit` | No |
| API integration tests | `make api-test-integration` | Yes (make api-test-integration-infra) |
| Full API test suite | `make docker-test` | Yes |
| Web type check | `make web-lint` | No |
| E2E tests | `cd web && bunx playwright test` | Yes (make docker-dev) |

## Docker

| Compose File | Use Case | Services |
|------------|----------|----------|
| `docker-compose.dev.yml` | Daily development | API + Postgres + Redis |
| `docker-compose.yml` | OAuth development | API + Postgres + Redis + Keycloak |
| `docker-compose.selfhost.yml` | Production deployment | Caddy + API + Web + Postgres + Redis |
| `docker-compose.test.yml` | CI/testing | Postgres + Redis + Test runner |

## API

- Health: `GET /health`
- API info: `GET /api/v1/`
- API docs: `GET /api/v1/docs`
- OpenAPI: `GET /api/v1/openapi.json`

## License

Private repository.

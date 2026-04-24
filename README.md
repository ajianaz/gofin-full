# Gofin

Self-hosted personal finance tracker. Go API + SvelteKit web frontend.

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go)
![SvelteKit](https://img.shields.io/badge/SvelteKit-5-FF3E00?logo=svelte)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17-4169E1?logo=postgresql)
![Redis](https://img.shields.io/badge/Redis-7-DC382D?logo=redis)
![License](https://img.shields.io/badge/License-Apache_2.0-blue)

Gofin is a self-hosted personal finance tracker. It uses **wallets** as the primary financial container, supports double-entry bookkeeping with split transactions, hierarchical role-based access control across groups and wallets, and real-time notifications via Server-Sent Events.

## Features

- **125+ API endpoints** with double-entry bookkeeping and split transactions
- **Hierarchical RBAC** -- 21 group-level roles + 3 wallet-level roles (owner/editor/viewer)
- **Real-time notifications** via Server-Sent Events (SSE)
- **40+ web pages** -- dashboard, wallets, transactions, budgets, piggy banks, analytics, settings
- **i18n** -- Indonesian and English, dark mode, responsive design
- **CSV/OFX export**, rules engine, recurring transactions, bill tracking
- **Auth** -- JWT with refresh tokens, OAuth2 (Google, GitHub), optional Keycloak OIDC
- **Self-hosted** -- single `docker compose` command with automated daily DB backup

## Tech Stack

| Layer       | Technology                                       |
|-------------|--------------------------------------------------|
| Backend     | Go 1.25, Fiber v2, PostgreSQL 17, Redis 7       |
| Frontend    | SvelteKit 5, Svelte 5, Tailwind CSS 4, shadcn-svelte |
| Auth        | JWT (golang-jwt/v5), OAuth2, Keycloak OIDC      |
| Deployment  | Docker Compose, Caddy (HTTPS)                    |

## Quick Start

### Self-Hosted (production)

```bash
git clone https://github.com/ajianaz/gofin-full.git
cd gofin-full
cp .env.example .env
# Edit .env -- set AUTH_JWT_SECRET, DOMAIN, ADMIN_EMAIL
make docker-selfhost
```

Open `https://your-domain` -- Caddy auto-provisions HTTPS.

### Development

```bash
# Terminal 1: API + PostgreSQL + Redis
make docker-dev

# Terminal 2: Frontend dev server
make web-dev
```

Open `http://localhost:5173` -- the web app proxies API requests to `http://localhost:8080`.

## Monorepo Structure

```
gofin-full/
├── api/                # Go backend (Fiber v2, PostgreSQL, Redis)
│   ├── cmd/            # server, migrate, seed binaries
│   ├── internal/       # auth, domain, handler, middleware, repository, service
│   ├── pkg/            # shared utilities
│   ├── migrations/     # PostgreSQL migrations
│   └── tests/          # unit + integration tests
├── web/                # SvelteKit 5 frontend
│   ├── src/routes/     # (auth)/ and (app)/ route groups
│   ├── src/lib/        # components, stores, services, i18n
│   └── tests/e2e/      # Playwright E2E tests
├── deployments/docker/ # Docker Compose configs
├── docs/               # OpenAPI spec, architecture, runbook, research
├── scripts/            # Utility scripts
└── mobile/             # Future Flutter app (placeholder)
```

## Documentation

| Document | Description |
|----------|-------------|
| [API Reference](docs/openapi.yaml) | OpenAPI 3.0 spec (79 paths) |
| [Architecture](docs/ARCHITECTURE.md) | System design and data flow |
| [Runbook](docs/RUNBOOK.md) | Operations, deployment, troubleshooting |
| [API README](api/README.md) | Backend details, RBAC, env vars |
| [Web README](web/README.md) | Frontend setup and conventions |
| [Research](docs/research/README.md) | Design research and decisions |

## Testing

| Test Type | Command | Requires Docker |
|-----------|---------|-----------------|
| API unit tests | `make api-test-unit` | No |
| API integration tests | `make api-test-integration` | Yes (`make api-test-integration-infra`) |
| Full API test suite | `make docker-test` | Yes |
| Web type check | `make web-lint` | No |
| E2E tests | `cd web && bunx playwright test` | Yes (`make docker-dev`) |

## Docker

| Compose File | Use Case | Services |
|--------------|----------|----------|
| `docker-compose.dev.yml` | Daily development | API + PostgreSQL + Redis |
| `docker-compose.yml` | OAuth development | API + PostgreSQL + Redis + Keycloak |
| `docker-compose.selfhost.yml` | Production deployment | Caddy + API + Web + PostgreSQL + Redis + Backup |
| `docker-compose.test.yml` | CI/testing | PostgreSQL + Redis + Test runner |

## License

Apache License, Version 2.0

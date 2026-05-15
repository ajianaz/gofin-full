# Gofin

> Self-hosted personal finance tracker with double-entry bookkeeping, hierarchical RBAC, and automated budgeting.

[![Documentation](https://img.shields.io/badge/docs-VitePress-blue)](https://ajianaz.github.io/gofin-full/)
[![Go](https://img.shields.io/badge/Go-1.25-00ADD8)](https://go.dev/)
[![License: Apache-2.0](https://img.shields.io/badge/License-Apache--2.0-green)](LICENSE)

## Quick Start

```bash
git clone https://github.com/ajianaz/gofin-full.git
cd gofin-full
cp .env.example .env
# Edit .env — set AUTH_JWT_SECRET, DOMAIN, ADMIN_EMAIL
make docker-selfhost
```

Open `https://your-domain` — Caddy handles HTTPS automatically.

## Screenshots

<!-- TODO: Add screenshots after first production deployment -->
<!-- To add: run `make docker-selfhost`, navigate to pages, capture with browser -->
<!-- Place in docs/public/ as dashboard.png, transactions.png, analytics.png -->

> 📸 Screenshots coming soon. Run `make docker-selfhost` to try it locally!

## Features

- 💰 **Double-entry bookkeeping** — source + destination wallet for full auditability
- 🛡️ **21-level RBAC** — group roles (read-only to owner) + wallet roles (viewer/editor/owner)
- 📊 **Analytics** — spending by category/period, net worth tracking, budget analysis
- 🔄 **Automation** — rules engine, recurring transactions, bill tracking
- 🌍 **i18n** — Indonesian + English, dark mode, responsive design
- 🔐 **Security** — JWT, OAuth2 (Google/GitHub), rate limiting, password policy
- 🐳 **Docker** — single `docker compose` command, auto-HTTPS via Caddy
- 📤 **Export** — CSV + OFX format
- 🔔 **Real-time** — Server-Sent Events notifications

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.25, Fiber v2, PostgreSQL 17, Redis 7 |
| Frontend | SvelteKit 5, Tailwind CSS 4, shadcn-svelte |
| Auth | JWT, OAuth2 (Google, GitHub), Keycloak OIDC |
| Deploy | Docker Compose, Caddy (auto-HTTPS) |

## Documentation

📖 **Full documentation** at [ajianaz.github.io/gofin-full](https://ajianaz.github.io/gofin-full/)

- [Getting Started](https://ajianaz.github.io/gofin-full/getting-started) — Deploy in minutes
- [Features](https://ajianaz.github.io/gofin-full/features) — Complete feature overview
- [Architecture](https://ajianaz.github.io/gofin-full/architecture) — System design & data flow
- [Configuration](https://ajianaz.github.io/gofin-full/configuration) — All environment variables
- [Deployment](https://ajianaz.github.io/gofin-full/deployment) — Production deployment guide
- [Security](https://ajianaz.github.io/gofin-full/security) — Security features & hardening
- [RBAC](https://ajianaz.github.io/gofin-full/rbac) — Permission system explained
- [API Reference](https://ajianaz.github.io/gofin-full/api/) — OpenAPI 3.0 specification (135 endpoints)

## Project Structure

```
gofin-full/
├── api/                  # Go backend
│   ├── internal/         # Handlers, services, repositories
│   ├── migrations/       # PostgreSQL migrations
│   └── tests/            # Unit + integration tests
├── web/                  # SvelteKit 5 frontend
│   ├── src/routes/       # File-based routing (43 pages)
│   └── tests/e2e/        # Playwright tests
├── docs/                 # VitePress documentation
├── deployments/docker/   # Docker Compose configs
└── scripts/              # Utility scripts
```

## Development

```bash
# Start infrastructure (PostgreSQL + Redis + API)
make docker-dev

# Start frontend dev server (port 5173)
make web-dev

# Run tests
make test-unit
make test-integration
```

See the [Development Guide](https://ajianaz.github.io/gofin-full/development) for full details.

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

[Apache-2.0](LICENSE)

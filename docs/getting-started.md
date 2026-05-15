# Getting Started

Get Gofin running in minutes. Choose your deployment method.

## Prerequisites

| Requirement | Version |
|------------|---------|
| [Docker](https://docs.docker.com/get-docker/) | 24+ |
| [Docker Compose](https://docs.docker.com/compose/install/) | v2+ |
| A domain (for HTTPS) | Any public domain |

::: tip No domain?
You can use an IP address or `localhost` for testing. Caddy will use self-signed certs automatically.
:::

## Self-Hosted (Production)

The recommended way to run Gofin. Includes Caddy (reverse proxy + HTTPS), API server, web frontend, PostgreSQL, Redis, and automated daily backups.

### 1. Clone the repository

```bash
git clone https://github.com/ajianaz/gofin-full.git
cd gofin-full
```

### 2. Configure environment

```bash
cp .env.example .env
```

Edit `.env` and set at minimum:

```ini
# REQUIRED — change these!
AUTH_JWT_SECRET=your-random-secret-at-least-32-chars
STATIC_CRON_TOKEN=your-random-cron-token
DOMAIN=your-domain.com

# RECOMMENDED — create admin user on first startup
ADMIN_EMAIL=admin@example.com
ADMIN_PASSWORD=your-secure-password
```

::: warning Important
- `AUTH_JWT_SECRET` must be at least 32 characters. Generate one with `openssl rand -hex 32`.
- `STATIC_CRON_TOKEN` is used for internal cron endpoints. Generate one with `openssl rand -hex 16`.
- If `ADMIN_PASSWORD` is empty, a random 16-char password is generated and printed to the logs.
:::

### 3. Start the stack

```bash
make docker-selfhost
```

This starts 6 services:

| Service | Description | Port |
|---------|-------------|------|
| **Caddy** | Reverse proxy + auto-HTTPS | 80, 443 |
| **API** | Go backend (Fiber) | Internal (8080) |
| **Web** | SvelteKit static build | Internal |
| **PostgreSQL** | Database | Internal (5432) |
| **Redis** | Cache + sessions | Internal (6379) |
| **Backup** | Daily DB backup | Cron (03:00 UTC) |

### 4. Open your browser

Navigate to `https://your-domain.com`. Log in with the admin credentials you set.

## Development Mode

For local development with hot-reload.

### Terminal 1 — Start infrastructure

```bash
make docker-dev
```

This starts PostgreSQL + Redis + the API server with auto-migrations.

### Terminal 2 — Start frontend dev server

```bash
make web-dev
```

Open `http://localhost:5173`. The Vite dev server proxies API requests to `http://localhost:8080`.

## First Steps After Login

1. **Create a wallet** — Go to Wallets → Create. Choose a type (bank account, cash, credit card, etc.)
2. **Set your currency** — Go to Settings → Preferences and select your default currency
3. **Add a transaction** — Go to Transactions → Create. Select source and destination wallets
4. **Create a budget** — Go to Budgets → Create. Set a spending limit per period
5. **Invite members** — (Optional) Go to a wallet → Members → Invite. Share access with role-based permissions

## What's Next?

- [Features Overview](/features) — See everything Gofin can do
- [Configuration Reference](/configuration) — All environment variables
- [Deployment Guide](/deployment) — Advanced deployment options
- [API Reference](/api/) — Full OpenAPI specification

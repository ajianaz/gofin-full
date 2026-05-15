# Architecture

System design, data flow, and key decisions behind Gofin.

::: tip Interactive Diagram
View the full interactive architecture diagram: [architecture.html](/architecture.html)
:::

## High-Level Overview

```
┌─────────────────────────────────────────────────────┐
│                    Internet                          │
└──────────┬──────────────────────┬───────────────────┘
           │                      │
     ┌─────▼─────┐          ┌─────▼─────┐
     │   HTTP    │          │   HTTPS   │
     └─────┬─────┘          └─────┬─────┘
           │                      │
     ┌─────▼──────────────────────▼──────┐
     │        Caddy (Reverse Proxy)       │
     │     Auto-HTTPS via Let's Encrypt    │
     └─────┬──────────────────────┬──────┘
           │ /api                  │ /*
     ┌─────▼─────┐          ┌─────▼─────┐
     │ Go API    │          │  SvelteKit │
     │ (Fiber)   │          │   Static   │
     └─────┬─────┘          └───────────┘
           │
     ┌─────┴─────┐
     │           │
┌────▼───┐ ┌────▼───┐
│Postgres│ │ Redis  │
│   17   │ │   7    │
└────────┘ └────────┘
```

## Monorepo Structure

```
gofin-full/
├── api/                  # Go backend
│   ├── cmd/              # Entry points (server, migrate, seed)
│   ├── internal/         # Application code
│   │   ├── auth/         # JWT, RBAC, context helpers
│   │   ├── config/       # Viper-based configuration
│   │   ├── domain/       # Domain models (Wallet, Transaction, Budget...)
│   │   ├── dto/          # Request/response DTOs with validation
│   │   ├── handler/      # HTTP handlers (request/response)
│   │   ├── middleware/    # CORS, auth, RBAC, rate limit, metrics
│   │   ├── repository/   # SQL queries, data access
│   │   ├── service/      # Business logic, orchestration
│   │   └── sse/          # Server-Sent Events hub
│   ├── pkg/              # Shared utilities
│   │   ├── pgxuuid/      # Custom pgx UUID codec
│   │   ├── uuid/         # UUID v7 generator
│   │   └── errors/       # App-level error types
│   ├── migrations/       # PostgreSQL migrations
│   └── tests/            # Unit + integration tests
├── web/                  # SvelteKit 5 frontend
│   ├── src/
│   │   ├── routes/       # File-based routing
│   │   │   ├── (auth)/   # Public routes (login, register, etc.)
│   │   │   └── (app)/    # Protected routes (dashboard, wallets, etc.)
│   │   ├── lib/
│   │   │   ├── components/  # UI components (shadcn-svelte + custom)
│   │   │   ├── services/    # API client functions
│   │   │   ├── stores/      # Svelte 5 state management
│   │   │   └── i18n/        # Translations (id, en)
│   │   └── app.html      # HTML shell
│   └── tests/e2e/        # Playwright E2E tests
├── deployments/docker/   # Docker Compose configs
├── docs/                 # Documentation (this site)
└── scripts/              # Utility scripts
```

## API Architecture

The Go backend follows a clean layered architecture:

```
HTTP Request
    │
    ▼
┌─────────────────┐
│  Fiber Router    │  Route registration (internal/router/)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Middleware      │  CORS → Auth → Group Role → Rate Limit → RBAC → Metrics
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Handler         │  HTTP request/response, input validation (internal/handler/)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Service         │  Business logic, orchestration (internal/service/)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Repository      │  SQL queries, data access (internal/repository/)
└────────┬────────┘
         │
    ┌────┴────┐
    ▼         ▼
PostgreSQL   Redis
```

### Middleware Chain

Each request passes through the middleware chain in order:

1. **CORS** — Validates origin against allowed list
2. **Request Logger** — Logs method, path, status, duration
3. **Auth** — Validates JWT, extracts `user_id` and `group_id`
4. **Group Role** — Looks up user's role in their active group
5. **Rate Limit** — Sliding window rate limiting (configurable)
6. **RBAC** — Checks group-level permission for the endpoint
7. **Wallet RBAC** — (Endpoint-specific) Checks wallet membership role

## Web Architecture

SvelteKit 5 with file-based routing and Svelte 5 runes.

### Route Groups

```
src/routes/
├── (auth)/                    # Public — no auth required
│   ├── login/
│   ├── register/
│   ├── forgot-password/
│   └── reset-password/
│
├── (app)/                     # Protected — auth required
│   ├── dashboard/
│   ├── wallets/
│   ├── transactions/
│   ├── budgets/
│   ├── piggy-banks/
│   ├── recurring/
│   ├── bills/
│   ├── rules/
│   ├── analytics/
│   ├── reports/
│   ├── categories/
│   ├── tags/
│   ├── currencies/
│   ├── groups/
│   ├── export/
│   ├── settings/
│   └── admin/
│
└── +page.svelte               # Root redirect
```

### State Management

Svelte 5 runes-based stores:

| Store | Purpose |
|-------|---------|
| `auth` | JWT token, user profile, login/logout actions |
| `theme` | Dark/light mode, system preference detection |
| `i18n` | Current locale, translation function |

### API Communication

All API calls go through a typed service layer in `src/lib/services/`:

- Automatic JWT header injection from auth store
- Token refresh on 401 responses
- Typed request/response via TypeScript interfaces
- Error handling with user-friendly messages

## Data Flow

```
Browser                    Vite Proxy              API Server
  │                           │                        │
  │── GET /api/v1/wallets ───│── GET /api/v1/wallets──│
  │   Authorization: Bearer   │                        │
  │   <token>                 │                   ┌────┴────┐
  │                           │                   │ Handler │
  │                           │                   │  →      │
  │                           │                   │ Service │
  │                           │                   │  →      │
  │                           │                   │  Repo   │
  │                           │                   │  →      │
  │◄── 200 JSON ─────────────│◄── 200 JSON ──────│  PG     │
  │                           │                   └─────────┘
```

## Authentication Flow

```
1. User submits login form
2. POST /api/v1/auth/login → { access_token, refresh_token }
3. Auth store saves access_token to localStorage
4. Every API request includes: Authorization: Bearer <token>
5. Middleware validates JWT, extracts user_id + group_id into context
6. RBAC middleware checks permissions for the requested resource
7. On token expiry: POST /api/v1/auth/refresh → new access_token
```

## Key Design Decisions

| Decision | Rationale |
|----------|-----------|
| **UUID v7** for all IDs | Globally unique, time-sortable, no sequence coordination |
| **Double-entry bookkeeping** | Source + destination wallet for full auditability |
| **Hierarchical RBAC (21 levels)** | Fine-grained group access from read-only to full owner |
| **Server-Sent Events** | Simpler than WebSocket for notifications (one-directional) |
| **Svelte 5 runes** | Fine-grained reactivity without wrapper objects |
| **shadcn-svelte** | Copy-paste components — full control, no heavy dependencies |
| **Custom i18n store** | Lightweight, no framework overhead, only 2 locales |
| **Caddy** | Automatic HTTPS via Let's Encrypt, simple Caddyfile config |
| **shopspring/decimal** | Exact decimal arithmetic for financial amounts (no float) |
| **Docker Compose** | Single-command deployment with all dependencies |

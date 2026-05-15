# Architecture

## Monorepo Layout

```
gofin-full/
├── api/          Go backend (Fiber v2)
├── web/          SvelteKit 5 frontend
├── deployments/  Docker configs
├── docs/         Documentation
├── scripts/      Utility scripts
└── mobile/       Future Flutter app (placeholder)
```

## API Architecture

The Go backend follows a layered architecture:

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
│  Middleware      │  CORS → Auth → Rate Limit → RBAC → Metrics
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Handler         │  HTTP request/response, validation (internal/handler/)
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

### Key Packages

| Package | Purpose |
|---------|---------|
| `internal/auth/` | JWT generation/validation, RBAC middleware, context helpers |
| `internal/config/` | Viper-based env configuration |
| `internal/domain/` | Domain models (Wallet, Transaction, Budget, etc.) |
| `internal/dto/` | Request/response DTOs with validation tags |
| `internal/middleware/` | CORS, auth, RBAC, metrics, caching middleware |
| `internal/sse/` | Server-Sent Events hub for real-time notifications |
| `pkg/pgxuuid/` | Custom pgx codec for UUID ↔ google/uuid |
| `pkg/uuid/` | UUID v7 generator |

## Web Architecture

SvelteKit 5 with file-based routing:

```
src/routes/
├── (auth)/                    # Public routes (no auth required)
│   ├── +layout.svelte         # Auth layout (centered card)
│   ├── login/+page.svelte
│   ├── register/+page.svelte
│   ├── forgot-password/+page.svelte
│   └── reset-password/+page.svelte
│
├── (app)/                     # Protected routes (auth required)
│   ├── +layout.svelte         # App shell (sidebar + topbar)
│   ├── dashboard/+page.svelte
│   ├── wallets/+page.svelte
│   ├── transactions/+page.svelte
│   ├── budgets/+page.svelte
│   ├── piggy-banks/+page.svelte
│   ├── recurring/+page.svelte
│   ├── bills/+page.svelte
│   ├── rules/+page.svelte
│   ├── analytics/+page.svelte
│   ├── settings/+page.svelte
│   └── ...
│
├── +layout.svelte             # Root layout (providers)
└── +page.svelte               # Redirect to /dashboard or /login
```

### State Management

Svelte 5 runes-based stores in `src/lib/stores/`:
- **auth store** — JWT token, user profile, login/logout actions
- **theme store** — Dark/light mode, system preference detection
- **i18n store** — Current locale, translation function

### API Communication

All API calls go through Vite's dev proxy:

```typescript
// vite.config.ts
server: {
  proxy: {
    '/api': 'http://localhost:8080'
  }
}
```

API client functions in `src/lib/services/` use `fetch()` with:
- Automatic JWT header injection from auth store
- Token refresh on 401 responses
- Typed request/response via TypeScript interfaces

## Data Flow

```
Browser                    Vite Proxy              API Server
  │                           │                        │
  │── GET /api/v1/wallets ───│── GET /api/v1/wallets──│
  │   Authorization: Bearer    │                        │
  │   <token>                  │                   ┌────┴────┐
  │                           │                   │ Handler │
  │                           │                   │  →      │
  │                           │                   │ Service │
  │                           │                   │  →      │
  │                           │                   │  Repo   │
  │                           │                   │  →      │
  │◄── 200 JSON ─────────────│◄── 200 JSON ──────│  PG     │
  │                           │                   └─────────┘
```

## Auth Flow

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
| UUID v7 for all IDs | Globally unique, time-sortable, no sequence coordination needed |
| Double-entry bookkeeping | Every transaction has source + destination wallet for auditability |
| Hierarchical RBAC (21 levels) | Fine-grained group access from read-only to full owner |
| Server-Sent Events | Simpler than WebSocket for notifications (one-directional) |
| Svelte 5 runes | Fine-grained reactivity without wrapper objects |
| shadcn-svelte | Copy-paste components — full control over code, no heavy dependencies |
| Custom i18n store | Lightweight, no framework overhead, only 2 locales needed |
| Caddy for self-host | Automatic HTTPS via Let's Encrypt, simple Caddyfile config |

## Deployment

### Self-Hosted (Docker Compose)

```
Internet → Caddy (443/80)
              │
              ├── / → Web (static SvelteKit build)
              └── /api → App (Go Fiber)
                          │
                     ┌────┴────┐
                     │ Postgres│
                     │ Redis   │
                     └─────────┘
```

See `deployments/docker/docker-compose.selfhost.yml` and [Runbook](RUNBOOK.md) for details.

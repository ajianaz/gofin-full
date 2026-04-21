# Gofin API

A self-hosted personal finance manager built with Go -- inspired by [Firefly III](https://www.firefly-iii.org/).

Gofin is designed for performance, multi-tenancy, and fine-grained access control. It uses **wallets** (not accounts) as the primary financial container, with double-entry bookkeeping, hierarchical RBAC, and real-time notifications.

Module: `github.com/ajianaz/gofin-full/api`

## Features

- **125+ API endpoints** covering the full personal finance domain
- Double-entry bookkeeping with split transactions
- Hierarchical RBAC with 21 group-level roles and 3 wallet-level roles (owner/editor/viewer)
- Real-time notifications via Server-Sent Events (SSE)
- Multi-group and multi-wallet support with wallet sharing between users
- Budgets, piggy banks, recurring transactions, rules engine, and bill tracking
- Currency management with exchange rates
- CSV/OFX export and reconciliation
- Analytics: spending by category, spending by period, net worth
- Audit trail and webhook support
- Prometheus metrics endpoint
- JWT authentication with refresh tokens
- Optional Keycloak (OIDC) integration
- OAuth2 login (Google, GitHub)
- API key support for long-lived integrations
- Feature flags (export, webhooks, debts, expression engine, running balance)

> All primary and foreign keys use UUID v7 for globally unique, time-sortable identifiers.

## Tech Stack

| Component     | Technology                          |
|---------------|-------------------------------------|
| Language      | Go 1.25                             |
| HTTP Framework| Fiber v2                            |
| Database      | PostgreSQL 17                       |
| Cache         | Redis 7                             |
| Auth          | JWT (golang-jwt/v5), Keycloak OIDC  |
| Migrations    | Goose                               |
| Logging       | Zerolog (JSON)                      |
| Config        | Viper (env-based)                   |
| Metrics       | Prometheus client_golang            |
| Observability | OpenTelemetry-ready middleware      |

## Local Development

Prerequisites: Go 1.25+, PostgreSQL 17, Redis 7

```bash
# Start API + PostgreSQL + Redis via the monorepo Makefile
cd .. && make docker-dev

# Or run manually:
cd api
cp .env.example .env
go mod tidy

# Run migrations (requires goose)
make migrate-up DB_DSN="postgres://gofin:gofin_secret@localhost:5432/gofin?sslmode=disable"

# Build and run
make run
```

Or run directly:

```bash
go run ./cmd/server
```

## API Documentation

Interactive docs and the OpenAPI spec are served at runtime:

- `GET /api/v1/docs` -- HTML documentation UI
- `GET /api/v1/openapi.json` -- OpenAPI 3.0 specification

All protected endpoints require a JWT token in the `Authorization: Bearer <token>` header, or an API key via `Authorization: Bearer gofin_...` or `X-API-Key: gofin_...`.

### Auth Endpoints (public)

```
GET    /api/v1/auth/provider    # Get configured auth provider
POST   /api/v1/auth/login       # Email/password login
POST   /api/v1/auth/logout      # Invalidate token
POST   /api/v1/auth/refresh     # Refresh JWT
```

### Core Resources

```
/api/v1/wallets           # Wallet CRUD (financial accounts)
/api/v1/wallet-types      # Wallet type reference data
/api/v1/transactions      # Transaction CRUD with double-entry
/api/v1/categories        # Category management
/api/v1/tags              # Tag management
/api/v1/budgets           # Budget tracking
/api/v1/bills             # Bill management
/api/v1/currencies        # Currency reference data
/api/v1/piggy_banks       # Savings goals (nested under wallets)
```

## Project Structure

```
api/
├── cmd/
│   ├── server/main.go          # Application entry point
│   ├── migrate/main.go         # Migration runner
│   └── seed/main.go            # Seed admin user
├── internal/
│   ├── auth/                   # JWT, RBAC middleware, role definitions
│   ├── config/                 # Viper-based configuration
│   ├── domain/                 # Domain models (wallet, transaction, etc.)
│   ├── dto/                    # Request/response DTOs
│   ├── handler/                # HTTP handlers (one per resource)
│   ├── middleware/             # CORS, auth, RBAC, metrics, caching
│   ├── repository/             # Database access (pgx)
│   ├── router/                 # Route registration
│   ├── service/                # Business logic
│   └── sse/                    # Server-Sent Events hub
├── pkg/
│   ├── bcrypt/                 # Password hashing
│   ├── crypto/                 # Encryption utilities
│   ├── currency/               # Currency formatting
│   ├── decimal/                # Precise decimal arithmetic
│   ├── errors/                 # Shared error types
│   ├── hash/                   # Hashing utilities
│   ├── hmac/                   # HMAC signatures
│   ├── pagination/             # Cursor/offset pagination
│   ├── pgxuuid/                # UUID type for pgx driver
│   └── uuid/                   # UUID generation (v7)
├── deployments/
│   └── docker/
│       ├── Dockerfile
│       ├── docker-compose.yml
│       ├── docker-compose.test.yml
│       └── entrypoint.sh
├── migrations/
│   └── postgres/               # Goose SQL migrations
├── tests/
│   ├── unit/                   # Unit tests
│   └── integration/            # Integration tests
├── docs/                       # Planning and research docs
├── entrypoint.sh               # Container entry point script
├── .env.example                # Environment variable template
├── Makefile                    # Build, test, and migration targets
└── go.mod
```

## Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Unit tests only
make test-unit

# Integration tests (starts isolated PostgreSQL/Redis on ports 5433/6380)
make test-integration-infra   # start containers
make test-integration         # run tests
make test-integration-teardown # stop and remove containers
```

## RBAC System

### Group Roles (21 levels)

Group roles follow a strict hierarchy. A higher role implicitly grants all permissions of lower roles.

| Level | Role                    | Scope                         |
|-------|-------------------------|-------------------------------|
| 1     | `read_only`             | View all data                 |
| 2     | `manage_transactions`   | Create/edit transactions      |
| 3     | `manage_meta`           | Manage categories, tags, etc. |
| 4-5   | `read_budgets` / `manage_budgets` | Budget CRUD          |
| 6-7   | `read_piggy_banks` / `manage_piggy_banks` | Piggy bank CRUD |
| 8-9   | `read_subscriptions` / `manage_subscriptions` | Subscription CRUD |
| 10-11 | `read_rules` / `manage_rules` | Rule engine CRUD     |
| 12-13 | `read_recurring` / `manage_recurring` | Recurring transactions |
| 14-15 | `read_webhooks` / `manage_webhooks` | Webhook CRUD   |
| 16-17 | `read_currencies` / `manage_currencies` | Currency management |
| 18    | `view_reports`          | Access analytics endpoints    |
| 19    | `view_memberships`      | View group membership, audit  |
| 20    | `full`                  | All group permissions         |
| 21    | `owner`                 | Group owner (delete group)    |

### Wallet Roles (3 levels)

Wallet-level access controls sharing between users within a group.

| Role     | Permissions                        |
|----------|------------------------------------|
| `owner`  | Full control, manage members       |
| `editor` | Create/modify transactions         |
| `viewer` | Read-only access                   |

## Environment Variables

Key variables (see `.env.example` for the full list):

| Variable            | Default   | Description                   |
|---------------------|-----------|-------------------------------|
| `APP_ENV`           | production| `local`, `testing`, `production` |
| `HTTP_PORT`         | 8080      | Server port                   |
| `DB_HOST`           | localhost | PostgreSQL host               |
| `DB_PORT`           | 5432      | PostgreSQL port               |
| `DB_DATABASE`       | gofin     | Database name                 |
| `REDIS_HOST`        | localhost | Redis host                    |
| `AUTH_PROVIDER`     | local     | `local` or `disabled`         |
| `AUTH_JWT_SECRET`   | (required)| Must be 32+ characters        |
| `LOG_LEVEL`         | info      | `debug`, `info`, `warn`, `error` |
| `KEYCLOAK_URL`      |           | Keycloak base URL (optional)  |

## License

MIT

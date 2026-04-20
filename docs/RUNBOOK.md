# Gofin Production Runbook

## Overview
Gofin is a personal finance management API built in Go (Fiber v2 + PostgreSQL + Redis).

## Architecture
- **HTTP Framework**: Fiber v2
- **Database**: PostgreSQL 15+ via pgxpool
- **Cache**: Redis (optional — app runs without it)
- **Auth**: JWT access + refresh tokens, pluggable providers (local, Google, GitHub, Keycloak)
- **Real-time**: Server-Sent Events (SSE) for notifications
- **Decimal Math**: shopspring/decimal for financial precision

## Local Development

### Docker Development (recommended)

The fastest way to get running is with Docker, which starts the API, PostgreSQL, and Redis together:

```bash
# Start backend services (API + Postgres + Redis)
make docker-dev

# In a separate terminal, start the frontend dev server
make web-dev
```

- API runs on http://localhost:8080
- Frontend runs on http://localhost:5173 and proxies `/api` requests to the backend
- PostgreSQL is available on `localhost:5432`
- Redis is available on `localhost:6379`

### OAuth Development

For Keycloak-based OAuth, use the full compose file instead:

```bash
docker compose -f deployments/docker/docker-compose.yml up -d
```

This adds Keycloak on http://localhost:8081. Set `AUTH_PROVIDER=keycloak` in your `.env`.

### Manual Setup

If you prefer to run services individually:

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

### Test Types

| Test Type | What it tests | Requires Docker |
|-----------|---------------|----------------|
| Unit tests | Domain models, handlers, middleware, services, config | No |
| Integration tests | Auth, RBAC, wallet permissions against real PostgreSQL + Redis | Yes |
| E2E tests | Full browser flows via Playwright | Yes |

### Running Tests

```bash
# Unit tests only (no Docker needed)
make api-test-unit

# Integration tests (start infra first)
make api-test-integration-infra
make api-test-integration

# Full test suite in Docker (spins up Postgres + Redis + runs all tests)
make docker-test

# Web type check
make web-lint

# E2E tests (requires running backend)
cd web && bunx playwright test
```

### Integration Test Details

Integration tests use `docker-compose.test.yml` which provides PostgreSQL on port 5433 and Redis on port 6380 (non-standard ports to avoid conflicts). Four seed users are created per test run: owner, full, manage_transactions, and read_only roles.

### Test Infrastructure

See `docs/tests/TEST_REPORT.md` for the latest test results and RBAC coverage matrix.

## Docker Commands

### Quick Reference

```bash
# Daily development (API + Postgres + Redis)
make docker-dev
docker compose -f deployments/docker/docker-compose.dev.yml up -d
docker compose -f deployments/docker/docker-compose.dev.yml down

# OAuth development (API + Postgres + Redis + Keycloak)
docker compose -f deployments/docker/docker-compose.yml up -d
docker compose -f deployments/docker/docker-compose.yml down

# Self-hosted production (Caddy + API + Web + Postgres + Redis)
make docker-selfhost
docker compose -f deployments/docker/docker-compose.selfhost.yml up -d
docker compose -f deployments/docker/docker-compose.selfhost.yml down

# Run tests (Postgres + Redis + Test runner)
make docker-test
docker compose -f deployments/docker/docker-compose.test.yml up --abort-on-container-exit
docker compose -f deployments/docker/docker-compose.test.yml down -v
```

### Compose Files

| Compose File | Use Case | Services |
|--------------|----------|----------|
| `docker-compose.dev.yml` | Daily development | API + Postgres + Redis |
| `docker-compose.yml` | OAuth development | API + Postgres + Redis + Keycloak |
| `docker-compose.selfhost.yml` | Production deployment | Caddy + API + Web + Postgres + Redis |
| `docker-compose.test.yml` | CI/testing | Postgres + Redis + Test runner |

### Useful Docker Tips

```bash
# View API logs
docker compose -f deployments/docker/docker-compose.dev.yml logs -f api

# Rebuild after code changes
docker compose -f deployments/docker/docker-compose.dev.yml up -d --build

# Reset database volumes
docker compose -f deployments/docker/docker-compose.dev.yml down -v
```

## Deployment

### Environment Variables
| Variable | Required | Description |
|----------|----------|-------------|
| `APP_ENV` | Yes | `development`, `staging`, `production` |
| `APP_PORT` | No | HTTP port (default 8080) |
| `DATABASE_URL` | Yes | PostgreSQL DSN |
| `REDIS_ADDR` | No | Redis address (default localhost:6379) |
| `REDIS_PASSWORD` | No | Redis password |
| `REDIS_DB` | No | Redis DB number (default 0) |
| `AUTH_JWT_SECRET` | Yes | JWT signing secret |
| `AUTH_JWT_EXPIRY` | No | Access token expiry (default 15m) |
| `AUTH_REFRESH_EXPIRY` | No | Refresh token expiry (default 168h) |
| `AUTH_PROVIDER` | No | Auth provider: local, google, github, keycloak, disabled |
| `LOG_FORMAT` | No | `json` or `console` (default console) |
| `APP_DEBUG` | No | Enable debug logging |

### Database Migrations
```bash
# Run all pending migrations
go run ./cmd/migrate -dsn "$DATABASE_URL" -dir up

# Rollback all migrations
go run ./cmd/migrate -dsn "$DATABASE_URL" -dir down
```

### Starting the Server
```bash
go run ./cmd/server
```

## Monitoring

### Health Check
```
GET /health
```
Returns database and Redis connectivity status.

### API Documentation
```
GET /api/v1/docs        — HTML documentation
GET /api/v1/openapi.json — OpenAPI 3.0 spec
```

### Key Metrics to Monitor
- **Request latency**: P50, P95, P99 via access logs
- **Database connections**: Pool utilization (MaxConns, MinConns)
- **Redis connectivity**: Cache hit/miss ratio
- **SSE connections**: Active real-time clients
- **Error rate**: 5xx responses per minute

## Common Operations

### User Login Issues
1. Check JWT secret is consistent across instances
2. Verify user exists in `users` table
3. Check `auth_provider` configuration matches

### Database Connection Pool Exhaustion
1. Check `DB_MAX_OPEN_CONNS` setting (default 25)
2. Monitor slow queries via `pg_stat_statements`
3. Consider increasing pool or optimizing queries

### Redis Unavailable
App degrades gracefully — cache misses hit database directly. No action needed unless performance degrades.

### Reset a User's Password
```sql
-- Hash generated with bcrypt
UPDATE users SET password_hash = '$2a$10$...' WHERE email = 'user@example.com';
```

## Incident Response

### High Error Rate
1. Check `/health` endpoint
2. Review recent logs for panic/recovery events
3. Check database and Redis connectivity
4. Review recent deployments

### Database Failover
1. Verify PostgreSQL is reachable
2. Pool will reconnect automatically on re-establishment
3. If persistent, check `pgxpool` configuration and connection limits

### Security Incident
1. Rotate `AUTH_JWT_SECRET` — all existing sessions will be invalidated
2. Review audit logs: `SELECT * FROM audit_logs ORDER BY created_at DESC LIMIT 100`
3. Check for suspicious login patterns

## API Endpoints Summary

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | /health | No | Health check |
| POST | /api/v1/auth/login | No | Login |
| POST | /api/v1/auth/refresh | No | Refresh token |
| GET | /api/v1/auth/provider | No | Auth provider info |
| GET/PUT | /api/v1/users/me | Yes | Current user profile |
| CRUD | /api/v1/groups | Yes | User groups |
| CRUD | /api/v1/accounts | Yes | Wallets/accounts |
| CRUD | /api/v1/transactions | Yes | Transactions |
| CRUD | /api/v1/categories | Yes | Categories |
| CRUD | /api/v1/tags | Yes | Tags |
| CRUD | /api/v1/budgets | Yes | Budgets |
| CRUD | /api/v1/bills | Yes | Bills |
| GET | /api/v1/notifications/stream | Yes | SSE real-time notifications |
| GET | /api/v1/analytics/* | Yes | Financial analytics |
| GET | /api/v1/export/* | Yes | CSV/OFX export |
| GET | /api/v1/audit-logs | Yes | Audit trail |
| GET | /api/v1/admin/* | Yes | Admin operations |

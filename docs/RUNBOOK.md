# Gofin Production Runbook

## Overview
Gofin is a personal finance management API built in Go (Fiber v2 + PostgreSQL + Redis).

## Architecture
- **HTTP Framework**: Fiber v2
- **Database**: PostgreSQL 17 via pgxpool
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
| `docker-compose.selfhost.yml` | Production deployment | Caddy + API + Web + Postgres + Redis + Backup |
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

See `.env.example` for the full list. Key variables:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `APP_ENV` | Yes | `production` | `local`, `testing`, `production` |
| `APP_URL` | Yes | — | Public URL (e.g., `https://your-domain.com`) |
| `HTTP_PORT` | No | `8080` | Server port |
| `DB_HOST` | Yes | `localhost` | PostgreSQL host |
| `DB_PORT` | No | `5432` | PostgreSQL port |
| `DB_DATABASE` | Yes | `gofin` | Database name |
| `DB_USERNAME` | Yes | `gofin` | Database user |
| `DB_PASSWORD` | Yes | — | Database password |
| `REDIS_HOST` | No | `localhost` | Redis host |
| `REDIS_PORT` | No | `6379` | Redis port |
| `REDIS_PASSWORD` | No | — | Redis password |
| `AUTH_PROVIDER` | No | `local` | `local`, `google`, `github`, `keycloak`, `disabled` |
| `AUTH_JWT_SECRET` | Yes | — | JWT signing secret (32+ characters) |
| `AUTH_ALLOW_REGISTRATION` | No | `false` | Allow public self-registration |
| `ADMIN_EMAIL` | No | — | Auto-seed admin user on first startup |
| `ADMIN_PASSWORD` | No | — | Admin password (random if empty) |
| `LOG_LEVEL` | No | `info` | `debug`, `info`, `warn`, `error` |
| `LOG_FORMAT` | No | `json` | `json` or `console` |
| `CORS_ALLOWED_ORIGINS` | Yes | — | Allowed frontend origins |
| `DOMAIN` | No | `localhost` | Public domain for Caddy HTTPS |

### Database Migrations
```bash
# Run all pending migrations
go run ./cmd/migrate -dsn "$DATABASE_URL" -dir up

# Rollback all migrations
go run ./cmd/migrate -dsn "$DATABASE_URL" -dir down
```

### Self-Hosted Deployment

The `docker-compose.selfhost.yml` provides a production-ready stack:

1. Copy `.env.example` to `.env` and configure:
   - `DOMAIN` — your public domain (used by Caddy for HTTPS)
   - `AUTH_JWT_SECRET` — strong random string (32+ chars)
   - `DB_PASSWORD` — database password
   - `ADMIN_EMAIL` — admin user email (auto-seeded on first run)

2. Start the stack:
   ```bash
   make docker-selfhost
   ```

3. This starts: Caddy (reverse proxy + HTTPS), API, Web, PostgreSQL, Redis, Backup.

4. The API entrypoint automatically runs migrations and seeds the admin user.

5. Caddy auto-provisions HTTPS via Let's Encrypt for your `DOMAIN`.

### Backup & Restore

Backups run automatically via a cron container (daily at 03:00 UTC by default).

**Manual backup:**
```bash
docker compose -f deployments/docker/docker-compose.selfhost.yml exec backup /backup.sh
```

**List backups:**
```bash
docker compose -f deployments/docker/docker-compose.selfhost.yml exec backup ls -la /backups/
```

**Restore from backup:**
```bash
# Copy backup out of the container
docker cp <container_id>:/backups/gofin_YYYYMMDD_HHMMSS.sql.gz ./backup.sql.gz
gunzip ./backup.sql.gz

# Restore into a running PostgreSQL
docker compose -f deployments/docker/docker-compose.selfhost.yml exec -T postgres \
  psql -U gofin -d gofin < ./backup.sql
```

**Backup configuration:**
| Variable | Default | Description |
|----------|---------|-------------|
| `BACKUP_CRON_SCHEDULE` | `0 3 * * *` | Cron schedule (daily at 03:00 UTC) |
| `BACKUP_RETENTION_DAYS` | `30` | Days to keep backups |

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
| CRUD | /api/v1/wallets | Yes | Wallets (financial accounts) |
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

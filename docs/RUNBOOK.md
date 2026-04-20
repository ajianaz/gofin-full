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

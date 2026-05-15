# Security

Security features, hardening measures, and best practices in Gofin.

## Authentication

### JWT Tokens

- **Access token:** Short-lived (configurable, default 60 minutes)
- **Refresh token:** Long-lived (configurable, default 30 days)
- **Library:** `golang-jwt/v5` with RS256 signing
- **Storage:** Access token in localStorage, refresh token as httpOnly cookie

### Password Policy

Registration and password change enforce:

| Requirement | Value |
|------------|-------|
| Minimum length | 8 characters |
| Character types | At least 3 of 4: uppercase, lowercase, digit, special |
| Hashing | bcrypt (cost factor 10) |

### OAuth2 Providers

| Provider | Scope |
|----------|-------|
| Google | `openid`, `email`, `profile` |
| GitHub | `user:email`, `read:user` |
| Keycloak OIDC | Configurable (realm-based) |

### Login Lockout

Failed login attempts are tracked per **email + IP combination** to prevent brute force:

- Lockout after configurable failed attempts
- Separate tracking per email/IP pair (prevents one attacker from locking out all users)
- **Redis-backed** with automatic **in-memory fallback** when Redis is unavailable (sliding window with periodic cleanup)

## Request Security

### CORS

- **Production:** Restricted to configured `CORS_ALLOWED_ORIGINS` (defaults to domain)
- **Development:** Allows `localhost:5173` (Vite dev server)
- **No wildcard** (`*`) in production — ever

### Rate Limiting

| Setting | Default | Description |
|---------|---------|-------------|
| `RATE_LIMIT_MAX` | 100 (prod), 20 (selfhost) | Max requests per window |
| `RATE_LIMIT_WINDOW_SECONDS` | 60 | Window duration in seconds |

Rate limiting uses a sliding window algorithm backed by Redis, with automatic in-memory fallback when Redis is unavailable.

### Request Body Limit

- **Default:** 10 MB (`MAX_REQUEST_BODY_BYTES`)
- Applied at the Fiber framework level
- Prevents oversized payload attacks

### HSTS (HTTP Strict Transport Security)

- Enabled automatically in production (`APP_ENV=production`)
- Not applied in development mode

## API Security

### Parameterized Queries

All database queries use parameterized SQL (`$1, $2, ...`). No string interpolation in queries.

```go
// ✅ Correct — parameterized
db.Query("SELECT * FROM wallets WHERE id = $1", walletID)

// ❌ Never done — string interpolation
db.Query("SELECT * FROM wallets WHERE id = " + walletID)
```

### Financial Precision

All monetary values use `shopspring/decimal.Decimal` — never `float64`. This prevents floating-point rounding errors in financial calculations.

### Error Responses

- **Production:** Generic error messages without stack traces or internal details
- **Development:** Full error details for debugging (`APP_DEBUG=true`)

### Recovery Middleware

Global panic recovery catches unexpected errors and returns a generic 500 response. Stack traces are only logged server-side, never sent to clients.

## Self-Host Defaults

The self-hosted Docker configuration uses secure defaults:

| Setting | Value | Reason |
|---------|-------|--------|
| `APP_ENV` | `production` | Enables all security features |
| `APP_DEBUG` | `false` | Hides error details |
| `AUTH_ALLOW_REGISTRATION` | `false` | Prevents open registration |
| `DISABLE_PROMETHEUS` | `true` | Disables metrics endpoint |
| `RATE_LIMIT_MAX` | `20` | Conservative rate limiting |
| `CORS_ALLOWED_ORIGINS` | Domain only | No wildcard CORS |

## Secrets Management

### Required Secrets

| Secret | Purpose | Generate With |
|--------|---------|--------------|
| `AUTH_JWT_SECRET` | JWT signing key | `openssl rand -hex 32` |
| `STATIC_CRON_TOKEN` | Internal cron auth | `openssl rand -hex 16` |
| `DB_PASSWORD` | PostgreSQL password | `openssl rand -hex 16` |

### Optional Secrets

| Secret | Purpose |
|--------|---------|
| `GOOGLE_CLIENT_ID` / `GOOGLE_CLIENT_SECRET` | Google OAuth2 |
| `GITHUB_CLIENT_ID` / `GITHUB_CLIENT_SECRET` | GitHub OAuth2 |
| `KEYCLOAK_CLIENT_ID` / `KEYCLOAK_CLIENT_SECRET` | Keycloak OIDC |
| `REDIS_PASSWORD` | Redis authentication |

::: warning Never commit secrets
All secrets are configured via environment variables (`.env` file). Never commit `.env` to version control. A `.env.example` template is provided with placeholder values.
:::

## Infrastructure Security

### Docker Isolation

- Each service runs in its own container
- Non-root users where supported
- Resource limits (memory) on all containers
- Internal network communication only (ports not exposed except 80/443 via Caddy)

### Database

- PostgreSQL 17 with `sslmode=prefer`
- Connection pooling with configurable limits
- Separate database per deployment (no shared databases)

### Backups

- Automated daily backups (configurable cron schedule)
- Configurable retention period (default: 30 days)
- Stored in a Docker volume (mount to host for persistence)

## Security Checklist

Use this checklist when deploying Gofin:

- [ ] `AUTH_JWT_SECRET` is set to a random 32+ char string
- [ ] `STATIC_CRON_TOKEN` is set to a random string
- [ ] `DB_PASSWORD` is different from the default
- [ ] `AUTH_ALLOW_REGISTRATION` is `false` (unless you want public registration)
- [ ] `DOMAIN` is set to your public domain (for CORS)
- [ ] `APP_DEBUG` is `false` in production
- [ ] `ADMIN_PASSWORD` is a strong password (or let it auto-generate)
- [ ] Backups are mounted to a persistent volume

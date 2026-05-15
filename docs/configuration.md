# Configuration

Complete reference for all environment variables. Copy `.env.example` to `.env` and customize.

```bash
cp .env.example .env
```

## Application Core

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_ENV` | `production` | Environment: `local`, `production` |
| `APP_DEBUG` | `false` | Enable debug logs and error details in responses |
| `APP_URL` | `http://localhost` | Public URL of the application |
| `TZ` | `UTC` | Server timezone |

::: tip
Set `APP_ENV=local` for development. It relaxes CORS and enables debug mode.
:::

## HTTP Server

| Variable | Default | Description |
|----------|---------|-------------|
| `HTTP_PORT` | `8080` | API server port |
| `HTTP_HOST` | `0.0.0.0` | Bind address |

## Database (PostgreSQL)

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_DATABASE` | `gofin` | Database name |
| `DB_USERNAME` | `gofin` | Database user |
| `DB_PASSWORD` | *(empty)* | Database password |
| `DB_SSL_MODE` | `prefer` | SSL mode: `disable`, `prefer`, `require` |
| `DB_SCHEMA` | `public` | PostgreSQL schema |
| `DB_MAX_OPEN_CONNS` | `5` | Max open connections |
| `DB_MAX_IDLE_CONNS` | `2` | Max idle connections |
| `DB_CONN_MAX_LIFETIME` | `300` | Connection max lifetime (seconds) |

## Redis

| Variable | Default | Description |
|----------|---------|-------------|
| `REDIS_HOST` | `localhost` | Redis host |
| `REDIS_PORT` | `6379` | Redis port |
| `REDIS_PASSWORD` | *(empty)* | Redis password |
| `REDIS_CACHE_DB` | `1` | Redis database index for cache |

## Authentication

| Variable | Default | Description |
|----------|---------|-------------|
| `AUTH_PROVIDER` | `local` | Auth provider: `local`, `keycloak` |
| `AUTH_JWT_SECRET` | *(required)* | JWT signing key — **must be ≥ 32 chars** |
| `AUTH_JWT_EXPIRY_MINUTES` | `60` | Access token expiry |
| `AUTH_REFRESH_EXPIRY_DAYS` | `30` | Refresh token expiry |
| `AUTH_ALLOW_REGISTRATION` | `false` | Allow public self-registration |
| `STATIC_CRON_TOKEN` | *(required)* | Token for internal cron endpoints |
| `ALLOW_2FA_BYPASS` | `false` | Allow skipping 2FA (not recommended) |

::: warning Required
`AUTH_JWT_SECRET` and `STATIC_CRON_TOKEN` have no safe defaults. The server will **refuse to start** if `AUTH_JWT_SECRET` is not changed from the placeholder.
:::

## Admin Seed (First Run)

| Variable | Default | Description |
|----------|---------|-------------|
| `ADMIN_EMAIL` | *(empty)* | Auto-create admin user on first startup |
| `ADMIN_PASSWORD` | *(empty)* | Admin password. If empty, auto-generates a 16-char password printed to logs |

## OAuth — Google

| Variable | Default | Description |
|----------|---------|-------------|
| `GOOGLE_CLIENT_ID` | *(empty)* | Google OAuth2 client ID |
| `GOOGLE_CLIENT_SECRET` | *(empty)* | Google OAuth2 client secret |

## OAuth — GitHub

| Variable | Default | Description |
|----------|---------|-------------|
| `GITHUB_CLIENT_ID` | *(empty)* | GitHub OAuth2 client ID |
| `GITHUB_CLIENT_SECRET` | *(empty)* | GitHub OAuth2 client secret |

## Keycloak OIDC

| Variable | Default | Description |
|----------|---------|-------------|
| `KEYCLOAK_URL` | `http://localhost:8088` | Keycloak server URL |
| `KEYCLOAK_REALM` | `gofin` | Keycloak realm |
| `KEYCLOAK_CLIENT_ID` | `gofin-api` | Keycloak client ID |
| `KEYCLOAK_CLIENT_SECRET` | *(empty)* | Keycloak client secret |

## Security

| Variable | Default | Description |
|----------|---------|-------------|
| `RATE_LIMIT_MAX` | `100` | Max requests per window |
| `RATE_LIMIT_WINDOW_SECONDS` | `60` | Rate limit window (seconds) |
| `MAX_REQUEST_BODY_BYTES` | `10485760` | Max request body (10 MB) |
| `CORS_ALLOWED_ORIGINS` | `http://localhost:5173` | Comma-separated allowed origins |
| `DISABLE_PROMETHEUS` | `false` | Disable `/metrics` endpoint |

## Logging

| Variable | Default | Description |
|----------|---------|-------------|
| `LOG_LEVEL` | `info` | Log level: `debug`, `info`, `warn`, `error` |
| `LOG_FORMAT` | `json` | Log format: `json`, `text` |

## Feature Flags

| Variable | Default | Description |
|----------|---------|-------------|
| `FEATURE_EXPORT` | `true` | Enable CSV/OFX export |
| `FEATURE_WEBHOOKS` | `true` | Enable webhook endpoints |
| `FEATURE_HANDLE_DEBTS` | `true` | Enable debt tracking |
| `FEATURE_EXPRESSION_ENGINE` | `true` | Enable expression evaluation |
| `FEATURE_RUNNING_BALANCE` | `true` | Enable running balance calculation |

## Business Logic

| Variable | Default | Description |
|----------|---------|-------------|
| `MAX_UPLOAD_SIZE` | `1073741824` | Max file upload (1 GB) |
| `ALLOW_WEBHOOKS` | `false` | Enable outgoing webhooks |
| `WEBHOOK_MAX_ATTEMPTS` | `3` | Max webhook retry attempts |
| `ENABLE_EXTERNAL_RATES` | `false` | Enable external exchange rate fetching |
| `ENABLE_EXCHANGE_RATES` | `false` | Enable exchange rate features |

## Database Backup (Self-Host)

| Variable | Default | Description |
|----------|---------|-------------|
| `BACKUP_CRON_SCHEDULE` | `0 3 * * *` | Cron schedule for automated backups |
| `BACKUP_RETENTION_DAYS` | `30` | Days to keep backup files |

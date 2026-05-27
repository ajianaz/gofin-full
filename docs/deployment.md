# Deployment

Production and development deployment guides for Gofin.

## Docker Images

Gofin publishes pre-built images to two registries:

| Registry | Images | Use Case |
|----------|--------|----------|
| **Docker Hub** | `ajianaz/gofin-api`, `ajianaz/gofin-web` | Production (`:latest`, `:v0.x.x`) |
| **GHCR** | `ghcr.io/ajianaz/gofin-api`, `ghcr.io/ajianaz/gofin-web` | Development (`:develop`) and production |

Both registries carry the same multi-arch images (amd64 + arm64). Choose based on your preference.

## Production — Self-Host (Caddy)

All-in-one deployment with built-in reverse proxy. No external dependencies.

### Architecture

```
Internet → Caddy (443) → API (8080)
                       → Web (static)
                  PostgreSQL (5432)
                  Redis (6379)
                  Backup (cron)
```

### Quick Deploy

```bash
git clone https://github.com/ajianaz/gofin-full.git
cd gofin-full
cp api/.env.example .env
# Edit .env — see Configuration page
docker compose -f deployments/docker/docker-compose.selfhost.yml up -d
```

> Default images: `ajianaz/gofin-api:latest` / `ajianaz/gofin-web:latest` (Docker Hub).
> Switch to GHCR: `DOCKER_IMAGE_API=ghcr.io/ajianaz/gofin-api:latest`

### What's Started

1. Pulls pre-built images (no build required)
2. Starts 6 Docker containers via `docker-compose.selfhost.yml`
3. Runs database migrations automatically
4. Seeds the admin user (if `ADMIN_EMAIL` is set)
5. Caddy provisions HTTPS via Let's Encrypt

### Docker Compose Override

To customize without modifying the main compose file, create a `docker-compose.override.yml`:

```yaml
# docker-compose.override.yml
services:
  api:
    environment:
      - LOG_LEVEL=debug
    deploy:
      resources:
        limits:
          memory: 512M
```

### Volume Mounts

| Volume | Purpose | Path in Container |
|--------|---------|-------------------|
| `pg_data` | PostgreSQL data persistence | `/var/lib/postgresql/data` |
| `redis_data` | Redis persistence | `/data` |
| `caddy_data` | Caddy certificates + config | `/data` |
| `caddy_config` | Caddy configuration | `/config` |
| `backups` | Database backup files | `/backups` |
| `uploads` | User file attachments | `/uploads` |

::: warning Persistent Storage
Make sure to mount the `backups` and `uploads` volumes to your host for data persistence. Without host mounts, data is lost when containers are removed.
:::

## Caddy (HTTPS)

Caddy automatically provisions SSL/TLS certificates via [Let's Encrypt](https://letsencrypt.org/).

### Default Caddyfile

```
{$DOMAIN} {
    reverse_proxy /api/* api:8080
    encode gzip

    root * /srv
    file_server {
        precompressed br gzip
    }

    try_files {path} /index.html
}
```

### Custom Caddy Configuration

To customize the Caddyfile, mount your own:

```yaml
# docker-compose.override.yml
services:
  caddy:
    volumes:
      - ./my-caddyfile:/etc/caddy/Caddyfile:ro
```

### DNS Requirements

Caddy needs port 80 and 443 accessible from the internet for certificate provisioning:

1. Point your domain's A record to your server's public IP
2. Ensure ports 80 (HTTP) and 443 (HTTPS) are open in your firewall
3. Caddy handles the rest automatically

## Development — Traefik (GHCR)

Lightweight dev deployment behind an external Traefik reverse proxy. Images pulled from GHCR (no Docker Hub rate limits).

### Architecture

```
Internet → Traefik (existing) → API (8080)
                              → Web (3000)
                         PostgreSQL (5432)
```

### Quick Deploy

```bash
git clone https://github.com/ajianaz/gofin-full.git
cd gofin-full
cp api/.env.example .env
# Edit .env — REQUIRED: AUTH_JWT_SECRET, DOMAIN, STATIC_CRON_TOKEN
# Optional: TRAEFIK_NETWORK=your-traefik-network-name
docker compose -f deployments/docker/docker-compose.traefik.yml up -d
```

> Default images: `ghcr.io/ajianaz/gofin-api:develop` / `ghcr.io/ajianaz/gofin-web:develop` (GHCR).
> Requires an existing Traefik instance. Set `TRAEFIK_NETWORK` to match your Traefik network.

### What's Included

- PostgreSQL (with health check + volume persistence)
- API server (behind Traefik labels)
- Web frontend (behind Traefik labels)
- No Redis (optional — set `REDIS_HOST` if available)
- No backup container

### Traefik Network

The compose file connects to an external Docker network for Traefik routing:

```bash
# Default: traefik-public
# Override via env:
TRAEFIK_NETWORK=my-proxy docker compose -f deployments/docker/docker-compose.traefik.yml up -d
```

## Database Migrations

Migrations run automatically on API startup. To run them manually:

```bash
make migrate
```

### Migration Files

Migrations are stored in `api/migrations/` and use a sequential naming scheme:

```
000001_init.up.sql
000001_init.down.sql
000002_add_budgets.up.sql
000002_add_budgets.down.sql
...
```

## Backups

### Automated Backups

Self-host deployments include an automated backup container:

- **Schedule:** Daily at 03:00 UTC (configurable via `BACKUP_CRON_SCHEDULE`)
- **Retention:** 30 days (configurable via `BACKUP_RETENTION_DAYS`)
- **Format:** Plain SQL dump (gzipped)
- **Storage:** Docker volume `backups`

### Manual Backup

```bash
docker compose -f deployments/docker/docker-compose.selfhost.yml exec backup pg_dump -U gofin gofin | gzip > backup.sql.gz
```

### Restore from Backup

```bash
gunzip -c backup.sql.gz | docker compose -f deployments/docker/docker-compose.selfhost.yml exec -T postgres psql -U gofin gofin
```

## Monitoring

### Health Check

```bash
curl https://your-domain/health
```

Returns:

```json
{
  "status": "healthy",
  "timestamp": "2026-05-15T00:00:00Z"
}
```

### Prometheus Metrics

If `DISABLE_PROMETHEUS=false` (default):

```
https://your-domain/metrics
```

Available metrics:
- HTTP request count, duration, status codes
- Active connections (DB, Redis)
- Custom business metrics

## Updating

### Standard Update

```bash
cd gofin-full
git pull origin main
docker compose -f deployments/docker/docker-compose.selfhost.yml pull
docker compose -f deployments/docker/docker-compose.selfhost.yml up -d
```

Docker Compose pulls the latest images and restarts containers. Data volumes persist.

### Zero-Downtime Update

For production environments, update one service at a time:

```bash
cd gofin-full
git pull origin main

# Pull latest images
docker compose -f deployments/docker/docker-compose.selfhost.yml pull

# Restart API first, wait for healthy, then web
docker compose -f deployments/docker/docker-compose.selfhost.yml up -d api
# Wait for API health check to pass, then:
docker compose -f deployments/docker/docker-compose.selfhost.yml up -d web
```

## Troubleshooting

### API won't start

```
Error: AUTH_JWT_SECRET must be at least 32 characters
```

→ Generate a proper secret: `openssl rand -hex 32`

### Caddy can't get certificates

1. Check DNS: `dig your-domain.com`
2. Check ports: `nc -zv your-domain.com 80` and `nc -zv your-domain.com 443`
3. Check Caddy logs: `docker compose logs caddy`

### Database connection refused

1. Check PostgreSQL is running: `docker compose ps`
2. Check DB credentials in `.env` match the compose file
3. Check network: containers must be on the same Docker network

### "Admin password" not in logs

The auto-generated password is printed once on first startup. Check:

```bash
docker compose logs api | grep "admin"
```

If you missed it, reset the admin password via the API or delete the user and let the seed recreate them.

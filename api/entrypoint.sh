#!/bin/sh
set -e

echo "[entrypoint] Running database migrations..."
/app/migrate -dsn "${DB_DSN}" -path /app/migrations/postgres

echo "[entrypoint] Seeding admin user (if needed)..."
if [ -n "$ADMIN_EMAIL" ]; then
  /app/seed -dsn "${DB_DSN}" -email "${ADMIN_EMAIL}" -password "${ADMIN_PASSWORD:-}"
else
  echo "[entrypoint] ADMIN_EMAIL not set, skipping admin seed."
  echo "[entrypoint] Set ADMIN_EMAIL and optionally ADMIN_PASSWORD to auto-create an admin user."
fi

echo "[entrypoint] Starting server..."
exec /app/gofin

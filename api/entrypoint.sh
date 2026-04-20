#!/bin/sh
set -e

echo "[entrypoint] Running database migrations..."
/app/migrate -dsn "${DB_DSN}" -path /app/migrations/postgres

echo "[entrypoint] Starting server..."
exec /app/gofin

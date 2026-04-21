#!/bin/sh
set -e

# Gofin database backup script
# Designed to run via cron inside the backup container.
# Backs up to /backups with automatic rotation.

DB_HOST="${DB_HOST:-postgres}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_DATABASE:-gofin}"
DB_USER="${DB_USERNAME:-gofin}"
DB_PASS="${DB_PASSWORD:-gofin_secret}"
BACKUP_DIR="/backups"
RETENTION_DAYS="${BACKUP_RETENTION_DAYS:-30}"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
FILE="${BACKUP_DIR}/${DB_NAME}_${TIMESTAMP}.sql.gz"

mkdir -p "${BACKUP_DIR}"

echo "[backup] Starting backup of ${DB_NAME}..."

PGPASSWORD="${DB_PASS}" pg_dump \
  -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" \
  --format=plain \
  --no-owner \
  --no-privileges \
  "${DB_NAME}" | gzip > "${FILE}"

SIZE=$(du -h "${FILE}" | cut -f1)
echo "[backup] Done: ${FILE} (${SIZE})"

# Rotate old backups
echo "[backup] Removing backups older than ${RETENTION_DAYS} days..."
find "${BACKUP_DIR}" -name "${DB_NAME}_*.sql.gz" -mtime +${RETENTION_DAYS} -delete

REMAINING=$(find "${BACKUP_DIR}" -name "${DB_NAME}_*.sql.gz" | wc -l)
echo "[backup] ${REMAINING} backups retained"

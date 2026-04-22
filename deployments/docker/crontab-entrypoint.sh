#!/bin/sh
set -e

# Write cron schedule with backup command
echo "${CRON_SCHEDULE} /backup.sh >> /var/log/cron.log 2>&1" > /etc/crontabs/root

echo "[backup] Cron schedule: ${CRON_SCHEDULE}"
echo "[backup] Retention: ${BACKUP_RETENTION_DAYS:-30} days"

# Run backup immediately on first start (unless SKIP_FIRST_BACKUP is set)
if [ "${SKIP_FIRST_BACKUP}" != "true" ]; then
  echo "[backup] Running initial backup..."
  /backup.sh
fi

echo "[backup] Starting cron daemon..."
exec crond -f -l 2

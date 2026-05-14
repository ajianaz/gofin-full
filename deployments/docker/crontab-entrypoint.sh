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
# On Docker Desktop (macOS/Windows), crond may fail due to setpgid restrictions.
# The initial backup already ran above; cron scheduling works on Linux hosts.
if crond -f -l 2 2>/dev/null; then
  : # cron running
else
  echo "[backup] Cron daemon unavailable (Docker Desktop). Backup completed. Use external scheduler for automated backups."
  # Keep container alive with sleep so docker considers it running
  tail -f /dev/null
fi

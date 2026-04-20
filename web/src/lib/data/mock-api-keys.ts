import type { ApiKey } from '$lib/types/index.js';

export const mockApiKeys: ApiKey[] = [
	{ id: 'ak1', title: 'Gofin Mobile App', key: 'gofin_mobil_*************', created_at: '2026-04-10T14:00:00Z', last_used_at: '2026-04-15T09:00:00Z' },
	{ id: 'ak2', title: 'Automation Script', key: 'gofin_auto_*************', created_at: '2026-03-15T10:00:00Z', last_used_at: '2026-04-14T08:00:00Z' },
	{ id: 'ak3', title: 'Backup Script', key: 'gofin_backup_*************', created_at: '2026-02-20T12:00:00Z' }
];

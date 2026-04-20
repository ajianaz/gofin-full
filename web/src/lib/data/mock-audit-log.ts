import type { AuditLogEntry } from '$lib/types/index.js';

export const mockAuditLog: AuditLogEntry[] = [
	{ id: 'al1', action: 'user.login', user_email: 'anaz@gofin.app', entity_type: 'user', entity_id: '1', changes: 'Login berhasil dari 192.168.1.100', created_at: '2026-04-15T08:30:00Z' },
	{ id: 'al2', action: 'transaction.create', user_email: 'anaz@gofin.app', entity_type: 'transaction', entity_id: '1', changes: 'Makan Siang di Warteg -Rp 25.000', created_at: '2026-04-15T09:15:00Z' },
	{ id: 'al3', action: 'budget.create', user_email: 'anaz@gofin.app', entity_type: 'budget', entity_id: 'b1', changes: 'Anggaran April 2026 - Rp 5.000.000', created_at: '2026-04-01T00:00:00Z' },
	{ id: 'al4', action: 'api_key.create', user_email: 'anaz@gofin.app', entity_type: 'api_key', entity_id: 'ak1', changes: 'API Key: Gofin Mobile App', created_at: '2026-04-10T14:00:00Z' },
	{ id: 'al5', action: 'user.update', user_email: 'anaz@gofin.app', entity_type: 'user', entity_id: '1', changes: 'Nama diubah: Anaz S. Aji', created_at: '2026-03-20T10:00:00Z' },
	{ id: 'al6', action: 'group.create', user_email: 'anaz@gofin.app', entity_type: 'group', entity_id: 'g1', changes: 'Grup Keuangan Pribadi dibuat', created_at: '2026-01-01T00:00:00Z' },
	{ id: 'al7', action: 'piggy_bank.create', user_email: 'anaz@gofin.app', entity_type: 'piggy_bank', entity_id: 'pb1', changes: 'Tabungan Liburan Bali - Target Rp 15.000.000', created_at: '2026-01-01T00:00:00Z' },
	{ id: 'al8', action: 'recurring.create', user_email: 'anaz@gofin.app', entity_type: 'recurring', entity_id: 'rt1', changes: 'Recurring: Gaji Bulanan Rp 8.500.000', created_at: '2026-01-10T00:00:00Z' },
	{ id: 'al9', action: 'rule.update', user_email: 'anaz@gofin.app', entity_type: 'rule', entity_id: 'r1', changes: 'Rule diaktifkan: Makanan > Warteg', created_at: '2026-02-15T11:00:00Z' },
	{ id: 'al10', action: 'currency.update', user_email: 'anaz@gofin.app', entity_type: 'currency', entity_id: 'c1', changes: 'IDR enabled', created_at: '2026-01-01T00:00:00Z' }
];

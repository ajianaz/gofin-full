import { api } from './client.js';
import { unwrapMany } from './helpers.js';

interface AdminUser {
	id: string;
	email: string;
	name: string;
	role: string;
	is_active: boolean;
	created_at: string;
}

interface AuditLogEntry {
	id: string;
	action: string;
	user_email: string;
	entity_type: string;
	entity_id: string;
	changes: string;
	created_at: string;
}

export const adminService = {
	async listUsers(): Promise<AdminUser[]> {
		const res = await api.get<{ data: { id: string; attributes: Record<string, unknown> }[] }>('/admin/users');
		return unwrapMany<AdminUser>(res).map((u) => ({
			...u,
			email: (u as any).email ?? '',
			name: (u as any).name ?? '',
			role: (u as any).role ?? 'user',
			is_active: (u as any).is_active ?? true,
			created_at: (u as any).created_at ?? ''
		}));
	},

	async listAuditLogs(entityType?: string): Promise<AuditLogEntry[]> {
		const params = new URLSearchParams();
		if (entityType) params.set('entity_type', entityType);
		const qs = params.toString();
		const res = await api.get<{ data: { id: string; attributes: Record<string, unknown> }[] }>(
			`/audit-logs${qs ? '?' + qs : ''}`
		);
		return unwrapMany<AuditLogEntry>(res).map((l) => ({
			id: l.id,
			action: (l as any).action ?? '',
			user_email: (l as any).user_id ?? (l as any).user_email ?? '',
			entity_type: (l as any).entity_type ?? '',
			entity_id: String((l as any).entity_id ?? ''),
			changes: (l as any).new_value ?? (l as any).changes ?? '',
			created_at: (l as any).created_at ?? ''
		}));
	}
};

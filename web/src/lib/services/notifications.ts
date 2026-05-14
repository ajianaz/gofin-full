import { api } from './client.js';
import { unwrapMany, unwrapOne } from './helpers.js';
import type { Notification } from '$lib/types/domain.js';

export const notificationService = {
	async list(): Promise<Notification[]> {
			const res = await api.get<{ data: { id: string; attributes: Record<string, unknown> }[] | null }>('/notifications');
		return unwrapMany<Notification>(res).map((n) => ({
			id: n.id,
			channel: (n as any).channel || '',
			type: (n as any).type || '',
			title: (n as any).title || '',
			message: (n as any).message || '',
			read: !!((n as any).read)
		}));
	},

	async listUnread(): Promise<Notification[]> {
			const res = await api.get<{ data: { id: string; attributes: Record<string, unknown> }[] | null }>('/notifications/unread');
		return unwrapMany<Notification>(res).map((n) => ({
			id: n.id,
			channel: (n as any).channel || '',
			type: (n as any).type || '',
			title: (n as any).title || '',
			message: (n as any).message || '',
			read: !!((n as any).read)
		}));
	},

	async markRead(id: string): Promise<void> {
		await api.put(`/notifications/${id}/read`);
	},

	async markAllRead(): Promise<void> {
		await api.put('/notifications/read-all');
	}
};

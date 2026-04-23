import { api } from './client.js';
import { unwrapMany } from './helpers.js';
import type { UserGroup } from '$lib/types/domain.js';

export const groupService = {
	async list(): Promise<UserGroup[]> {
		const res = await api.get<{ data: { id: string; attributes: Record<string, unknown> }[] }>(
			'/groups'
		);
		return unwrapMany<UserGroup>(res).map((g) => ({
			id: g.id,
			title: (g as any).title ?? '',
			member_count: (g as any).member_count ?? 0,
			is_current: (g as any).is_current ?? false
		}));
	},

	async create(data: { title: string }): Promise<UserGroup> {
		const res = await api.post<{ data: { id: string; attributes: Record<string, unknown> } }>(
			'/groups',
			data
		);
		const g = { id: res.data.id, ...res.data.attributes } as UserGroup & { id: string };
		return { id: g.id, title: (g as any).title ?? data.title, member_count: 1, is_current: false };
	},

	async switch(userGroupId: string): Promise<void> {
		await api.post('/groups/switch', { user_group_id: userGroupId });
	}
};

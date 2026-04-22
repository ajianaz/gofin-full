import { api } from './client.js';
import { unwrapMany, unwrapOne } from './helpers.js';
import type { Tag } from '$lib/types/domain.js';

export const tagService = {
	async list(): Promise<Tag[]> {
		const res = await api.get<{ data: { id: string; attributes: Record<string, unknown> }[] }>('/tags');
		return unwrapMany<Tag>(res).map((t) => ({
			...t,
			date: (t as any).date || new Date().toISOString().split('T')[0]
		}));
	},

	async create(data: { tag: string; date?: string }): Promise<Tag> {
		const payload: Record<string, string> = { tag: data.tag };
		if (data.date) {
			payload.date = new Date(data.date).toISOString();
		}
		const res = await api.post<{ data: { id: string; attributes: Record<string, unknown> } }>('/tags', payload);
		const t = unwrapOne<Tag>(res);
		return { ...t, date: (t as any).date || new Date().toISOString().split('T')[0] };
	}
};

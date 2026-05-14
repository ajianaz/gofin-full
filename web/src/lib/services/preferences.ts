import { api } from './client.js';
import type { PreferenceItem } from '$lib/types/domain.js';

export const preferenceService = {
	async list(): Promise<PreferenceItem[]> {
			const res = await api.get<{ data: { id: string; attributes: { name: string; data: string } }[] | null }>('/preferences');
		return (res.data || []).map((p) => ({
			id: p.id,
			name: p.attributes.name,
			data: p.attributes.data
		}));
	},

	async get(name: string): Promise<PreferenceItem> {
		const res = await api.get<{ data: { id: string; attributes: { name: string; data: string } } }>(`/preferences/${name}`);
		return { id: res.data.id, name: res.data.attributes.name, data: res.data.attributes.data };
	},

	async set(name: string, data: string): Promise<void> {
		await api.post('/preferences', { name, data });
	},

	async delete(name: string): Promise<void> {
		await api.delete(`/preferences/${name}`);
	}
};

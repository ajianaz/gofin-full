import { api } from './client.js';
import type { ApiKeyListItem, ApiKeyCreateResponse } from '$lib/types/domain.js';

export const apiKeyService = {
	async list(): Promise<ApiKeyListItem[]> {
		const res = await api.get<{ data: ApiKeyListItem[] }>('/api-keys');
		return (res.data || []).map((k) => ({
			id: k.id,
			name: k.name,
			key_prefix: k.key_prefix || '',
			last_used: k.last_used || '',
			created_at: k.created_at || ''
		}));
	},

	async create(name: string): Promise<ApiKeyCreateResponse> {
		const res = await api.post<{ data: ApiKeyCreateResponse }>('/api-keys', { name });
		return res.data;
	},

	async delete(id: string): Promise<void> {
		await api.delete(`/api-keys/${id}`);
	}
};

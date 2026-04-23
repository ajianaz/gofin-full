import { api } from './client.js';
import { unwrapMany, unwrapOne } from './helpers.js';
import type { Transaction } from '$lib/types/domain.js';

interface TransactionListMeta {
	meta: { pagination: { total: number; count: number; per_page: number; current_page: number; total_pages: number } };
}

export const transactionService = {
	async list(params?: { page?: number; per_page?: number; start?: string; end?: string; type?: string }): Promise<{ data: Transaction[]; meta?: TransactionListMeta['meta'] }> {
		const query = new URLSearchParams();
		if (params?.page) query.set('page', String(params.page));
		if (params?.per_page) query.set('per_page', String(params.per_page));
		if (params?.start) query.set('start', params.start);
		if (params?.end) query.set('end', params.end);
		if (params?.type) query.set('type', params.type);
		const qs = query.toString();
		const path = `/transactions${qs ? `?${qs}` : ''}`;
		const res = await api.get<{ data: { id: string; attributes: Record<string, unknown> }[] } & TransactionListMeta>(path);
		const items = unwrapMany<Transaction>(res).map((t) => ({
			...t,
			amount: (t as any).amount || '0',
			date: (t as any).date || (t as any).created_at || '',
			description: (t as any).description || (t as any).group_title || '',
			tags: []
		}));
		return { data: items, meta: res.meta };
	},

	async create(data: {
		type: string;
		description?: string;
		amount: string;
		source_id: string;
		destination_id: string;
		date?: string;
		category_ids?: string[];
		tag_ids?: string[];
	}): Promise<Transaction> {
		const payload: Record<string, unknown> = { ...data };
		if (data.date) {
			payload.date = new Date(data.date).toISOString();
		}
		const res = await api.post<{ data: { id: string; attributes: Record<string, unknown> } }>('/transactions', payload);
		return unwrapOne<Transaction>(res);
	},

	async update(id: string, data: { description?: string }): Promise<void> {
		await api.put(`/transactions/${id}`, data);
	},

	async delete(id: string): Promise<void> {
		await api.delete(`/transactions/${id}`);
	},

	async split(id: string, splits: Array<{ amount: string; description?: string }>): Promise<void> {
		await api.post(`/transactions/${id}/split`, { splits });
	}
};

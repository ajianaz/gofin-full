import { api } from './client.js';
import { unwrapMany, unwrapOne } from './helpers.js';
import type { Category } from '$lib/types/domain.js';

export const categoryService = {
	async list(): Promise<Category[]> {
		const res = await api.get<{ data: { id: string; attributes: Record<string, unknown> }[] }>('/categories');
		return unwrapMany<Category>(res).map((c) => ({
			...c,
			type: 'expense',
			transaction_count: 0
		}));
	},

	async create(data: { name: string }): Promise<Category> {
		const res = await api.post<{ data: { id: string; attributes: Record<string, unknown> } }>('/categories', data);
		const c = unwrapOne<Category>(res);
		return { ...c, type: 'expense', transaction_count: 0 };
	},

	async update(id: string, data: { name: string }): Promise<void> {
		await api.put(`/categories/${id}`, data);
	},

	async delete(id: string): Promise<void> {
		await api.delete(`/categories/${id}`);
	}
};

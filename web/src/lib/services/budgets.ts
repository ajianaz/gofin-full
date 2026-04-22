import { api } from './client.js';
import { unwrapMany, unwrapOne } from './helpers.js';
import type { Budget } from '$lib/types/domain.js';

export const budgetService = {
	async list(): Promise<Budget[]> {
		const res = await api.get<{ data: { id: string; attributes: Record<string, unknown> }[] }>('/budgets');
		return unwrapMany<Budget>(res).map((b) => ({
			...b,
			spend_amount: '0',
			budget_amount: '0',
			limits: []
		}));
	},

	async create(data: { name: string; order?: number }): Promise<Budget> {
		const res = await api.post<{ data: { id: string; attributes: Record<string, unknown> } }>('/budgets', data);
		const b = unwrapOne<Budget>(res);
		return { ...b, spend_amount: '0', budget_amount: '0', limits: [] };
	},

	async update(id: string, data: { name?: string; active?: boolean }): Promise<void> {
		await api.put(`/budgets/${id}`, data);
	},

	async delete(id: string): Promise<void> {
		await api.delete(`/budgets/${id}`);
	}
};

import { api } from './client.js';
import { unwrapMany, unwrapOne } from './helpers.js';
import type { RecurringTransaction } from '$lib/types/domain.js';

export const recurringService = {
	async list(): Promise<RecurringTransaction[]> {
		const res = await api.get<{ data: { id: string; attributes: Record<string, unknown> }[] }>('/recurrences');
		return unwrapMany<RecurringTransaction>(res).map((r) => ({
			...r,
			type: 'withdrawal',
			amount: '0',
			currency_code: 'USD'
		}));
	},

	async create(data: {
		title: string;
		first_date: string;
		repeat_freq?: string;
		repeat_until?: string;
		transactions?: {
			type: string;
			description?: string;
			amount: string;
			source_id: string;
			destination_id: string;
			category_id?: string;
			piggy_bank_id?: string;
		}[];
	}): Promise<RecurringTransaction> {
		const res = await api.post<{ data: { id: string; attributes: Record<string, unknown> } }>('/recurrences', data);
		const r = unwrapOne<RecurringTransaction>(res);
		return { ...r, type: 'withdrawal', amount: '0', currency_code: 'USD' };
	}
};

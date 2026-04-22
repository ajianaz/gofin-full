import { api } from './client.js';
import { unwrapMany, unwrapOne } from './helpers.js';
import type { Bill } from '$lib/types/domain.js';

export const billService = {
	async list(): Promise<Bill[]> {
		const res = await api.get<{ data: { id: string; attributes: Record<string, unknown> }[] }>('/bills');
		return unwrapMany<Bill>(res).map((b) => ({
			...b,
			next_date: (b as any).date || '',
			currency_code: (b as any).currency_id || 'USD',
			currency_symbol: 'Rp'
		}));
	},

	async create(data: {
		name: string;
		amount_min?: string;
		amount_max?: string;
		date?: string;
		repeat_freq?: string;
	}): Promise<Bill> {
		const res = await api.post<{ data: { id: string; attributes: Record<string, unknown> } }>('/bills', {
			...data,
			amount_min: data.amount_min ? String(data.amount_min) : undefined,
			amount_max: data.amount_max ? String(data.amount_max) : undefined,
		});
		const b = unwrapOne<Bill>(res);
		return { ...b, next_date: (b as any).date || '', currency_code: 'USD', currency_symbol: 'Rp' };
	}
};

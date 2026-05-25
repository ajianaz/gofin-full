import { api } from './client.js';
import { unwrapMany, unwrapOne } from './helpers.js';
import type { Currency, ExchangeRate } from '$lib/types/domain.js';

export const currencyService = {
	async list(): Promise<Currency[]> {
		const res = await api.get<{ data: { id: string; attributes: Record<string, unknown> }[] }>('/currencies');
		return unwrapMany<Currency>(res).map((c) => ({
			...c,
			code: c.id,
			decimal_places: (c as any).decimal_places ?? 2,
			enabled: (c as any).enabled ?? true
		}));
	},

	async exchangeRates(): Promise<ExchangeRate[]> {
		const res = await api.get<{ data: { id: string; attributes: Record<string, unknown> }[] }>('/exchange-rates');
		return unwrapMany<ExchangeRate>(res).map((r) => ({
			...r,
			from_code: (r as any).from_currency_id ?? (r as any).from_code ?? '',
			to_code: (r as any).to_currency_id ?? (r as any).to_code ?? '',
			rate: parseFloat(String((r as any).rate ?? '0')),
			date: (r as any).date ?? ''
		}));
	},

	async createExchangeRate(data: { from_currency_id: string; to_currency_id: string; rate: string; date?: string }): Promise<ExchangeRate> {
		const res = await api.post<{ data: { id: string; attributes: Record<string, unknown> } }>('/exchange-rates', data);
		const r = unwrapOne<ExchangeRate>(res);
		return {
			...r,
			from_code: (r as any).from_currency_id ?? (r as any).from_code ?? '',
			to_code: (r as any).to_currency_id ?? (r as any).to_code ?? '',
			rate: parseFloat(String((r as any).rate ?? '0')),
			date: (r as any).date ?? ''
		};
	},

	async deleteExchangeRate(id: string): Promise<void> {
		await api.delete(`/exchange-rates/${id}`);
	}
};

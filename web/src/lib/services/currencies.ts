import { api } from './client.js';
import { unwrapMany } from './helpers.js';
import type { Currency } from '$lib/types/domain.js';

export const currencyService = {
	async list(): Promise<Currency[]> {
		const res = await api.get<{ data: { id: string; attributes: Record<string, unknown> }[] }>('/currencies');
		return unwrapMany<Currency>(res).map((c) => ({
			...c,
			code: c.id,
			decimal_places: (c as any).decimal_places ?? 2,
			enabled: (c as any).enabled ?? true
		}));
	}
};

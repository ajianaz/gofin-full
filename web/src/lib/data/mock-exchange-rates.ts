import type { ExchangeRate } from '$lib/types/index.js';

export const mockExchangeRates: ExchangeRate[] = [
	{ id: 'er1', from_code: 'USD', to_code: 'IDR', rate: 15850.0, date: '2026-04-15' },
	{ id: 'er2', from_code: 'EUR', to_code: 'IDR', rate: 17200.5, date: '2026-04-15' },
	{ id: 'er3', from_code: 'GBP', to_code: 'IDR', rate: 20100.0, date: '2026-04-15' },
	{ id: 'er4', from_code: 'JPY', to_code: 'IDR', rate: 105.2, date: '2026-04-15' },
	{ id: 'er5', from_code: 'SGD', to_code: 'IDR', rate: 11800.0, date: '2026-04-15' },
	{ id: 'er6', from_code: 'USD', to_code: 'EUR', rate: 0.92, date: '2026-04-15' },
	{ id: 'er7', from_code: 'USD', to_code: 'JPY', rate: 150.8, date: '2026-04-15' },
	{ id: 'er8', from_code: 'USD', to_code: 'SGD', rate: 1.34, date: '2026-04-15' }
];

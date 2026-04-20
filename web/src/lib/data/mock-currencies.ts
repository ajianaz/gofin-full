import type { Currency, ExchangeRate } from '$lib/types/index.js';

export const mockCurrencies: Currency[] = [
	{ id: 'c1', code: 'IDR', name: 'Rupiah Indonesia', symbol: 'Rp', decimal_places: 0, enabled: true },
	{ id: 'c2', code: 'USD', name: 'Dolar Amerika', symbol: '$', decimal_places: 2, enabled: true },
	{ id: 'c3', code: 'EUR', name: 'Euro', symbol: '€', decimal_places: 2, enabled: true },
	{ id: 'c4', code: 'SGD', name: 'Dolar Singapura', symbol: 'S$', decimal_places: 2, enabled: true },
	{ id: 'c5', code: 'GBP', name: 'Pound Sterling', symbol: '£', decimal_places: 2, enabled: true },
	{ id: 'c6', code: 'JPY', name: 'Yen Jepang', symbol: '¥', decimal_places: 0, enabled: false }
];

export const mockExchangeRates: ExchangeRate[] = [
	{ id: 'er1', from_code: 'USD', to_code: 'IDR', rate: 16150, date: '2026-04-15' },
	{ id: 'er2', from_code: 'EUR', to_code: 'IDR', rate: 17400, date: '2026-04-15' },
	{ id: 'er3', from_code: 'SGD', to_code: 'IDR', rate: 11900, date: '2026-04-15' },
	{ id: 'er4', from_code: 'GBP', to_code: 'IDR', rate: 20400, date: '2026-04-15' },
	{ id: 'er5', from_code: 'USD', to_code: 'EUR', rate: 0.92, date: '2026-04-15' }
];

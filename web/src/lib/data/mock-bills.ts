import type { Bill } from '$lib/types/index.js';

export const mockBills: Bill[] = [
	{ id: 'bl1', name: 'Listrik PLN', amount_min: '400000', amount_max: '500000', next_date: '2026-04-20', repeat_freq: 'monthly', active: true, currency_code: 'IDR', currency_symbol: 'Rp' },
	{ id: 'bl2', name: 'Internet Indihome', amount_min: '350000', amount_max: '350000', next_date: '2026-04-25', repeat_freq: 'monthly', active: true, currency_code: 'IDR', currency_symbol: 'Rp' },
	{ id: 'bl3', name: 'BPJS Kesehatan', amount_min: '150000', amount_max: '150000', next_date: '2026-04-05', repeat_freq: 'monthly', active: true, currency_code: 'IDR', currency_symbol: 'Rp' },
	{ id: 'bl4', name: 'Cicilan Motor', amount_min: '1500000', amount_max: '1500000', next_date: '2026-04-27', repeat_freq: 'monthly', active: true, currency_code: 'IDR', currency_symbol: 'Rp' },
	{ id: 'bl5', name: 'Spotify', amount_min: '54990', amount_max: '54990', next_date: '2026-05-06', repeat_freq: 'monthly', active: true, currency_code: 'IDR', currency_symbol: 'Rp' },
	{ id: 'bl6', name: 'Netflix', amount_min: '186000', amount_max: '186000', next_date: '2026-05-02', repeat_freq: 'monthly', active: false, currency_code: 'IDR', currency_symbol: 'Rp' }
];

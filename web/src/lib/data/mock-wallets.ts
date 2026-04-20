import type { Account } from '$lib/types/index.js';

export const mockWallets: Account[] = [
	{ id: '1', name: 'BCA', type: 'asset', active: true, balance: '12500000', currency_code: 'IDR', currency_symbol: 'Rp', currency_decimal_places: 0 },
	{ id: '2', name: 'Mandiri', type: 'asset', active: true, balance: '8500000', currency_code: 'IDR', currency_symbol: 'Rp', currency_decimal_places: 0 },
	{ id: '3', name: 'GoPay', type: 'asset', active: true, balance: '1500000', currency_code: 'IDR', currency_symbol: 'Rp', currency_decimal_places: 0 },
	{ id: '4', name: 'OVO', type: 'asset', active: true, balance: '750000', currency_code: 'IDR', currency_symbol: 'Rp', currency_decimal_places: 0 },
	{ id: '5', name: 'Kartu Kredit BCA', type: 'liability', active: true, balance: '-3200000', currency_code: 'IDR', currency_symbol: 'Rp', currency_decimal_places: 0 },
	{ id: '6', name: 'Dana', type: 'asset', active: true, balance: '2100000', currency_code: 'IDR', currency_symbol: 'Rp', currency_decimal_places: 0 },
	{ id: '7', name: 'Kas', type: 'cash', active: true, balance: '500000', currency_code: 'IDR', currency_symbol: 'Rp', currency_decimal_places: 0 }
];

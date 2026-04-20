import type { Category } from '$lib/types/index.js';

export const mockCategories: Category[] = [
	{ id: '1', name: 'Makanan', type: 'expense', transaction_count: 45 },
	{ id: '2', name: 'Gaji', type: 'income', transaction_count: 12 },
	{ id: '3', name: 'Utilitas', type: 'expense', transaction_count: 24 },
	{ id: '4', name: 'Belanja', type: 'expense', transaction_count: 30 },
	{ id: '5', name: 'Freelance', type: 'income', transaction_count: 5 },
	{ id: '6', name: 'Telepon', type: 'expense', transaction_count: 15 },
	{ id: '7', name: 'Hiburan', type: 'expense', transaction_count: 8 },
	{ id: '8', name: 'Transportasi', type: 'expense', transaction_count: 20 },
	{ id: '9', name: 'Kesehatan', type: 'expense', transaction_count: 12 },
	{ id: '10', name: 'Cicilan', type: 'expense', transaction_count: 6 },
	{ id: '11', name: 'Investasi', type: 'income', transaction_count: 4 },
	{ id: '12', name: 'Pendidikan', type: 'expense', transaction_count: 3 }
];

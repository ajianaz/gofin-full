import type { PiggyBank } from '$lib/types/index.js';

export const mockPiggyBanks: PiggyBank[] = [
	{ id: 'pb1', name: 'Liburan Bali', account_id: '2', account_name: 'Mandiri', target_amount: '15000000', current_amount: '8500000', start_date: '2026-01-01', target_date: '2026-12-31', status: 'active' },
	{ id: 'pb2', name: 'Laptop Baru', account_id: '1', account_name: 'BCA', target_amount: '12000000', current_amount: '9500000', start_date: '2026-02-01', target_date: '2026-08-01', status: 'active' },
	{ id: 'pb3', name: 'Dana Darurat', account_id: '1', account_name: 'BCA', target_amount: '5000000', current_amount: '3200000', start_date: '2026-03-01', status: 'active' },
	{ id: 'pb4', name: 'Kursus Online', account_id: '6', account_name: 'Dana', target_amount: '3000000', current_amount: '1500000', start_date: '2026-01-15', target_date: '2026-06-30', status: 'active' },
	{ id: 'pb5', name: 'Hadiah Ulang Tahun', account_id: '1', account_name: 'BCA', target_amount: '2000000', current_amount: '2000000', start_date: '2025-10-01', target_date: '2026-05-01', status: 'completed' }
];

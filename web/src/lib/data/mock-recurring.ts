import type { RecurringTransaction } from '$lib/types/index.js';

export const mockRecurring: RecurringTransaction[] = [
	{ id: 'rt1', title: 'Gaji Bulanan', type: 'deposit', first_date: '2026-01-14', repeat_freq: 'monthly', active: true, amount: '8500000', currency_code: 'IDR', description: 'Gaji dari kantor', destination_account_id: '1', destination_account_name: 'BCA' },
	{ id: 'rt2', title: 'Bayar Kos', type: 'withdrawal', first_date: '2026-01-01', repeat_freq: 'monthly', active: true, amount: '1500000', currency_code: 'IDR', description: 'Bayar kos bulanan', source_account_id: '1', source_account_name: 'BCA' },
	{ id: 'rt3', title: 'BPJS Kesehatan', type: 'withdrawal', first_date: '2026-01-05', repeat_freq: 'monthly', active: true, amount: '150000', currency_code: 'IDR', description: 'Iuran BPJS', source_account_id: '2', source_account_name: 'Mandiri' },
	{ id: 'rt4', title: 'Internet Indihome', type: 'withdrawal', first_date: '2026-01-26', repeat_freq: 'monthly', active: true, amount: '350000', currency_code: 'IDR', description: 'Tagihan internet', source_account_id: '2', source_account_name: 'Mandiri' },
	{ id: 'rt5', title: 'Netflix', type: 'withdrawal', first_date: '2025-12-02', repeat_freq: 'monthly', active: false, amount: '186000', currency_code: 'IDR', description: 'Langganan Netflix', source_account_id: '5', source_account_name: 'Kartu Kredit BCA' }
];

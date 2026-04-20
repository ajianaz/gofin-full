import type { Preference } from '$lib/types/index.js';

export const mockPreferences: Preference[] = [
	{ name: 'language', value: 'id', description: 'Bahasa tampilan', type: 'select', options: ['id', 'en'] },
	{ name: 'currency', value: 'IDR', description: 'Mata uang utama', type: 'select', options: ['IDR', 'USD', 'EUR'] },
	{ name: 'date_format', value: 'DD MMM YYYY', description: 'Format tanggal', type: 'select', options: ['DD MMM YYYY', 'YYYY-MM-DD', 'MM/DD/YYYY'] },
	{ name: 'group_style', value: 'default', description: 'Gaya tampilan grup', type: 'select', options: ['default', 'compact'] },
	{ name: 'budget_indicator', value: true, description: 'Indikator anggaran di sidebar', type: 'boolean' },
	{ name: 'show_news', value: false, description: 'Tampilkan berita di dashboard', type: 'boolean' },
	{ name: 'fiscal_year_start', value: '01', description: 'Bulan awal tahun fiskal', type: 'select', options: ['01', '04', '07', '10'] },
	{ name: 'transaction_count_per_page', value: '25', description: 'Jumlah transaksi per halaman', type: 'select', options: ['10', '15', '25', '50'] },
	{ name: 'two_factor_enabled', value: false, description: 'Aktifkan autentikasi dua faktor', type: 'boolean' },
	{ name: 'email_digest', value: 'weekly', description: 'Frekuensi email ringkasan', type: 'select', options: ['daily', 'weekly', 'monthly'] }
];

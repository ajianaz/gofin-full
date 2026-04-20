import type { NotificationSetting } from '$lib/types/index.js';

export const mockNotifications: NotificationSetting[] = [
	{ id: 'n1', title: 'Transaksi Baru', description: 'Notifikasi saat ada transaksi baru', enabled: true },
	{ id: 'n2', title: 'Tagihan Mendatang', description: 'Pengingat tagihan 3 hari sebelum jatuh tempo', enabled: true },
	{ id: 'n3', title: 'Anggaran Melebihi', description: 'Peringatan jika anggaran melebihi 80%', enabled: true },
	{ id: 'n4', title: 'Tabungan Tercapai', description: 'Notifikasi saat target tabungan tercapai', enabled: true },
	{ id: 'n5', title: 'Bill Gagal', description: 'Peringatan saat pembayaran tagihan gagal', enabled: true },
	{ id: 'n6', title: 'Login Baru', description: 'Notifikasi login dari perangkat baru', enabled: false },
	{ id: 'n7', title: 'Laporan Mingguan', description: 'Ringkasan keuangan setiap minggu', enabled: true },
	{ id: 'n8', title: 'Recurring Gagal', description: 'Peringatan saat transaksi berulang gagal', enabled: false },
	{ id: 'n9', title: 'Update Sistem', description: 'Notifikasi update dan maintenance', enabled: true },
	{ id: 'n10', title: 'Pemberitahuan API', description: 'Notifikasi aktivitas API key', enabled: true }
];

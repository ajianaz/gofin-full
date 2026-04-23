import { api } from '$lib/services/client.js';

export const exportService = {
	async downloadCSV(startDate?: string, endDate?: string, walletId?: string): Promise<void> {
		const params = new URLSearchParams();
		if (startDate) params.set('start', startDate);
		if (endDate) params.set('end', endDate);
		if (walletId) params.set('wallet_id', walletId);
		const qs = params.toString();
		const url = `/export/csv${qs ? '?' + qs : ''}`;

		const response = await fetch(`/api/v1${url}`, {
			headers: {
				'Content-Type': 'application/json',
				...(localStorage.getItem('access_token') ? { Authorization: `Bearer ${localStorage.getItem('access_token')}` } : {})
			}
		});

		if (!response.ok) throw new Error(`Export failed: ${response.statusText}`);

		const blob = await response.blob();
		const a = document.createElement('a');
		a.href = URL.createObjectURL(blob);
		a.download = 'transactions.csv';
		a.click();
		URL.revokeObjectURL(a.href);
	},

	async downloadOFX(startDate?: string, endDate?: string, walletId?: string): Promise<void> {
		const params = new URLSearchParams();
		if (startDate) params.set('start', startDate);
		if (endDate) params.set('end', endDate);
		if (walletId) params.set('wallet_id', walletId);
		const qs = params.toString();
		const url = `/export/ofx${qs ? '?' + qs : ''}`;

		const response = await fetch(`/api/v1${url}`, {
			headers: {
				'Content-Type': 'application/json',
				...(localStorage.getItem('access_token') ? { Authorization: `Bearer ${localStorage.getItem('access_token')}` } : {})
			}
		});

		if (!response.ok) throw new Error(`Export failed: ${response.statusText}`);

		const blob = await response.blob();
		const a = document.createElement('a');
		a.href = URL.createObjectURL(blob);
		a.download = 'transactions.ofx';
		a.click();
		URL.revokeObjectURL(a.href);
	}
};

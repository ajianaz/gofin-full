export const exportService = {
	async downloadCSV(startDate?: string, endDate?: string, walletId?: string): Promise<void> {
		const params = new URLSearchParams();
		if (startDate) params.set('start', startDate);
		if (endDate) params.set('end', endDate);
		if (walletId) params.set('wallet_id', walletId);
		const qs = params.toString();
		const url = `/api/v1/export/csv${qs ? '?' + qs : ''}`;

		const token = localStorage.getItem('access_token');
		const response = await fetch(url, {
			headers: {
				...(token ? { Authorization: `Bearer ${token}` } : {})
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
		const url = `/api/v1/export/ofx${qs ? '?' + qs : ''}`;

		const token = localStorage.getItem('access_token');
		const response = await fetch(url, {
			headers: {
				...(token ? { Authorization: `Bearer ${token}` } : {})
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

import { api } from '$lib/services/client.js';

async function downloadExport(format: string, filename: string, startDate?: string, endDate?: string, walletId?: string): Promise<void> {
	const allowedFormats = ['csv', 'ofx'];
	if (!allowedFormats.includes(format)) {
		throw new Error(`Unsupported export format: ${format}`);
	}

	const params = new URLSearchParams();
	if (startDate) params.set('start', startDate);
	if (endDate) params.set('end', endDate);
	if (walletId) params.set('wallet_id', walletId);
	const qs = params.toString();
	const path = `/export/${format}${qs ? '?' + qs : ''}`;

	// Use the central API client which handles auth headers and token refresh
	const response = await api.blob(path);
	if (!response.ok) throw new Error(`Export failed: ${response.statusText}`);

	const blob = await response.blob();
	const a = document.createElement('a');
	a.href = URL.createObjectURL(blob);
	a.download = filename;
	a.click();
	URL.revokeObjectURL(a.href);
}

export const exportService = {
	downloadCSV: (startDate?: string, endDate?: string, walletId?: string) =>
		downloadExport('csv', 'transactions.csv', startDate, endDate, walletId),
	downloadOFX: (startDate?: string, endDate?: string, walletId?: string) =>
		downloadExport('ofx', 'transactions.ofx', startDate, endDate, walletId)
};

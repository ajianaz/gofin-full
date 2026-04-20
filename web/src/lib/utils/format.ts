import { localeStore } from '$lib/stores/i18n.svelte.js';

export function formatCurrency(
	amount: string,
	symbol: string = 'Rp',
	decimalPlaces: number = 0
): string {
	const num = parseFloat(amount);
	if (isNaN(num)) return `${symbol}0`;
	return `${symbol}${Math.abs(num).toLocaleString(localeStore.localeCode, {
		minimumFractionDigits: decimalPlaces,
		maximumFractionDigits: decimalPlaces
	})}`;
}

export function formatAmount(amount: string): { text: string; color: string } {
	const num = parseFloat(amount);
	if (isNaN(num)) return { text: 'Rp0', color: 'text-foreground' };
	const isNegative = num < 0;
	return {
		text: `${isNegative ? '-' : '+'}Rp ${Math.abs(num).toLocaleString(localeStore.localeCode)}`,
		color: isNegative ? 'text-red-600' : 'text-green-600'
	};
}

export function formatDate(date: string): string {
	const d = new Date(date);
	if (isNaN(d.getTime())) return date;
	return d.toLocaleDateString(localeStore.localeCode, {
		day: 'numeric',
		month: 'short',
		year: 'numeric'
	});
}

export function formatPercentage(value: number): string {
	return `${value.toFixed(1)}%`;
}

export function formatNumber(value: number): string {
	return value.toLocaleString(localeStore.localeCode);
}

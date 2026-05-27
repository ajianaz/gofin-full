import { localeStore } from '$lib/stores/i18n.svelte.js';

const DEFAULT_SYMBOLS: Record<string, string> = {
	en: '$',
	id: 'Rp'
};

export function getDefaultSymbol(): string {
	return DEFAULT_SYMBOLS[localeStore.locale] ?? '$';
}

/**
 * Format a numeric amount string with currency symbol.
 * Always prefer passing the wallet/transaction's own symbol.
 * Falls back to locale-based default only when no symbol is available.
 */
export function formatCurrency(
	amount: string,
	symbol?: string,
	decimalPlaces: number = 2
): string {
	const sym = symbol || getDefaultSymbol();
	const num = parseFloat(amount);
	if (isNaN(num)) return `${sym}0`;
	return `${sym}${Math.abs(num).toLocaleString(localeStore.localeCode, {
		minimumFractionDigits: decimalPlaces,
		maximumFractionDigits: decimalPlaces
	})}`;
}

/**
 * Format a balance from a wallet, using its own currency symbol.
 */
export function formatBalance(balance: string, currencySymbol?: string, decimalPlaces: number = 2): string {
	const sym = currencySymbol || getDefaultSymbol();
	const num = parseFloat(balance);
	if (isNaN(num)) return `${sym}0`;
	const isNeg = num < 0;
	return `${isNeg ? '-' : ''}${sym}${Math.abs(num).toLocaleString(localeStore.localeCode, {
		minimumFractionDigits: decimalPlaces,
		maximumFractionDigits: decimalPlaces
	})}`;
}

export function formatAmount(amount: string, symbol?: string): { text: string; color: string } {
	const sym = symbol || getDefaultSymbol();
	const num = parseFloat(amount);
	if (isNaN(num)) return { text: `${sym}0`, color: 'text-foreground' };
	const isNegative = num < 0;
	return {
		text: `${isNegative ? '-' : '+'}${sym} ${Math.abs(num).toLocaleString(localeStore.localeCode)}`,
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

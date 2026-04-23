import { api } from './client.js';
import { unwrapMany } from './helpers.js';

interface CategorySpending {
	category_id: string;
	category_name: string;
	total: number;
	count: number;
}

interface PeriodSpending {
	period: string;
	income: number;
	expense: number;
}

interface NetWorthSummary {
	total_income: number;
	total_expense: number;
	net_income: number;
	transaction_count: number;
}

export const reportService = {
	async spendingByCategory(start?: string, end?: string): Promise<CategorySpending[]> {
		const params = new URLSearchParams();
		if (start) params.set('start', start);
		if (end) params.set('end', end);
		const qs = params.toString();
		const res = await api.get<{ data: { id: string; attributes: Record<string, unknown> }[] }>(
			`/analytics/spending-by-category${qs ? '?' + qs : ''}`
		);
		return unwrapMany<CategorySpending>(res).map((r) => ({
			category_id: String((r as any).category_id ?? ''),
			category_name: String((r as any).category_name ?? ''),
			total: parseFloat(String((r as any).total ?? '0')),
			count: parseInt(String((r as any).count ?? '0'), 10)
		}));
	},

	async spendingByPeriod(start?: string, end?: string): Promise<PeriodSpending[]> {
		const params = new URLSearchParams();
		if (start) params.set('start', start);
		if (end) params.set('end', end);
		const qs = params.toString();
		const res = await api.get<{ data: { id: string; attributes: Record<string, unknown> }[] }>(
			`/analytics/spending-by-period${qs ? '?' + qs : ''}`
		);
		return unwrapMany<PeriodSpending>(res).map((r) => ({
			period: String((r as any).period ?? ''),
			income: parseFloat(String((r as any).income ?? '0')),
			expense: parseFloat(String((r as any).expense ?? '0'))
		}));
	},

	async netWorth(start?: string, end?: string): Promise<NetWorthSummary> {
		const params = new URLSearchParams();
		if (start) params.set('start', start);
		if (end) params.set('end', end);
		const qs = params.toString();
		const res = await api.get<{ data: { type: string; attributes: Record<string, unknown> } }>(
			`/analytics/net-worth${qs ? '?' + qs : ''}`
		);
		const attrs = res.data.attributes;
		return {
			total_income: parseFloat(String(attrs.total_income ?? '0')),
			total_expense: parseFloat(String(attrs.total_expense ?? '0')),
			net_income: parseFloat(String(attrs.net_income ?? '0')),
			transaction_count: parseInt(String(attrs.transaction_count ?? '0'), 10)
		};
	}
};

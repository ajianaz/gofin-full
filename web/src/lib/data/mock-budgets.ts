import type { Budget, BudgetLimit } from '$lib/types/index.js';

const makanLimit: BudgetLimit = { id: 'bl1', budget_id: 'b1', start_date: '2026-04-01', end_date: '2026-04-30', amount: '2000000', spend: '1250000', category: '1', category_name: 'Makanan' };
const belanjaLimit: BudgetLimit = { id: 'bl2', budget_id: 'b1', start_date: '2026-04-01', end_date: '2026-04-30', amount: '1500000', spend: '630000', category: '4', category_name: 'Belanja' };
const utilitasLimit: BudgetLimit = { id: 'bl3', budget_id: 'b1', start_date: '2026-04-01', end_date: '2026-04-30', amount: '1000000', spend: '800000', category: '3', category_name: 'Utilitas' };
const transportLimit: BudgetLimit = { id: 'bl4', budget_id: 'b1', start_date: '2026-04-01', end_date: '2026-04-30', amount: '500000', spend: '115000', category: '8', category_name: 'Transportasi' };

export const mockBudgets: Budget[] = [
	{ id: 'b1', name: 'Anggaran April 2026', active: true, order: 1, auto_budget_type: 'none', spend_amount: '2795000', budget_amount: '5000000', limits: [makanLimit, belanjaLimit, utilitasLimit, transportLimit] },
	{ id: 'b2', name: 'Anggaran Maret 2026', active: true, order: 2, auto_budget_type: 'none', spend_amount: '4250000', budget_amount: '5000000', limits: [] },
	{ id: 'b3', name: 'Tabungan Liburan', active: true, order: 3, auto_budget_type: 'none', spend_amount: '3500000', budget_amount: '5000000', limits: [] },
	{ id: 'b4', name: 'Kebutuhan Pokok', active: true, order: 4, auto_budget_type: 'fixed', auto_budget_period: 'monthly', spend_amount: '2100000', budget_amount: '3000000', limits: [] },
	{ id: 'b5', name: 'Dana Darurat', active: true, order: 5, spend_amount: '0', budget_amount: '2000000', limits: [] }
];

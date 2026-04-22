import { api } from './client.js';
import { unwrapMany, unwrapOne } from './helpers.js';
import type { PiggyBank } from '$lib/types/domain.js';

export const piggyBankService = {
	async list(walletId: string): Promise<PiggyBank[]> {
		const res = await api.get<{ data: { id: string; attributes: Record<string, unknown> }[] }>(`/wallets/${walletId}/piggy_banks`);
		return unwrapMany<PiggyBank>(res).map((p) => ({
			...p,
			account_id: (p as any).wallet_id || walletId,
			account_name: '',
			current_amount: '0',
			status: 'active'
		}));
	},

	async create(data: {
		wallet_id: string;
		name: string;
		target_amount?: string;
	}): Promise<PiggyBank> {
		const res = await api.post<{ data: { id: string; attributes: Record<string, unknown> } }>(`/wallets/${data.wallet_id}/piggy_banks`, data);
		const p = unwrapOne<PiggyBank>(res);
		return { ...p, account_id: data.wallet_id, account_name: '', current_amount: '0', status: 'active' };
	}
};

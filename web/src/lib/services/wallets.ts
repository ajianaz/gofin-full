import { api } from './client.js';
import { unwrapMany, unwrapOne } from './helpers.js';
import type { Account } from '$lib/types/domain.js';

interface WalletTypeRaw { id: string; attributes: { type: string } }

function mapWalletAttrs(w: Record<string, unknown>): Partial<Account> {
	const a = (w as any) || {};
	return {
		type: a.wallet_type || a.type || 'asset',
		balance: a.virtual_balance || '0',
		currency_code: a.currency_code || '',
		currency_symbol: a.currency_symbol || '',
		currency_decimal_places: a.currency_decimal_places ?? 2
	};
}

export const walletService = {
	async list(): Promise<Account[]> {
		const res = await api.get<{ data: { id: string; attributes: Record<string, unknown> }[] }>('/wallets');
		return unwrapMany<Account>(res).map((w) => ({
			...w,
			...mapWalletAttrs(w as unknown as Record<string, unknown>)
		}));
	},

	async create(data: { name: string; wallet_type?: string; currency_id?: string; active?: boolean }): Promise<Account> {
		const res = await api.post<{ data: { id: string; attributes: Record<string, unknown> } }>('/wallets', data);
		const w = unwrapOne<Account>(res);
		return {
			...w,
			...mapWalletAttrs(w as unknown as Record<string, unknown>)
		};
	},

	async update(id: string, data: { name?: string; active?: boolean; currency_id?: string }): Promise<void> {
		await api.put(`/wallets/${id}`, data);
	},

	async delete(id: string): Promise<void> {
		await api.delete(`/wallets/${id}`);
	},

	async types(): Promise<WalletTypeRaw[]> {
		const res = await api.get<{ data: WalletTypeRaw[] }>('/wallet-types');
		return res.data;
	}
};

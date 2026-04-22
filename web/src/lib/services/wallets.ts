import { api } from './client.js';
import { unwrapMany, unwrapOne } from './helpers.js';
import type { Account } from '$lib/types/domain.js';

interface WalletTypeRaw { id: string; attributes: { type: string } }

export const walletService = {
	async list(): Promise<Account[]> {
		const res = await api.get<{ data: { id: string; attributes: Record<string, unknown> }[] }>('/wallets');
		return unwrapMany<Account>(res).map((w) => ({
			...w,
			type: (w as any).wallet_type || w.type || 'asset',
			balance: (w as any).virtual_balance || '0',
			currency_code: (w as any).currency_id || 'USD',
			currency_symbol: 'Rp',
			currency_decimal_places: 0
		}));
	},

	async create(data: { name: string; wallet_type?: string; active?: boolean }): Promise<Account> {
		const res = await api.post<{ data: { id: string; attributes: Record<string, unknown> } }>('/wallets', data);
		const w = unwrapOne<Account>(res);
		return {
			...w,
			type: (w as any).wallet_type || w.type || 'asset',
			balance: '0',
			currency_code: 'USD',
			currency_symbol: 'Rp',
			currency_decimal_places: 0
		};
	},

	async types(): Promise<WalletTypeRaw[]> {
		const res = await api.get<{ data: WalletTypeRaw[] }>('/wallet-types');
		return res.data;
	}
};

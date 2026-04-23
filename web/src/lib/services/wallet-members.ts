import { api } from './client.js';
import { unwrapMany } from './helpers.js';
import type { WalletMember } from '$lib/types/domain.js';

export const walletMemberService = {
	async list(walletId: string): Promise<WalletMember[]> {
			const res = await api.get<{ data: { id: string; attributes: Record<string, unknown> }[] | null }>(`/wallets/${walletId}/members`);
		return unwrapMany<WalletMember>(res).map((m) => ({
			id: m.id,
			wallet_id: (m as any).wallet_id || walletId,
			user_id: (m as any).user_id || '',
			role: ((m as any).role as WalletMember['role']) || 'viewer'
		}));
	},

	async add(walletId: string, userId: string, role: string): Promise<void> {
		await api.post(`/wallets/${walletId}/members`, { user_id: userId, role });
	},

	async updateRole(walletId: string, id: string, role: string): Promise<void> {
		await api.put(`/wallets/${walletId}/members/${id}`, { role });
	},

	async remove(walletId: string, userId: string): Promise<void> {
		await api.delete(`/wallets/${walletId}/members/${userId}`);
	}
};

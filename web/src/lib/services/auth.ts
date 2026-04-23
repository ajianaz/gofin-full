import { api } from './client.js';
import type { LoginRequest, RegisterRequest, TokenResponse, User } from '$lib/types/index.js';
import { unwrapOne } from './helpers.js';

export const authService = {
	async login(data: LoginRequest): Promise<TokenResponse> {
		return api.post<TokenResponse>('/auth/login', data);
	},

	async register(data: RegisterRequest): Promise<TokenResponse> {
		return api.post<TokenResponse>('/auth/register', data);
	},

	async logout(refreshToken?: string): Promise<{ message: string }> {
		return api.post<{ message: string }>('/auth/logout', { refresh_token: refreshToken });
	},

	async refresh(refreshToken: string): Promise<TokenResponse> {
		return api.post<TokenResponse>('/auth/refresh', { refresh_token: refreshToken });
	},

	async getMe(): Promise<User> {
		const res = await api.get<{ data: { id: string; attributes: { email: string; created_at?: string } } }>('/users/me');
		const raw = unwrapOne<{ email: string; created_at?: string }>(res);
		const email = raw.email || '';
		const name = email.split('@')[0] || '';
		return {
			id: raw.id,
			email: raw.email,
			name,
			created_at: raw.created_at || new Date().toISOString()
		};
	},

	async getProvider(): Promise<{ provider: string }> {
		return api.get<{ provider: string }>('/auth/provider');
	},

	async updateProfile(data: { name?: string; email?: string }): Promise<User> {
		await api.put('/users/me', data);
		return this.getMe();
	},

	async changePassword(data: { current_password: string; new_password: string }): Promise<void> {
		await api.post('/users/me/password', data);
	}
};

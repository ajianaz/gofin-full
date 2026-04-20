import type { LoginRequest, RegisterRequest, TokenResponse, User, ApiError } from '$lib/types/index.js';

const MOCK_USERS: Record<string, { password: string; user: User }> = {
	'admin@gofin.id': {
		password: 'admin123',
		user: { id: 'u1', email: 'admin@gofin.id', name: 'Admin Gofin', created_at: '2025-01-01' }
	},
	'budi@gofin.id': {
		password: 'budi123',
		user: { id: 'u2', email: 'budi@gofin.id', name: 'Budi Santoso', created_at: '2025-03-15' }
	},
	'sari@gofin.id': {
		password: 'sari123',
		user: { id: 'u3', email: 'sari@gofin.id', name: 'Sari Dewi', created_at: '2025-06-20' }
	}
};

function mockTokens(): TokenResponse {
	return {
		access_token: 'mock_access_token_' + Date.now(),
		refresh_token: 'mock_refresh_token_' + Date.now(),
		expires_in: 3600,
		token_type: 'Bearer'
	};
}

function delay(ms = 300) {
	return new Promise((r) => setTimeout(r, ms));
}

export const authService = {
	async login(data: LoginRequest): Promise<TokenResponse> {
		await delay();
		const entry = MOCK_USERS[data.email];
		if (!entry || entry.password !== data.password) {
			const err: ApiError = { status: 401, detail: 'Email atau password salah.' };
			throw err;
		}
		return mockTokens();
	},

	async register(data: RegisterRequest): Promise<TokenResponse> {
		await delay();
		if (MOCK_USERS[data.email]) {
			const err: ApiError = { status: 409, detail: 'Email sudah terdaftar.' };
			throw err;
		}
		return mockTokens();
	},

	async logout(): Promise<{ message: string }> {
		await delay(100);
		return { message: 'ok' };
	},

	async refresh(_refreshToken: string): Promise<TokenResponse> {
		await delay();
		return mockTokens();
	},

	async getMe(): Promise<User> {
		await delay(100);
		return { id: 'u1', email: 'admin@gofin.id', name: 'Admin Gofin', created_at: '2025-01-01' };
	},

	async getProvider(): Promise<{ provider: string }> {
		return { provider: 'local' };
	}
};

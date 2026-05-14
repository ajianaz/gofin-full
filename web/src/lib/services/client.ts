import type { ApiError, TokenResponse } from '$lib/types/index.js';

const API_BASE = '/api/v1';

let isRefreshing = false;
let refreshPromise: Promise<TokenResponse | null> | null = null;

async function refreshAccessToken(): Promise<TokenResponse | null> {
	const refreshToken = localStorage.getItem('refresh_token');
	if (!refreshToken) return null;

	try {
		const response = await fetch(`${API_BASE}/auth/refresh`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ refresh_token: refreshToken })
		});

		if (!response.ok) {
			// Refresh failed — clear tokens
			localStorage.removeItem('access_token');
			localStorage.removeItem('refresh_token');
			return null;
		}

		const tokens: TokenResponse = await response.json();
		localStorage.setItem('access_token', tokens.access_token);
		localStorage.setItem('refresh_token', tokens.refresh_token);
		return tokens;
	} catch {
		return null;
	}
}

async function getRefreshedToken(): Promise<TokenResponse | null> {
	// If a refresh is already in flight, reuse the same promise
	if (isRefreshing && refreshPromise) {
		return refreshPromise;
	}

	isRefreshing = true;
	refreshPromise = refreshAccessToken().finally(() => {
		isRefreshing = false;
		refreshPromise = null;
	});

	return refreshPromise;
}

async function request<T>(
	path: string,
	options: RequestInit = {}
): Promise<T> {
	const url = `${API_BASE}${path}`;

	const headers: Record<string, string> = {
		'Content-Type': 'application/json',
		...((options.headers as Record<string, string>) || {})
	};

	const token = localStorage.getItem('access_token');
	if (token) {
		headers['Authorization'] = `Bearer ${token}`;
	}

	const response = await fetch(url, {
		...options,
		headers
	});

	if (response.status === 401 && path !== '/auth/refresh' && path !== '/auth/login' && path !== '/auth/register') {
		const newTokens = await getRefreshedToken();

		if (newTokens) {
			// Retry the original request with the new token
			const retryHeaders: Record<string, string> = {
				'Content-Type': 'application/json',
				...((options.headers as Record<string, string>) || {})
			};
			retryHeaders['Authorization'] = `Bearer ${newTokens.access_token}`;

			const retryResponse = await fetch(url, {
				...options,
				headers: retryHeaders
			});

			if (!retryResponse.ok) {
				let error: ApiError = { status: retryResponse.status };
				try {
					error = await retryResponse.json();
				} catch {
					error.detail = retryResponse.statusText;
				}
				throw error;
			}

			if (retryResponse.status === 204) return undefined as T;
			return retryResponse.json();
		}

		// Refresh failed — clear tokens and redirect to login
		localStorage.removeItem('access_token');
		localStorage.removeItem('refresh_token');
		if (typeof window !== 'undefined') {
			window.location.href = '/login';
		}

		let error: ApiError = { status: 401 };
		try {
			error = await response.json();
		} catch {
			error.detail = 'Session expired';
		}
		throw error;
	}

	if (!response.ok) {
		let error: ApiError = { status: response.status };
		try {
			error = await response.json();
		} catch {
			error.detail = response.statusText;
		}
		throw error;
	}

	if (response.status === 204) return undefined as T;
	return response.json();
}

export const api = {
	get<T>(path: string, options?: RequestInit): Promise<T> {
		return request<T>(path, { ...options, method: 'GET' });
	},

	post<T>(path: string, body?: unknown, options?: RequestInit): Promise<T> {
		return request<T>(path, {
			...options,
			method: 'POST',
			body: body ? JSON.stringify(body) : undefined
		});
	},

	put<T>(path: string, body?: unknown, options?: RequestInit): Promise<T> {
		return request<T>(path, {
			...options,
			method: 'PUT',
			body: body ? JSON.stringify(body) : undefined
		});
	},

	delete<T>(path: string, options?: RequestInit): Promise<T> {
		return request<T>(path, { ...options, method: 'DELETE' });
	}
};

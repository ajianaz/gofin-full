import type { ApiError, TokenResponse } from '$lib/types/index.js';
import { authStore } from '$lib/stores/auth.svelte.js';
import { handleApiError } from '$lib/stores/toast.js';

const API_BASE = '/api/v1';

let isRefreshing = false;
let refreshPromise: Promise<TokenResponse | null> | null = null;

async function refreshAccessToken(): Promise<TokenResponse | null> {
	// Try in-memory refresh token first (available right after login/register).
	// If null (page reload), send request anyway — server reads httpOnly cookie.
	const refreshToken = authStore.refreshToken;

	const body = refreshToken
		? JSON.stringify({ refresh_token: refreshToken })
		: '{}';

	try {
		const response = await fetch(`${API_BASE}/auth/refresh`, {
			method: 'POST',
			credentials: 'same-origin', // ensure httpOnly cookies are sent
			headers: { 'Content-Type': 'application/json' },
			body
		});

		if (!response.ok) {
			// Refresh failed — clear reactive state
			authStore.clearTokens();
			return null;
		}

		const tokens: TokenResponse = await response.json();
		// Sync with reactive auth store (cookie persistence handled by server)
		authStore.setTokens(tokens);
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

	// Set Authorization header from in-memory token for backward compatibility.
	// The httpOnly cookie is always sent by the browser and is the primary
	// auth mechanism. This header ensures existing API key / header-based
	// clients continue to work during the migration.
	const token = authStore.accessToken;
	if (token) {
		headers['Authorization'] = `Bearer ${token}`;
	}

	const response = await fetch(url, {
		...options,
		credentials: 'same-origin', // ensure httpOnly cookies are sent
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
				credentials: 'same-origin',
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
		authStore.clearTokens();
		if (typeof window !== 'undefined') {
			window.location.href = '/login';
		}

		let error: ApiError = { status: 401 };
		try {
			error = await response.json();
		} catch {
			error.detail = 'Session expired';
		}
		handleApiError(error);
		throw error;
	}

	if (!response.ok) {
		let error: ApiError = { status: response.status };
		try {
			error = await response.json();
		} catch {
			error.detail = response.statusText;
		}
		handleApiError(error);
		throw error;
	}

	if (response.status === 204) return undefined as T;
	return response.json();
}

async function requestBlob(
	path: string,
	options: RequestInit = {}
): Promise<Response> {
	const url = `${API_BASE}${path}`;

	const headers: Record<string, string> = {
		'Content-Type': 'application/json',
		...((options.headers as Record<string, string>) || {})
	};

	const token = authStore.accessToken;
	if (token) {
		headers['Authorization'] = `Bearer ${token}`;
	}

	const response = await fetch(url, {
		...options,
		credentials: 'same-origin',
		headers
	});

	if (response.status === 401 && path !== '/auth/refresh' && path !== '/auth/login' && path !== '/auth/register') {
		const newTokens = await getRefreshedToken();

		if (newTokens) {
			const retryHeaders: Record<string, string> = {
				'Content-Type': 'application/json',
				...((options.headers as Record<string, string>) || {})
			};
			retryHeaders['Authorization'] = `Bearer ${newTokens.access_token}`;

			return fetch(url, {
				...options,
				credentials: 'same-origin',
				headers: retryHeaders
			});
		}

		// Refresh failed — clear tokens and redirect to login
		authStore.clearTokens();
		if (typeof window !== 'undefined') {
			window.location.href = '/login';
		}
	}

	return response;
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
	},

	/** Fetch a binary endpoint; returns the raw Response for blob/stream handling.
	 *  Token refresh is handled transparently. */
	blob(path: string, options?: RequestInit): Promise<Response> {
		return requestBlob(path, { ...options, method: 'GET' });
	}
};

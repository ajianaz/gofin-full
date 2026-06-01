import type { ApiError, TokenResponse } from '$lib/types/index.js';
import { authStore } from '$lib/stores/auth.svelte.js';
import { handleApiError } from '$lib/stores/toast.js';

const API_BASE = '/api/v1';

// Paths excluded from token refresh on 401 (self-managed auth endpoints)
const AUTH_PATHS = new Set(['/auth/refresh', '/auth/login', '/auth/register']);

let isRefreshing = false;
let refreshPromise: Promise<TokenResponse | null> | null = null;

function isAuthPath(path: string): boolean {
	return AUTH_PATHS.has(path);
}

async function refreshAccessToken(): Promise<TokenResponse | null> {
	const refreshToken = authStore.refreshToken;
	if (!refreshToken) return null;

	try {
		const response = await fetch(`${API_BASE}/auth/refresh`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ refresh_token: refreshToken })
		});

		if (!response.ok) {
			authStore.clearTokens();
			return null;
		}

		const tokens: TokenResponse = await response.json();
		authStore.setTokens(tokens);
		return tokens;
	} catch {
		return null;
	}
}

async function getRefreshedToken(): Promise<TokenResponse | null> {
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

/**
 * Build request headers with auth token from authStore (single source of truth).
 * authStore is updated synchronously by refreshAccessToken before this is called for retries.
 */
function buildHeaders(options: RequestInit): Record<string, string> {
	const headers: Record<string, string> = {
		'Content-Type': 'application/json',
		...((options.headers as Record<string, string>) || {})
	};

	const token = authStore.accessToken;
	if (token) {
		headers['Authorization'] = `Bearer ${token}`;
	}

	return headers;
}

/**
 * Make an authenticated API request with automatic token refresh on 401.
 *
 * Flow:
 * 1. Read token from authStore (single source of truth)
 * 2. If 401 on non-auth path, refresh token via API
 * 3. refreshAccessToken updates authStore.setTokens (synchronous state update)
 * 4. Retry with buildHeaders which reads the now-fresh token from authStore
 * 5. If refresh also fails, clear tokens and redirect to login
 */
async function request<T>(
	path: string,
	options: RequestInit = {}
): Promise<T> {
	const url = `${API_BASE}${path}`;
	const headers = buildHeaders(options);

	const response = await fetch(url, { ...options, headers });

	if (response.status === 401 && !isAuthPath(path)) {
		const newTokens = await getRefreshedToken();

		if (newTokens) {
			// At this point authStore.setTokens was already called by refreshAccessToken,
			// so buildHeaders() reads the fresh token from authStore.
			const retryHeaders = buildHeaders(options);

			const retryResponse = await fetch(url, { ...options, headers: retryHeaders });

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

/**
 * Make an authenticated blob request with automatic token refresh on 401.
 *
 * Returns the raw Response for non-auth cases (caller handles HTTP status).
 * On auth failure (401 + refresh failed), throws and redirects to login.
 *
 * Why not throw on all non-ok? The caller (export.ts) needs to check
 * response.ok and handle specific error statuses (e.g., 422 for invalid format).
 */
async function requestBlob(
	path: string,
	options: RequestInit = {}
): Promise<Response> {
	const url = `${API_BASE}${path}`;
	const headers = buildHeaders(options);

	const response = await fetch(url, { ...options, headers });

	if (response.status === 401 && !isAuthPath(path)) {
		const newTokens = await getRefreshedToken();

		if (newTokens) {
			const retryHeaders = buildHeaders(options);
			const retryResponse = await fetch(url, { ...options, headers: retryHeaders });
			if (retryResponse.ok) return retryResponse;
			// Retry failed — fall through to auth failure handling
		}

		// Refresh failed or retry still 401 — redirect to login
		authStore.clearTokens();
		if (typeof window !== 'undefined') {
			window.location.href = '/login';
		}
		throw new Error('Authentication failed — session expired');
	}

	// Return raw response — caller handles non-2xx (export.ts checks response.ok)
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

	/** Fetch a binary endpoint. Returns raw Response for non-auth cases. Throws on auth failure. */
	blob(path: string, options?: RequestInit): Promise<Response> {
		return requestBlob(path, { ...options, method: 'GET' });
	}
};

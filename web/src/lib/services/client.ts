import type { ApiError } from '$lib/types/index.js';

const API_BASE = '/api/v1';

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

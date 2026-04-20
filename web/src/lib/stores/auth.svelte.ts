import { browser } from '$app/environment';
import type { User, TokenResponse } from '$lib/types/index.js';
import { authService } from '$lib/services/auth.js';

function getStoredToken(key: string): string | null {
	if (!browser) return null;
	return localStorage.getItem(key);
}

function setStoredToken(key: string, value: string) {
	if (!browser) return;
	localStorage.setItem(key, value);
}

function removeStoredToken(key: string) {
	if (!browser) return;
	localStorage.removeItem(key);
}

function createAuthStore() {
	let user = $state<User | null>(null);
	let accessToken = $state<string | null>(getStoredToken('access_token'));
	let refreshToken = $state<string | null>(getStoredToken('refresh_token'));
	let isLoading = $state(false);

	const isAuthenticated = $derived(!!accessToken);

	function setTokens(tokens: TokenResponse) {
		accessToken = tokens.access_token;
		refreshToken = tokens.refresh_token;
		setStoredToken('access_token', tokens.access_token);
		setStoredToken('refresh_token', tokens.refresh_token);
	}

	function clearTokens() {
		user = null;
		accessToken = null;
		refreshToken = null;
		removeStoredToken('access_token');
		removeStoredToken('refresh_token');
	}

	async function login(email: string, password: string) {
		isLoading = true;
		try {
			const tokens = await authService.login({ email, password });
			setTokens(tokens);
			await fetchUser();
			return tokens;
		} finally {
			isLoading = false;
		}
	}

	async function register(email: string, password: string) {
		isLoading = true;
		try {
			const tokens = await authService.register({ email, password });
			setTokens(tokens);
			await fetchUser();
			return tokens;
		} finally {
			isLoading = false;
		}
	}

	async function logout() {
		try {
			await authService.logout();
		} catch {
			// ignore logout errors
		}
		clearTokens();
	}

	async function fetchUser(): Promise<User | null> {
		if (!accessToken) return null;
		try {
			user = await authService.getMe();
			return user;
		} catch {
			clearTokens();
			return null;
		}
	}

	async function restore() {
		if (!accessToken) return;
		isLoading = true;
		try {
			await fetchUser();
		} finally {
			isLoading = false;
		}
	}

	return {
		get user() { return user; },
		get accessToken() { return accessToken; },
		get refreshToken() { return refreshToken; },
		get isLoading() { return isLoading; },
		get isAuthenticated() { return isAuthenticated; },
		login,
		register,
		logout,
		fetchUser,
		restore
	};
}

export const authStore = createAuthStore();

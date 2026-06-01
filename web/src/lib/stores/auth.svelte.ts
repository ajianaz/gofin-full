import type { User, TokenResponse } from '$lib/types/index.js';
import { authService } from '$lib/services/auth.js';
import { api } from '$lib/services/client.js';

function createAuthStore() {
	let user = $state<User | null>(null);
	let accessToken = $state<string | null>(null);
	let refreshToken = $state<string | null>(null);
	let isLoading = $state(false);

	const isAuthenticated = $derived(!!accessToken);

	// Token persistence is now handled server-side via httpOnly cookies.
	// The reactive state serves two purposes:
	//   1. Authorization header fallback during migration (backward compat)
	//   2. Client-side isAuthenticated derived state for UI reactivity
	// After page reload, accessToken starts null; the server authenticates
	// via httpOnly cookie. The in-memory token is populated by setTokens()
	// when tokens come back from login/register/refresh responses.
	function setTokens(tokens: TokenResponse) {
		accessToken = tokens.access_token;
		refreshToken = tokens.refresh_token;
	}

	function clearTokens() {
		user = null;
		accessToken = null;
		refreshToken = null;
	}

	async function setupGroup(): Promise<void> {
		if (!refreshToken) return;
		try {
			const groupsRes = await api.get<{ data: { id: string; attributes: { title: string } }[] }>('/groups');
			let groupId: string | undefined;

			if (groupsRes.data.length === 0) {
				const created = await api.post<{ data: { id: string } }>('/groups', { title: 'Default' });
				groupId = created.data.id;
			} else {
				groupId = groupsRes.data[0].id;
			}

			if (groupId) {
				const switchRes = await api.post<{ tokens?: TokenResponse }>('/groups/switch', { user_group_id: groupId });
				if (switchRes?.tokens) {
					setTokens(switchRes.tokens);
				}
			}
		} catch (err) {
			console.error('Failed to setup group:', err);
		}
	}

	async function login(email: string, password: string) {
		isLoading = true;
		try {
			const tokens = await authService.login({ email, password });
			setTokens(tokens);
			await fetchUser();
			await setupGroup();
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
			await setupGroup();
			return tokens;
		} finally {
			isLoading = false;
		}
	}

	async function logout() {
		try {
			await authService.logout(refreshToken || undefined);
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
		// On page load, attempt to fetch user via httpOnly cookies.
		// We call authService.getMe() directly instead of fetchUser() because
		// fetchUser() guards on `if (!accessToken)` which is null after reload.
		// The server authenticates the request via the httpOnly cookie.
		isLoading = true;
		try {
			user = await authService.getMe();
			await setupGroup();
		} catch {
			// No valid cookie session — user is unauthenticated
			clearTokens();
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
		setTokens,
		clearTokens,
		login,
		register,
		logout,
		fetchUser,
		restore,
		setupGroup,
		/** Restore session from httpOnly cookies — bypasses accessToken guard.
		 *  Use in OAuth callback and similar flows where tokens are server-set. */
		async restoreSession() {
			isLoading = true;
			try {
				user = await authService.getMe();
				await setupGroup();
				return user;
			} catch {
				clearTokens();
				return null;
			} finally {
				isLoading = false;
			}
		}
	};
}

export const authStore = createAuthStore();

import { browser } from '$app/environment';
import { authStore } from '$lib/stores/auth.svelte.js';
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = async ({ data }) => {
	// On the server: use the cookie-based hint from +layout.server.ts.
	// On the client: prefer the live auth store (reactive to actual API state).
	// Once the client hydrates, authStore.isAuthenticated reflects real auth
	// state from API calls (httpOnly cookies). The SSR hint is only used for
	// the initial SSR render to prevent flash-of-login-page.
	const serverData = data as { session?: { isAuthenticated?: boolean } } | undefined;
	const isServerAuthenticated = serverData?.session?.isAuthenticated ?? false;

	return {
		session: {
			// Client: always trust the live auth store after hydration.
			// It tracks actual API responses (401 → false, valid session → true).
			// Server: use the cookie-presence hint for SSR routing only.
			isAuthenticated: browser ? authStore.isAuthenticated : isServerAuthenticated
		}
	};
};

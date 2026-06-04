import { browser } from '$app/environment';
import { authStore } from '$lib/stores/auth.svelte.js';
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = async ({ data }) => {
	// Prefer server-provided session state (from +layout.server.ts cookie check)
	// over the client-only default. On the client, once hydrated, the auth store
	// takes over reactivity via httpOnly cookie-based API calls.
	const isServerAuthenticated = data?.session?.isAuthenticated ?? false;

	return {
		session: {
			isAuthenticated: browser
				? authStore.isAuthenticated || isServerAuthenticated
				: isServerAuthenticated
		}
	};
};

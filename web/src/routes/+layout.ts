import { browser } from '$app/environment';
import { authStore } from '$lib/stores/auth.svelte.js';
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = () => {
	// On the client the auth store reads localStorage synchronously at
	// creation time, so `isAuthenticated` is already available for child
	// load functions that need to make access-control decisions before the
	// page renders (no flash of unauthorized content).
	return {
		session: {
			isAuthenticated: browser ? authStore.isAuthenticated : false
		}
	};
};

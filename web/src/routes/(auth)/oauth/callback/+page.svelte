<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { authStore } from '$lib/stores/auth.svelte.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import { Card, CardContent } from '$lib/components/ui/card/index.js';

	const t = localeStore.t;

	let error = $state<string | null>(null);

	onMount(async () => {
		const hash = window.location.hash;

		// Immediately strip tokens from the URL fragment to prevent them from
		// persisting in browser history.  Use replaceState so the token-bearing
		// URL is never recorded.
		if (hash) {
			window.history.replaceState(null, '', window.location.pathname + window.location.search);
		}

		if (!hash) {
			// No tokens in hash — try reading from URL params as fallback
			// (in case BE returns JSON and we need to fetch manually)
			const params = new URLSearchParams(window.location.search);
			const code = params.get('code');
			const state = params.get('state');

			if (code && state) {
				// Strip query params immediately too
				window.history.replaceState(null, '', window.location.pathname);

				// Fetch provider info to know which provider to call
				try {
					const providerRes = await fetch('/api/v1/auth/provider');
					if (providerRes.ok) {
						const providerData = await providerRes.json();
						const provider = providerData.provider;
						if (provider && provider.toLowerCase() !== 'local' && provider.toLowerCase() !== 'disabled') {
							const cbRes = await fetch(`/api/v1/auth/${provider.toLowerCase()}/callback?code=${encodeURIComponent(code)}&state=${encodeURIComponent(state)}`);
							if (cbRes.ok) {
								const tokens = await cbRes.json();
								authStore.setTokens(tokens);
								await authStore.fetchUser();
								await authStore.setupGroup();
								goto('/dashboard');
								return;
							}
						}
					}
				} catch (err) {
					console.error('OAuth callback failed:', err);
				}
				error = t('auth.oauth.callbackError');
				return;
			}

			// Nothing to process — shouldn't happen normally
			error = t('auth.oauth.callbackError');
			return;
		}

		// Parse tokens from URL fragment: #access_token=...&refresh_token=...
		const fragmentParams = new URLSearchParams(hash.slice(1));
		const accessToken = fragmentParams.get('access_token');
		const refreshToken = fragmentParams.get('refresh_token');

		if (!accessToken || !refreshToken) {
			error = t('auth.oauth.callbackError');
			return;
		}

		// Store tokens using the auth store
		authStore.setTokens({
			access_token: accessToken,
			refresh_token: refreshToken,
			expires_in: 0,
			token_type: 'Bearer'
		});

		// Fetch user info and setup group
		try {
			await authStore.fetchUser();
			await authStore.setupGroup();
		} catch (err) {
			console.error('Failed to fetch user after OAuth:', err);
		}

		// Redirect to dashboard (fragment already cleared above)
		goto('/dashboard');
	});
</script>

<div class="flex items-center justify-center min-h-[calc(100vh-4rem)]">
	<Card class="w-full max-w-md mx-4">
		<CardContent class="flex flex-col items-center justify-center py-12 px-8">
			{#if error}
				<div class="text-center space-y-4">
					<div class="text-destructive text-lg font-medium">
						{error}
					</div>
					<a href="/login" class="text-primary font-medium hover:underline">
						← {t('auth.login.submit')}
					</a>
				</div>
			{:else}
				<div class="animate-spin h-8 w-8 border-4 border-primary border-t-transparent rounded-full" role="status"></div>
				<p class="mt-4 text-sm text-muted-foreground">{t('auth.oauth.processing')}</p>
			{/if}
		</CardContent>
	</Card>
</div>

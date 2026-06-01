<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { authStore } from '$lib/stores/auth.svelte.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import { Card, CardContent } from '$lib/components/ui/card/index.js';

	const t = localeStore.t;

	let error = $state<string | null>(null);

	onMount(async () => {
		// Clear any residual URL fragment/params from old OAuth flow
		window.history.replaceState(null, '', window.location.pathname);

		// With httpOnly cookies, tokens are set server-side on OAuth callback.
		// restoreSession() calls getMe() directly, bypassing the accessToken guard.
		// It also handles group setup internally.
		try {
			const user = await authStore.restoreSession();
			if (user) {
				goto('/dashboard');
				return;
			}
		} catch (err) {
			console.error('Failed to verify OAuth session:', err);
		}

		error = t('auth.oauth.callbackError');
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

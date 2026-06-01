<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { authStore } from '$lib/stores/auth.svelte.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import { api } from '$lib/services/client.js';
	const t = localeStore.t;

	let isAllowed = $state(false);
	let isChecking = $state(true);

	onMount(async () => {
		if (!authStore.isAuthenticated) {
			goto('/dashboard');
			return;
		}
		// Defense-in-depth UX check — actual enforcement is via backend RBAC (auth.AdminMiddleware)
		try {
			await api.get('/admin/users');
			isAllowed = true;
		} catch {
			goto('/dashboard');
		} finally {
			isChecking = false;
		}
	});
</script>

{#if isChecking}
	<div class="flex items-center justify-center min-h-[50vh]">
		<p class="text-muted-foreground">{t('common.loading')}</p>
	</div>
{:else if isAllowed}
	<slot />
{:else}
	<div class="flex items-center justify-center min-h-[50vh]">
		<p class="text-muted-foreground">{t('common.forbidden')}</p>
	</div>
{/if}

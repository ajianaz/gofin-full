<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { authStore } from '$lib/stores/auth.svelte.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import { LanguageSwitcher } from '$lib/components/shared/index.js';

	let { children } = $props();
	const t = localeStore.t;

	onMount(async () => {
		if (authStore.isAuthenticated) {
			goto('/dashboard');
		}
	});
</script>

<div class="flex min-h-screen">
	<div class="hidden lg:flex lg:w-1/2 bg-primary items-center justify-center p-12">
		<div class="text-center space-y-6">
			<div class="flex items-center justify-center gap-3">
				<div class="flex size-12 items-center justify-center rounded-xl bg-primary-foreground">
					<span class="text-2xl font-bold text-primary">G</span>
				</div>
				<span class="text-3xl font-bold text-primary-foreground">{t('app.name')}</span>
			</div>
			<p class="text-primary-foreground text-lg">{t('auth.layout.tagline')}</p>
			<p class="text-primary-foreground/80 text-sm max-w-xs mx-auto">
				{t('auth.layout.description')}
			</p>
		</div>
	</div>

	<div class="flex flex-1 flex-col items-center justify-center p-8 gap-6">
		<div class="self-end">
			<LanguageSwitcher />
		</div>
		<div class="w-full max-w-md">{@render children()}</div>
	</div>
</div>

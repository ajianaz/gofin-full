<script lang="ts">
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import type { Locale } from '$lib/stores/i18n.svelte.js';
	import { Globe } from '@lucide/svelte';

	const locales: { value: Locale; label: string }[] = [
		{ value: 'id', label: 'ID' },
		{ value: 'en', label: 'EN' }
	];
</script>

<div class="flex items-center gap-1.5">
	<Globe class="size-3.5 text-muted-foreground" />
	<div class="flex rounded-md border border-border text-xs">
		{#each locales as loc}
			<button
				class="px-1.5 py-0.5 transition-colors {localeStore.locale === loc.value
					? 'bg-primary text-primary-foreground'
					: 'text-muted-foreground hover:bg-muted'}"
				onclick={() => localeStore.setLocale(loc.value)}
			>
				{loc.label}
			</button>
		{/each}
	</div>
</div>

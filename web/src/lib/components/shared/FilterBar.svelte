<script lang="ts">
	import type { Snippet } from 'svelte';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Search } from '@lucide/svelte';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	let {
		searchPlaceholder,
		actions,
		children
	}: {
		searchPlaceholder?: string;
		actions?: Snippet;
		children?: Snippet;
	} = $props();

	let searchValue = $state('');
</script>

<div class="flex flex-wrap items-center gap-3 mb-4">
	<div class="relative">
		<Search class="absolute left-2.5 top-2.5 size-4 text-muted-foreground" />
		<Input
			placeholder={searchPlaceholder ?? t('common.search')}
			class="pl-9 w-64"
			bind:value={searchValue}
		/>
	</div>
	{#if children}
		{@render children()}
	{/if}
	{#if actions}
		<div class="ml-auto flex items-center gap-2">
			{@render actions()}
		</div>
	{/if}
</div>

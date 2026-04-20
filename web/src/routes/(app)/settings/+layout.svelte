<script lang="ts">
	import { page } from '$app/stores';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	let { children } = $props();

	const items = $derived([
		{ label: t('settings.profile'), href: '/settings/profile' },
		{ label: t('settings.preferences'), href: '/settings/preferences' },
		{ label: t('settings.notifications'), href: '/settings/notifications' },
		{ label: t('settings.apiKeys'), href: '/settings/api-keys' }
	]);
</script>

<div class="flex flex-col gap-4">
	<nav class="flex items-center gap-1 border-b pb-0">
		{#each items as item}
			{@const active = $page.url.pathname === item.href}
			<a
				href={item.href}
				class="px-4 py-2.5 text-sm font-medium transition-colors {active
					? 'border-b-2 border-primary text-foreground'
					: 'text-muted-foreground hover:text-foreground'}"
			>
				{item.label}
			</a>
		{/each}
	</nav>

	{@render children()}
</div>

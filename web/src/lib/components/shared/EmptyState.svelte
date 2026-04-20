<script lang="ts">
	import { Inbox } from '@lucide/svelte';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	let {
		title,
		description,
		icon = Inbox
	}: {
		title?: string;
		description?: string;
		icon?: typeof Inbox;
	} = $props();

	const IconComponent = icon;
</script>

<div class="flex flex-col items-center justify-center py-16 text-center">
	<div class="flex size-16 items-center justify-center rounded-full bg-muted">
		<IconComponent class="size-8 text-muted-foreground" />
	</div>
	<h3 class="mt-4 text-lg font-semibold text-foreground">{title ?? t('common.noData')}</h3>
	<p class="mt-1 text-sm text-muted-foreground">{description ?? t('common.noDataDesc')}</p>
</div>

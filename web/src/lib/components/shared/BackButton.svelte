<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import { ArrowLeft } from '@lucide/svelte';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	let { href, label }: { href: string; label?: string } = $props();
</script>

<Button variant="ghost" size="sm" class="mb-4 gap-2" asChild>
	<a href={href}>
		<ArrowLeft class="size-4" />
		{label ?? t('common.back')}
	</a>
</Button>

<script lang="ts">
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	let {
		status,
		label
	}: {
		status: 'active' | 'inactive' | 'overdue' | 'completed' | 'cancelled';
		label?: string;
	} = $props();

	const variants: Record<string, 'default' | 'secondary' | 'destructive' | 'outline'> = {
		active: 'default',
		inactive: 'secondary',
		overdue: 'destructive',
		completed: 'secondary',
		cancelled: 'outline'
	};

	const defaultLabels: Record<string, string> = {
		active: t('common.active'),
		inactive: t('common.inactive'),
		overdue: 'Overdue',
		completed: 'Completed',
		cancelled: 'Cancelled'
	};
</script>

<Badge variant={variants[status]}>
	{label ?? defaultLabels[status]}
</Badge>

<script lang="ts">
	import { goto } from '$app/navigation';
	import { Card, CardContent } from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Plus } from '@lucide/svelte';
	import { mockRecurring } from '$lib/data/mock-recurring.js';
	import { formatCurrency, formatDate } from '$lib/utils/format.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	const freqLabels: Record<string, string> = {
		daily: t('recurring.list.freqDaily'),
		weekly: t('recurring.list.freqWeekly'),
		monthly: t('recurring.list.freqMonthly'),
		quarterly: t('recurring.list.freqQuarterly'),
		yearly: t('recurring.list.freqYearly')
	};
</script>

<div class="flex flex-col gap-4">
	<div class="flex flex-wrap items-center justify-between gap-3">
		<div class="flex items-center gap-3">
			<h2 class="text-lg font-semibold text-foreground">{t('recurring.list.title')}</h2>
			<span class="inline-flex items-center rounded-2xl bg-primary px-2.5 py-0.5 text-xs font-medium text-primary-foreground">
				{mockRecurring.length}
			</span>
		</div>
		<Button size="sm" onclick={() => goto('/recurring/create')}>
			<Plus class="size-4" />
			{t('recurring.list.add')}
		</Button>
	</div>

	<Card>
		<CardContent class="p-0">
			{#each mockRecurring as rec}
				<div class="flex items-center justify-between px-5 py-4 border-b last:border-b-0">
					<div class="flex flex-col gap-1 min-w-0">
						<p class="text-sm font-semibold text-foreground truncate">{rec.title}</p>
						<div class="flex items-center gap-3 text-[13px] text-muted-foreground">
							<span class="font-medium">{freqLabels[rec.repeat_freq] ?? rec.repeat_freq}</span>
							<span class="text-muted-foreground/50">|</span>
							<span>{t('recurring.list.next')} {formatDate(rec.first_date)}</span>
						</div>
					</div>
					<div class="flex shrink-0 items-center gap-3">
						<span class="text-sm font-semibold text-green-600">+ {formatCurrency(rec.amount)}</span>
						{#if rec.active}
							<Badge variant="secondary" class="text-xs">{t('recurring.list.active')}</Badge>
						{:else}
							<Badge variant="outline" class="text-xs">{t('recurring.list.inactive')}</Badge>
						{/if}
					</div>
				</div>
			{/each}
		</CardContent>
	</Card>
</div>

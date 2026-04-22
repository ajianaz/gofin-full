<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Card, CardContent } from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import { Plus } from '@lucide/svelte';
	import { budgetService } from '$lib/services/index.js';
	import type { Budget } from '$lib/types/domain.js';
	import { formatCurrency, formatPercentage } from '$lib/utils/format.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	let items = $state<Budget[]>([]);
	let isLoading = $state(true);

	onMount(async () => {
		try {
			items = await budgetService.list();
		} catch (e) {
			console.error(e);
		} finally {
			isLoading = false;
		}
	});
</script>

<div class="flex flex-col gap-4">
	<div class="flex flex-wrap items-center justify-between gap-3">
		<div class="flex items-center gap-3">
			<h2 class="text-base font-semibold text-foreground">{t('budgets.list.title')}</h2>
			<span class="inline-flex items-center rounded-2xl bg-primary px-2.5 py-0.5 text-xs font-medium text-primary-foreground">
				{items.length}
			</span>
		</div>
		<Button size="sm" onclick={() => goto('/budgets/create')}>
			<Plus class="size-4" />
			{t('budgets.list.add')}
		</Button>
	</div>

	<div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
		{#if isLoading}
				<p class="col-span-full py-8 text-center text-sm text-muted-foreground">Memuat...</p>
			{:else}
			{#each items as budget}
			{@const pct = parseFloat(budget.budget_amount) > 0
				? (parseFloat(budget.spend_amount) / parseFloat(budget.budget_amount)) * 100
				: 0}
			{@const remaining = parseFloat(budget.budget_amount) - parseFloat(budget.spend_amount)}
			<Card class="hover:shadow-md transition-shadow">
				<CardContent class="p-5">
					<div class="flex items-center justify-between mb-3">
						<p class="text-base font-semibold text-foreground">{budget.name}</p>
						{#if budget.auto_budget_type && budget.auto_budget_type !== 'none'}
							<Badge variant="outline" class="text-xs">{t('budgets.list.auto')}</Badge>
						{/if}
					</div>

					<div class="mb-3">
						<div class="flex justify-between text-sm mb-1">
							<span class="text-muted-foreground">{t('budgets.list.used')}</span>
							<span class="font-medium text-foreground">{formatCurrency(budget.spend_amount)}</span>
						</div>
						<div class="flex justify-between text-sm mb-1">
							<span class="text-muted-foreground">{t('budgets.list.budget')}</span>
							<span class="font-medium text-foreground">{formatCurrency(budget.budget_amount)}</span>
						</div>
						<div class="flex justify-between text-sm">
							<span class="text-muted-foreground">{t('budgets.list.remaining')}</span>
							<span class="font-medium {remaining >= 0 ? 'text-green-600' : 'text-red-600'}">
								{formatCurrency(remaining.toString())}
							</span>
						</div>
					</div>

					<Progress value={Math.min(pct, 100)} class="h-2" />
					<p class="text-xs text-muted-foreground mt-1">{t('budgets.list.usedPercent', { pct: Math.round(pct) })}</p>
				</CardContent>
			</Card>
		{/each}
			{/if}
	</div>
</div>

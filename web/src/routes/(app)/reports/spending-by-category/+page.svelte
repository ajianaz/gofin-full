<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { ArrowLeft } from '@lucide/svelte';
	import { reportService } from '$lib/services/index.js';
	import { formatCurrency } from '$lib/utils/format.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	let loading = $state(true);
	let error = $state('');

	let categoryData: { name: string; amount: number }[] = $state([]);
	let totalSpent = $state(0);

	onMount(async () => {
		try {
			const data = await reportService.spendingByCategory();
			categoryData = data
				.map((c) => ({ name: c.category_name, amount: c.total }))
				.filter((c) => c.amount > 0)
				.sort((a, b) => b.amount - a.amount);
			totalSpent = categoryData.reduce((s, c) => s + c.amount, 0);
		} catch (e: any) {
			error = e?.detail || e?.message || 'Failed to load category spending';
		} finally {
			loading = false;
		}
	});
</script>

<div class="flex flex-col gap-4">
	<div class="flex flex-wrap items-center justify-between gap-3">
		<div class="flex items-center gap-3">
			<Button variant="ghost" size="sm" onclick={() => goto('/reports')}>
				<ArrowLeft class="size-4" />
				{t('common.back')}
			</Button>
			<h2 class="text-base font-semibold text-foreground">{t('reports.spendingByCategory.title')}</h2>
		</div>
	</div>

	{#if loading}
		<div class="flex items-center justify-center py-12">
			<p class="text-sm text-muted-foreground">{t('common.loading')}</p>
		</div>
	{:else if error}
		<div class="flex items-center justify-center py-12">
			<p class="text-sm text-destructive">{error}</p>
		</div>
	{:else}
		<Card>
			<CardHeader class="pb-1"><CardTitle class="text-sm font-semibold">{t('reports.spendingByCategory.totalSpending')}</CardTitle></CardHeader>
			<CardContent>
				<p class="text-xl font-bold text-destructive">{formatCurrency(totalSpent.toString())}</p>
				<p class="text-xs text-muted-foreground">{t('reports.spendingByCategory.categoryCount', { count: categoryData.length })}</p>
			</CardContent>
		</Card>

		<Card>
			<CardHeader>
				<CardTitle class="text-base">{t('reports.spendingByCategory.details')}</CardTitle>
			</CardHeader>
			<CardContent>
				{#if categoryData.length === 0}
					<p class="text-sm text-muted-foreground">{t('common.noData')}</p>
				{:else}
					<div class="flex flex-col gap-4">
						{#each categoryData as cat}
							<div>
								<div class="flex items-center justify-between mb-1.5">
									<span class="text-sm font-medium text-foreground">{cat.name}</span>
									<span class="text-sm font-medium text-foreground">{formatCurrency(cat.amount.toString())}</span>
								</div>
								<div class="h-0.5 w-full rounded-full bg-muted"></div>
							</div>
						{/each}
					</div>
				{/if}
			</CardContent>
		</Card>
	{/if}
</div>

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

	interface PeriodRow {
		period: string;
		income: number;
		expense: number;
		diff: number;
	}

	let periodData: PeriodRow[] = $state([]);

	onMount(async () => {
		try {
			const data = await reportService.spendingByPeriod();
			periodData = data
				.map((p) => ({
					period: p.period,
					income: p.income,
					expense: p.expense,
					diff: p.income - p.expense
				}))
				.sort((a, b) => b.period.localeCompare(a.period));
		} catch (e: any) {
			error = e?.detail || e?.message || 'Failed to load period spending';
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
			<h2 class="text-base font-semibold text-foreground">{t('reports.spendingByPeriod.title')}</h2>
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
			<CardHeader>
				<CardTitle class="text-base">{t('reports.spendingByPeriod.history')}</CardTitle>
			</CardHeader>
			<CardContent class="p-0">
				{#if periodData.length === 0}
					<p class="px-4 py-6 text-sm text-muted-foreground text-center">{t('common.noData')}</p>
				{:else}
					<div class="overflow-x-auto">
						<table class="w-full text-sm">
							<thead>
								<tr class="border-b bg-muted/50">
									<th class="text-left px-4 py-3 font-medium text-muted-foreground">{t('reports.spendingByPeriod.period')}</th>
									<th class="text-right px-4 py-3 font-medium text-muted-foreground">{t('reports.spendingByPeriod.income')}</th>
									<th class="text-right px-4 py-3 font-medium text-muted-foreground">{t('reports.spendingByPeriod.expense')}</th>
									<th class="text-right px-4 py-3 font-medium text-muted-foreground">{t('reports.spendingByPeriod.diff')}</th>
								</tr>
							</thead>
							<tbody>
								{#each periodData as row}
									<tr class="border-b last:border-b-0 hover:bg-muted/30">
										<td class="px-4 py-3 font-medium text-foreground">{row.period}</td>
										<td class="px-4 py-3 text-right text-green-600">{formatCurrency(row.income.toString())}</td>
										<td class="px-4 py-3 text-right text-destructive">{formatCurrency(row.expense.toString())}</td>
										<td class="px-4 py-3 text-right font-medium {row.diff >= 0 ? 'text-green-600' : 'text-destructive'}">
											{row.diff >= 0 ? '+' : '-'}{formatCurrency(Math.abs(row.diff).toString())}
										</td>
									</tr>
								{/each}
							</tbody>
						</table>
					</div>
				{/if}
			</CardContent>
		</Card>
	{/if}
</div>

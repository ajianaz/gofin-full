<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '$lib/components/ui/table/index.js';
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
					<Table>
							<TableHeader>
								<TableRow>
									<TableHead>{t('reports.spendingByPeriod.period')}</TableHead>
									<TableHead class="text-right">{t('reports.spendingByPeriod.income')}</TableHead>
									<TableHead class="text-right">{t('reports.spendingByPeriod.expense')}</TableHead>
									<TableHead class="text-right">{t('reports.spendingByPeriod.diff')}</TableHead>
								</TableRow>
							</TableHeader>
							<TableBody>
								{#each periodData as row}
						<TableRow>
							<TableCell class="px-4 py-3 font-medium text-foreground">{row.period}</TableCell>
							<TableCell class="px-4 py-3 text-right text-green-600">{formatCurrency(row.income.toString())}</TableCell>
							<TableCell class="px-4 py-3 text-right text-destructive">{formatCurrency(row.expense.toString())}</TableCell>
							<TableCell class="px-4 py-3 text-right font-medium {row.diff >= 0 ? 'text-green-600' : 'text-destructive'}">
								{row.diff >= 0 ? '+' : '-'}{formatCurrency(Math.abs(row.diff).toString())}
							</TableCell>
						</TableRow>
								{/each}
							</TableBody>
						</Table>
				{/if}
			</CardContent>
		</Card>
	{/if}
</div>

<script lang="ts">
	import { onMount } from 'svelte';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { reportService } from '$lib/services/index.js';
	import { formatCurrency } from '$lib/utils/format.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	let loading = $state(true);
	let error = $state('');

	let totalIncome = $state(0);
	let totalExpense = $state(0);
	let diff = $state(0);
	let netWorthVal = $state(0);
	let transactionCount = $state(0);

	let categorySpending: { name: string; amount: number }[] = $state([]);
	let monthData: { label: string; income: number; expense: number }[] = $state([]);

	let maxCatAmount = $state(1);
	let maxMonthVal = $state(1);

	const barColors = ['#3b82f6', '#ef4444', '#f59e0b', '#10b981', '#8b5cf6', '#ec4899'];

	onMount(async () => {
		try {
			const [netWorthRes, categoryRes, periodRes] = await Promise.all([
				reportService.netWorth(),
				reportService.spendingByCategory(),
				reportService.spendingByPeriod()
			]);

			totalIncome = netWorthRes.total_income;
			totalExpense = netWorthRes.total_expense;
			diff = netWorthRes.net_income;
			netWorthVal = netWorthRes.net_income;
			transactionCount = netWorthRes.transaction_count;

			categorySpending = categoryRes
				.map((c) => ({ name: c.category_name, amount: c.total }))
				.filter((c) => c.amount > 0)
				.sort((a, b) => b.amount - a.amount);

			maxCatAmount = Math.max(...categorySpending.map((c) => c.amount), 1);

			// Take last 6 periods for the trend chart
			const last6 = periodRes.slice(-6);
			monthData = last6.map((p) => ({
				label: p.period,
				income: p.income,
				expense: p.expense
			}));
			maxMonthVal = Math.max(...monthData.flatMap((m) => [m.income, m.expense]), 1);
		} catch (e: any) {
			error = e?.detail || e?.message || 'Failed to load reports';
		} finally {
			loading = false;
		}
	});
</script>

<div class="flex flex-col gap-4">
	<div class="flex flex-wrap items-center justify-between gap-3">
		<h2 class="text-lg font-semibold text-foreground">{t('reports.title')}</h2>
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
		<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
			<Card>
				<CardHeader class="pb-1"><CardTitle class="text-sm font-semibold">{t('reports.income')}</CardTitle></CardHeader>
				<CardContent>
					<p class="text-xl font-bold text-green-600">{formatCurrency(totalIncome.toString())}</p>
					<p class="text-xs text-muted-foreground">{transactionCount} transactions</p>
				</CardContent>
			</Card>
			<Card>
				<CardHeader class="pb-1"><CardTitle class="text-sm font-semibold">{t('reports.expense')}</CardTitle></CardHeader>
				<CardContent>
					<p class="text-xl font-bold text-destructive">{formatCurrency(totalExpense.toString())}</p>
					<p class="text-xs text-muted-foreground">{t('reports.spendingByPeriod.title')}</p>
				</CardContent>
			</Card>
			<Card>
				<CardHeader class="pb-1"><CardTitle class="text-sm font-semibold">{t('reports.diff')}</CardTitle></CardHeader>
				<CardContent>
					<p class="text-xl font-bold {diff >= 0 ? 'text-green-600' : 'text-destructive'}">{formatCurrency(Math.abs(diff).toString())}</p>
					<p class="text-xs text-muted-foreground">{t('reports.monthlySavings')}</p>
				</CardContent>
			</Card>
			<Card>
				<CardHeader class="pb-1"><CardTitle class="text-sm font-semibold">{t('reports.netWorth')}</CardTitle></CardHeader>
				<CardContent>
					<p class="text-xl font-bold text-foreground">{formatCurrency(netWorthVal.toString())}</p>
					<p class="text-xs text-muted-foreground">{t('reports.netWorth.title')}</p>
				</CardContent>
			</Card>
		</div>

		<div class="grid gap-4 lg:grid-cols-2">
			<Card>
				<CardHeader class="pb-1"><CardTitle class="text-sm font-semibold">{t('reports.spendingByCategory')}</CardTitle></CardHeader>
				<CardContent>
					{#if categorySpending.length === 0}
						<p class="text-sm text-muted-foreground">{t('common.noData')}</p>
					{:else}
						<div class="flex flex-col gap-2">
							{#each categorySpending as cat, i}
								<div class="flex items-center gap-2">
									<span class="w-24 shrink-0 text-xs text-foreground truncate">{cat.name}</span>
									<div
										class="h-5 rounded"
										style="width: {(cat.amount / maxCatAmount) * 70}%; min-width: 8px; background-color: {barColors[i % barColors.length]}"
									></div>
									<span class="text-xs text-muted-foreground">{formatCurrency(cat.amount.toString())}</span>
								</div>
							{/each}
						</div>
					{/if}
				</CardContent>
			</Card>

			<Card>
				<CardHeader class="pb-1"><CardTitle class="text-sm font-semibold">{t('reports.sixMonthTrend')}</CardTitle></CardHeader>
				<CardContent>
					{#if monthData.length === 0}
						<p class="text-sm text-muted-foreground">{t('common.noData')}</p>
					{:else}
						<div class="flex items-end gap-2" style="height: 160px">
							{#each monthData as m}
								<div class="flex flex-1 flex-col items-center gap-1">
									<div class="flex w-full items-end justify-center gap-[2px]" style="height: 140px">
										<div
											class="w-full rounded-t-sm bg-green-600"
											style="height: {(m.income / maxMonthVal) * 100}%; min-height: 4px"
										></div>
										<div
											class="w-full rounded-t-sm bg-destructive"
											style="height: {(m.expense / maxMonthVal) * 100}%; min-height: 4px"
										></div>
									</div>
									<span class="text-[11px] text-muted-foreground">{m.label}</span>
								</div>
							{/each}
						</div>
						<div class="mt-3 flex items-center gap-4">
							<div class="flex items-center gap-2">
								<div class="size-3 rounded-sm bg-green-600"></div>
								<span class="text-xs text-muted-foreground">{t('reports.income')}</span>
							</div>
							<div class="flex items-center gap-2">
								<div class="size-3 rounded-sm bg-destructive"></div>
								<span class="text-xs text-muted-foreground">{t('reports.expense')}</span>
							</div>
						</div>
					{/if}
				</CardContent>
			</Card>
		</div>
	{/if}
</div>

<script lang="ts">
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { ChevronDown, Download } from '@lucide/svelte';
	import { mockTransactions } from '$lib/data/mock-transactions.js';
	import { mockWallets } from '$lib/data/mock-wallets.js';
	import { mockCategories } from '$lib/data/mock-categories.js';
	import { formatCurrency } from '$lib/utils/format.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	let period = $state('monthly');

	const totalIncome = mockTransactions
		.filter((t) => t.type === 'deposit')
		.reduce((s, t) => s + parseFloat(t.amount), 0);

	const totalExpense = mockTransactions
		.filter((t) => t.type === 'withdrawal')
		.reduce((s, t) => s + Math.abs(parseFloat(t.amount)), 0);

	const diff = totalIncome - totalExpense;
	const netWorth = mockWallets
		.filter((w) => w.type !== 'liability')
		.reduce((s, w) => s + parseFloat(w.balance), 0)
		- mockWallets
			.filter((w) => w.type === 'liability')
			.reduce((s, w) => s + Math.abs(parseFloat(w.balance)), 0);

	const categorySpending = mockCategories.map((cat) => {
		const amount = mockTransactions
			.filter((t) => t.category === cat.id && t.type === 'withdrawal')
			.reduce((s, t) => s + Math.abs(parseFloat(t.amount)), 0);
		return { name: cat.name, amount };
	}).filter((c) => c.amount > 0).sort((a, b) => b.amount - a.amount);

	const maxCatAmount = Math.max(...categorySpending.map((c) => c.amount), 1);

	const barColors = ['#3b82f6', '#ef4444', '#f59e0b', '#10b981', '#8b5cf6', '#ec4899'];

	const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun'];
	const monthData = months.map((m, i) => {
		const income = Math.round(totalIncome * (0.7 + Math.random() * 0.6));
		const expense = Math.round(totalExpense * (0.6 + Math.random() * 0.8));
		return { label: m, income, expense };
	});
	const maxMonthVal = Math.max(...monthData.flatMap((m) => [m.income, m.expense]), 1);
</script>

<div class="flex flex-col gap-4">
	<div class="flex flex-wrap items-center justify-between gap-3">
		<h2 class="text-lg font-semibold text-foreground">{t('reports.title')}</h2>
		<div class="flex items-center gap-3">
			<div class="relative">
				<select
					bind:value={period}
					class="cn-input h-9 w-32 appearance-none bg-background pr-8 text-sm"
				>
					<option value="monthly">{t('reports.period.thisMonth')}</option>
					<option value="quarterly">{t('reports.period.quarterly')}</option>
					<option value="yearly">{t('reports.period.thisYear')}</option>
				</select>
				<ChevronDown class="pointer-events-none absolute right-2.5 top-1/2 size-3.5 -translate-y-1/2 text-muted-foreground" />
			</div>
			<Button variant="outline" size="sm">
				<Download class="size-4" />
				{t('reports.download')}
			</Button>
		</div>
	</div>

	<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
		<Card>
			<CardHeader class="pb-1"><CardTitle class="text-sm font-semibold">{t('reports.income')}</CardTitle></CardHeader>
			<CardContent>
				<p class="text-xl font-bold text-green-600">{formatCurrency(totalIncome.toString())}</p>
				<p class="text-xs text-muted-foreground">+12% {t('reports.fromLastMonth')}</p>
			</CardContent>
		</Card>
		<Card>
			<CardHeader class="pb-1"><CardTitle class="text-sm font-semibold">{t('reports.expense')}</CardTitle></CardHeader>
			<CardContent>
				<p class="text-xl font-bold text-destructive">{formatCurrency(totalExpense.toString())}</p>
				<p class="text-xs text-muted-foreground">-5% {t('reports.fromLastMonth')}</p>
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
				<p class="text-xl font-bold text-foreground">{formatCurrency(netWorth.toString())}</p>
				<p class="text-xs text-muted-foreground">+2.3% {t('reports.fromLastMonth')}</p>
			</CardContent>
		</Card>
	</div>

	<div class="grid gap-4 lg:grid-cols-2">
		<Card>
			<CardHeader class="pb-1"><CardTitle class="text-sm font-semibold">{t('reports.spendingByCategory')}</CardTitle></CardHeader>
			<CardContent>
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
			</CardContent>
		</Card>

		<Card>
			<CardHeader class="pb-1"><CardTitle class="text-sm font-semibold">{t('reports.sixMonthTrend')}</CardTitle></CardHeader>
			<CardContent>
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
			</CardContent>
		</Card>
	</div>
</div>

<script lang="ts">
	import { goto } from '$app/navigation';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import { ArrowLeft, ChevronDown } from '@lucide/svelte';
	import { mockTransactions } from '$lib/data/mock-transactions.js';
	import { mockCategories } from '$lib/data/mock-categories.js';
	import { mockBudgets } from '$lib/data/mock-budgets.js';
	import { formatCurrency } from '$lib/utils/format.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	let period = $state('monthly');

	const budget = mockBudgets[0];

	const categoryData = mockCategories.map((cat) => {
		const amount = mockTransactions
			.filter((t) => t.category === cat.id && t.type === 'withdrawal')
			.reduce((s, t) => s + Math.abs(parseFloat(t.amount)), 0);
		const budgetLimit = budget?.limits.find((l) => l.category === cat.id);
		return {
			name: cat.name,
			amount,
			budget: budgetLimit ? parseFloat(budgetLimit.amount) : 0,
			pct: budgetLimit ? (amount / parseFloat(budgetLimit.amount)) * 100 : 0
		};
	}).filter((c) => c.amount > 0).sort((a, b) => b.amount - a.amount);

	const totalSpent = categoryData.reduce((s, c) => s + c.amount, 0);
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
		<div class="relative">
			<select bind:value={period} class="cn-input h-9 w-32 appearance-none bg-background pr-8 text-sm">
				<option value="monthly">{t('reports.period.thisMonth')}</option>
				<option value="quarterly">{t('reports.period.quarterly')}</option>
				<option value="yearly">{t('reports.period.thisYear')}</option>
			</select>
			<ChevronDown class="pointer-events-none absolute right-2.5 top-1/2 size-3.5 -translate-y-1/2 text-muted-foreground" />
		</div>
	</div>

	<div class="grid gap-4 sm:grid-cols-2">
		<Card>
			<CardHeader class="pb-1"><CardTitle class="text-sm font-semibold">{t('reports.spendingByCategory.totalSpending')}</CardTitle></CardHeader>
			<CardContent>
				<p class="text-xl font-bold text-destructive">{formatCurrency(totalSpent.toString())}</p>
				<p class="text-xs text-muted-foreground">{t('reports.spendingByCategory.categoryCount', { count: categoryData.length })}</p>
			</CardContent>
		</Card>
		<Card>
			<CardHeader class="pb-1"><CardTitle class="text-sm font-semibold">{t('reports.spendingByCategory.budget')}</CardTitle></CardHeader>
			<CardContent>
				<p class="text-xl font-bold text-foreground">{budget ? formatCurrency(budget.budget_amount) : 'Rp 0'}</p>
				<p class="text-xs text-muted-foreground">{budget?.name ?? '-'}</p>
			</CardContent>
		</Card>
	</div>

	<Card>
		<CardHeader>
			<CardTitle class="text-base">{t('reports.spendingByCategory.details')}</CardTitle>
		</CardHeader>
		<CardContent>
			<div class="flex flex-col gap-4">
				{#each categoryData as cat}
					<div>
						<div class="flex items-center justify-between mb-1.5">
							<span class="text-sm font-medium text-foreground">{cat.name}</span>
							<div class="text-right">
								<span class="text-sm font-medium text-foreground">{formatCurrency(cat.amount.toString())}</span>
								{#if cat.budget > 0}
									<span class="text-xs text-muted-foreground"> / {formatCurrency(cat.budget.toString())}</span>
								{/if}
							</div>
						</div>
						{#if cat.budget > 0}
							<Progress value={Math.min(cat.pct, 100)} class="h-2" />
							<p class="text-xs text-muted-foreground mt-0.5">{cat.pct.toFixed(1)}%</p>
						{:else}
							<div class="h-0.5 w-full rounded-full bg-muted"></div>
						{/if}
					</div>
				{/each}
			</div>
		</CardContent>
	</Card>
</div>

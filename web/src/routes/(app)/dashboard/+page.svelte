<script lang="ts">
	import { PageHeader, StatCard, AmountDisplay } from '$lib/components/shared/index.js';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import { Wallet, TrendingUp, TrendingDown, PiggyBank } from '@lucide/svelte';
	import { mockWallets } from '$lib/data/mock-wallets.js';
	import { mockTransactions } from '$lib/data/mock-transactions.js';
	import { mockBudgets } from '$lib/data/mock-budgets.js';
	import { formatCurrency, formatDate } from '$lib/utils/format.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	const totalBalance = mockWallets.reduce((sum, w) => sum + parseFloat(w.balance), 0);
	const totalIncome = mockTransactions.filter((t) => t.type === 'deposit').reduce((sum, t) => sum + parseFloat(t.amount), 0);
	const totalExpense = mockTransactions.filter((t) => t.type === 'withdrawal').reduce((sum, t) => sum + Math.abs(parseFloat(t.amount)), 0);
	const recentTransactions = mockTransactions.slice(0, 8);
	const budget = mockBudgets[0];
	const budgetPercent = budget ? (parseFloat(budget.spend_amount) / parseFloat(budget.budget_amount)) * 100 : 0;
</script>

<PageHeader title={t('dashboard.title')} />

<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4 mb-6">
	<StatCard title={t('dashboard.totalBalance')} value={formatCurrency(totalBalance.toString())} icon={Wallet} trend={{ value: '+12%', positive: true }} />
	<StatCard title={t('dashboard.income')} value={formatCurrency(totalIncome.toString())} icon={TrendingUp} trend={{ value: '+5%', positive: true }} />
	<StatCard title={t('dashboard.expense')} value={formatCurrency(totalExpense.toString())} icon={TrendingDown} trend={{ value: '-8%', positive: false }} />
	<StatCard title={t('dashboard.savings')} value={formatCurrency('12000000')} icon={PiggyBank} trend={{ value: '+3%', positive: true }} />
</div>

<div class="grid gap-4 lg:grid-cols-3">
	<Card class="lg:col-span-2">
		<CardHeader>
			<CardTitle class="text-base font-semibold">{t('dashboard.recentTransactions')}</CardTitle>
		</CardHeader>
		<CardContent>
			<div class="space-y-3">
				{#each recentTransactions as tx}
					<div class="flex items-center justify-between">
						<div class="min-w-0 flex-1">
							<p class="text-sm font-medium text-foreground truncate">{tx.description}</p>
							<p class="text-xs text-muted-foreground">{formatDate(tx.date)} &middot; {tx.category_name}</p>
						</div>
						<AmountDisplay amount={tx.amount} class="text-sm" />
					</div>
				{/each}
			</div>
			<a href="/transactions" class="mt-4 block text-sm text-primary font-medium hover:underline">{t('dashboard.viewAll')}</a>
		</CardContent>
	</Card>

	<Card>
		<CardHeader>
			<CardTitle class="text-base font-semibold">{t('dashboard.spendingByCategory')}</CardTitle>
		</CardHeader>
		<CardContent>
			{#if budget && budget.limits.length > 0}
				<div class="space-y-4">
					{#each budget.limits as limit}
						{@const pct = (parseFloat(limit.spend) / parseFloat(limit.amount)) * 100}
						<div>
							<div class="flex items-center justify-between mb-1">
								<span class="text-sm text-foreground">{limit.category_name}</span>
								<span class="text-xs text-muted-foreground">{formatCurrency(limit.spend)} / {formatCurrency(limit.amount)}</span>
							</div>
							<Progress value={pct} class="h-2" />
						</div>
					{/each}
				</div>
			{/if}
		</CardContent>
	</Card>
</div>

<script lang="ts">
	import { onMount } from 'svelte';
	import { PageHeader, StatCard, AmountDisplay } from '$lib/components/shared/index.js';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import { Wallet, TrendingUp, TrendingDown, PiggyBank, Plus, ArrowRight } from '@lucide/svelte';
	import { walletService, transactionService, budgetService } from '$lib/services/index.js';
	import { formatCurrency, formatDate } from '$lib/utils/format.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import type { Account, Transaction, Budget } from '$lib/types/domain.js';
import { Skeleton } from '$lib/components/ui/skeleton/index.js';
import { Button } from '$lib/components/ui/button/index.js';
import EmptyState from '$lib/components/shared/EmptyState.svelte';

const t = localeStore.t;

let wallets = $state<Account[]>([]);
let transactions = $state<Transaction[]>([]);
let budgets = $state<Budget[]>([]);
let isLoading = $state(true);

const isEmpty = $derived(wallets.length === 0);

	// Group wallets by currency for accurate totals
	const walletCurrencyGroups = $derived(() => {
		const groups = new Map<string, { wallets: Account[]; symbol: string; decimal: number }>();
		for (const w of wallets) {
			const code = w.currency_code || 'USD';
			if (!groups.has(code)) {
				groups.set(code, { wallets: [], symbol: w.currency_symbol || '$', decimal: w.currency_decimal_places ?? 2 });
			}
			groups.get(code)!.wallets.push(w);
		}
		return groups;
	});

	const totalBalance = $derived(wallets.reduce((sum, w) => sum + parseFloat(w.balance || '0'), 0));
	const totalBalanceSymbol = $derived(wallets.length > 0 ? (wallets[0].currency_symbol || '$') : '$');
	const totalIncome = $derived(transactions.filter((tx) => tx.type === 'deposit').reduce((sum, tx) => sum + parseFloat(tx.amount || '0'), 0));
	const totalExpense = $derived(transactions.filter((tx) => tx.type === 'withdrawal').reduce((sum, tx) => sum + Math.abs(parseFloat(tx.amount || '0')), 0));
	const recentTransactions = $derived(transactions.slice(0, 8));
	const budget = $derived(budgets[0] || null);
	const budgetPercent = $derived(budget ? (parseFloat(budget.spend_amount || '0') / parseFloat(budget.budget_amount || '1')) * 100 : 0);

	onMount(async () => {
		try {
			const [wRes, tRes, bRes] = await Promise.all([
				walletService.list(),
				transactionService.list({ per_page: 50 }),
				budgetService.list()
			]);
			wallets = wRes;
			transactions = tRes.data;
			budgets = bRes;
		} catch (e) {
			console.error('Failed to load dashboard data:', e);
		} finally {
			isLoading = false;
		}
	});
</script>

<PageHeader title={t('dashboard.title')} />

{#if isLoading}
<div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
		{#each Array(4) as _}
			<Card>
				<CardContent class="p-5">
					<Skeleton class="mb-2 h-4 w-24" />
					<Skeleton class="h-8 w-32" />
				</CardContent>
			</Card>
		{/each}
	</div>
	<Card class="mt-4">
		<CardContent class="p-5">
			<Skeleton class="mb-4 h-5 w-40" />
				{#each Array(5) as _}
					<div class="flex items-center justify-between py-2 border-b last:border-b-0">
						<Skeleton class="h-4 w-32" />
						<Skeleton class="h-4 w-20" />
					</div>
			{/each}
		</CardContent>
	</Card>
{:else if isEmpty}
	<div class="flex flex-col items-center justify-center py-16 text-center">
		<div class="flex size-20 items-center justify-center rounded-full bg-muted mb-4">
			<Wallet class="size-10 text-muted-foreground" />
		</div>
		<h3 class="text-xl font-semibold text-foreground mb-2">{t('dashboard.emptyState.title')}</h3>
		<p class="text-sm text-muted-foreground max-w-md mb-6">{t('dashboard.emptyState.description')}</p>
		<div class="flex gap-3">
			<a href="/wallets/create">
				<Button>
					<Plus class="size-4 mr-2" />
					{t('dashboard.emptyState.createWallet')}
				</Button>
			</a>
			<a href="/transactions/create">
				<Button variant="outline">
					{t('dashboard.emptyState.addTransaction')}
					<ArrowRight class="size-4 ml-2" />
				</Button>
			</a>
		</div>
	</div>
{:else}
<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4 mb-6">
		<StatCard title={t('dashboard.totalBalance')} value={formatCurrency(totalBalance.toString(), totalBalanceSymbol)} icon={Wallet} />
		<StatCard title={t('dashboard.income')} value={formatCurrency(totalIncome.toString(), totalBalanceSymbol)} icon={TrendingUp} />
		<StatCard title={t('dashboard.expense')} value={formatCurrency(totalExpense.toString(), totalBalanceSymbol)} icon={TrendingDown} />
		<StatCard title={t('dashboard.savings')} value={formatCurrency(Math.max(0, totalIncome - totalExpense).toString(), totalBalanceSymbol)} icon={PiggyBank} />
	</div>

	<div class="grid gap-4 lg:grid-cols-3 mb-6">
		{#each wallets as w}
			<Card>
				<CardContent class="p-5">
					<div class="flex items-center justify-between">
						<div class="min-w-0 flex-1">
							<p class="text-sm font-medium text-foreground truncate">{w.name}</p>
							<p class="text-xs text-muted-foreground">{w.currency_code || 'USD'}</p>
						</div>
						<AmountDisplay amount={w.balance || '0'} symbol={w.currency_symbol || '$'} class="text-lg font-semibold" />
					</div>
				</CardContent>
			</Card>
		{/each}
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
								<p class="text-xs text-muted-foreground">{formatDate(tx.date)} &middot; {tx.category_name || ''}</p>
							</div>
							<AmountDisplay amount={tx.amount} symbol={tx.currency_symbol} class="text-sm" />
						</div>
					{:else}
						<p class="text-sm text-muted-foreground">{t('common.noData')}</p>
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
							{@const pct = (parseFloat(limit.spend || '0') / parseFloat(limit.amount || '1')) * 100}
							<div>
								<div class="flex items-center justify-between mb-1">
									<span class="text-sm text-foreground">{limit.category_name}</span>
									<span class="text-xs text-muted-foreground">{formatCurrency(limit.spend || '0', totalBalanceSymbol)} / {formatCurrency(limit.amount || '0', totalBalanceSymbol)}</span>
								</div>
								<Progress value={pct} class="h-2" />
							</div>
						{/each}
					</div>
				{:else}
					<p class="text-sm text-muted-foreground">{t('common.noData')}</p>
				{/if}
			</CardContent>
		</Card>
	</div>
{/if}

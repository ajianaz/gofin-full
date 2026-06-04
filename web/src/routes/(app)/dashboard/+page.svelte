<script lang="ts">
	import { onMount } from 'svelte';
	import { PageHeader, StatCard, AmountDisplay } from '$lib/components/shared/index.js';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import { Wallet, TrendingUp, TrendingDown, PiggyBank, Plus, ArrowRight, AlertTriangle, RefreshCw } from '@lucide/svelte';
	import { walletService, transactionService, budgetService } from '$lib/services/index.js';
	import { formatCurrency, formatDate, getDefaultSymbol } from '$lib/utils/format.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import type { Account, Transaction, Budget } from '$lib/types/domain.js';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	import { handleApiError, showSuccessToast } from '$lib/stores/toast.js';
	const t = localeStore.t;

	let wallets = $state<Account[]>([]);
	let transactions = $state<Transaction[]>([]);
	let budgets = $state<Budget[]>([]);
	let isLoading = $state(true);
	let loadError = $state(false);

	const isEmpty = $derived(wallets.length === 0);

	// Detect unique currencies
	const uniqueCurrencies = $derived(() => {
		const set = new Set<string>();
		for (const w of wallets) {
			set.add(w.currency_code || 'USD');
		}
		return Array.from(set).sort();
	});

	const hasMultipleCurrencies = $derived(uniqueCurrencies().length > 1);

	// Group wallets by currency for accurate totals
	const walletCurrencyGroups = $derived(() => {
		const groups = new Map<string, { wallets: Account[]; symbol: string; decimal: number }>();
		for (const w of wallets) {
			const code = w.currency_code || 'USD';
			if (!groups.has(code)) {
				groups.set(code, { wallets: [], symbol: w.currency_symbol || getDefaultSymbol(), decimal: w.currency_decimal_places ?? 2 });
			}
			groups.get(code)!.wallets.push(w);
		}
		return groups;
	});

	// Per-currency totals
	const currencyBalances = $derived(() => {
		const result = new Map<string, { total: number; symbol: string; decimal: number; code: string }>();
		for (const w of wallets) {
			const code = w.currency_code || 'USD';
			if (!result.has(code)) {
				result.set(code, { total: 0, symbol: w.currency_symbol || getDefaultSymbol(), decimal: w.currency_decimal_places ?? 2, code });
			}
			result.get(code)!.total += parseFloat(w.balance || '0');
		}
		return result;
	});

	// Per-currency income/expense totals
	const currencyIncomeExpense = $derived(() => {
		const result = new Map<string, { income: number; expense: number; symbol: string; decimal: number; code: string }>();
		for (const tx of transactions) {
			const code = tx.currency_code || 'USD';
			if (!result.has(code)) {
				result.set(code, { income: 0, expense: 0, symbol: tx.currency_symbol || getDefaultSymbol(), decimal: tx.currency_decimal_places ?? 2, code });
			}
			if (tx.type === 'deposit') {
				result.get(code)!.income += parseFloat(tx.amount || '0');
			} else if (tx.type === 'withdrawal') {
				result.get(code)!.expense += Math.abs(parseFloat(tx.amount || '0'));
			}
		}
		return result;
	});

	// Legacy single-currency totals (for single currency display)
	const singleCurrencyCode = $derived(uniqueCurrencies()[0] || 'USD');
	const singleSymbol = $derived(wallets.length > 0 ? (wallets[0].currency_symbol || getDefaultSymbol()) : getDefaultSymbol());
	const singleDecimal = $derived(wallets.length > 0 ? (wallets[0].currency_decimal_places ?? 2) : 2);

	const totalBalance = $derived(currencyBalances().get(singleCurrencyCode)?.total ?? 0);
	const totalIncome = $derived(currencyIncomeExpense().get(singleCurrencyCode)?.income ?? 0);
	const totalExpense = $derived(currencyIncomeExpense().get(singleCurrencyCode)?.expense ?? 0);

	const recentTransactions = $derived(transactions.slice(0, 8));
	const budget = $derived(budgets[0] || null);
	const budgetPercent = $derived(budget ? (parseFloat(budget.spend_amount || '0') / parseFloat(budget.budget_amount || '1')) * 100 : 0);

	async function loadDashboard() {
		isLoading = true;
		loadError = false;
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
			loadError = true;
			handleApiError(e);
		} finally {
			isLoading = false;
		}
	}

	onMount(() => {
		loadDashboard();
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
{:else if loadError}
	<div class="flex flex-col items-center justify-center py-16 text-center">
		<div class="flex size-20 items-center justify-center rounded-full bg-destructive/10 mb-4">
			<AlertTriangle class="size-10 text-destructive" />
		</div>
		<h3 class="text-xl font-semibold text-foreground mb-2">{t('common.errorLoading')}</h3>
		<p class="text-sm text-muted-foreground max-w-md mb-6">{t('common.errorLoadingDescription')}</p>
		<Button onclick={loadDashboard}>
			<RefreshCw class="size-4 mr-2" />
			{t('common.retry')}
		</Button>
	</div>
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
	{#if hasMultipleCurrencies}
		<div class="mb-4 flex items-center gap-2 rounded-lg border border-amber-200 bg-amber-50 px-4 py-3 dark:border-amber-800 dark:bg-amber-950">
			<AlertTriangle class="size-4 shrink-0 text-amber-600 dark:text-amber-400" />
			<p class="text-sm text-amber-800 dark:text-amber-200">
				{t('dashboard.multipleCurrencyWarning')}
			</p>
		</div>
	{/if}

	{#if hasMultipleCurrencies}
		<!-- Per-currency stat cards -->
		{#each Array.from(currencyBalances().entries()) as [code, entry]}
			{@const ie = currencyIncomeExpense().get(code)}
			<div class="mb-6">
				<p class="mb-2 text-xs font-medium text-muted-foreground uppercase tracking-wide">{code}</p>
				<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
					<StatCard title={t('dashboard.totalBalance')} value={formatCurrency(entry.total.toString(), entry.symbol, entry.decimal)} icon={Wallet} />
					<StatCard title={t('dashboard.income')} value={formatCurrency((ie?.income ?? 0).toString(), entry.symbol, entry.decimal)} icon={TrendingUp} />
					<StatCard title={t('dashboard.expense')} value={formatCurrency((ie?.expense ?? 0).toString(), entry.symbol, entry.decimal)} icon={TrendingDown} />
					<StatCard title={t('dashboard.savings')} value={formatCurrency(Math.max(0, (ie?.income ?? 0) - (ie?.expense ?? 0)).toString(), entry.symbol, entry.decimal)} icon={PiggyBank} />
				</div>
			</div>
		{/each}
	{:else}
		<!-- Single currency stat cards -->
		<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4 mb-6">
			<StatCard title={t('dashboard.totalBalance')} value={formatCurrency(totalBalance.toString(), singleSymbol, singleDecimal)} icon={Wallet} />
			<StatCard title={t('dashboard.income')} value={formatCurrency(totalIncome.toString(), singleSymbol, singleDecimal)} icon={TrendingUp} />
			<StatCard title={t('dashboard.expense')} value={formatCurrency(totalExpense.toString(), singleSymbol, singleDecimal)} icon={TrendingDown} />
			<StatCard title={t('dashboard.savings')} value={formatCurrency(Math.max(0, totalIncome - totalExpense).toString(), singleSymbol, singleDecimal)} icon={PiggyBank} />
		</div>
	{/if}

	<div class="grid gap-4 lg:grid-cols-3 mb-6">
		{#each wallets as w}
			<Card>
				<CardContent class="p-5">
					<div class="flex items-center justify-between">
						<div class="min-w-0 flex-1">
							<p class="text-sm font-medium text-foreground truncate">{w.name}</p>
							<p class="text-xs text-muted-foreground">{w.currency_code || 'USD'}</p>
						</div>
						<AmountDisplay amount={w.balance || '0'} symbol={w.currency_symbol || getDefaultSymbol()} class="text-lg font-semibold" />
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
									<span class="text-xs text-muted-foreground">{formatCurrency(limit.spend || '0', singleSymbol, singleDecimal)} / {formatCurrency(limit.amount || '0', singleSymbol, singleDecimal)}</span>
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

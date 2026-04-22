<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Search, Plus, ChevronLeft, ChevronRight, ChevronDown } from '@lucide/svelte';
	import { transactionService, walletService, categoryService } from '$lib/services/index.js';
	import { formatAmount, formatDate } from '$lib/utils/format.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import type { Transaction, Account, Category } from '$lib/types/domain.js';
	const t = localeStore.t;

	let isLoading = $state(true);
	let items = $state<Transaction[]>([]);
	let wallets = $state<Account[]>([]);
	let categories = $state<Category[]>([]);

	let searchQuery = $state('');
	let typeFilter = $state('all');
	let accountFilter = $state('all');
	let categoryFilter = $state('all');
	let periodFilter = $state('this_month');

	const PER_PAGE = 5;
	let currentPage = $state(1);

	onMount(async () => {
		try {
			const [txRes, walletList, catList] = await Promise.all([
				transactionService.list(),
				walletService.list(),
				categoryService.list()
			]);
			items = txRes.data;
			wallets = walletList;
			categories = catList;
		} catch (e) {
			console.error('Failed to load transactions:', e);
		} finally {
			isLoading = false;
		}
	});

	let filtered = $derived(() => {
		let result = [...items];
		if (searchQuery) {
			const q = searchQuery.toLowerCase();
			result = result.filter((t) => t.description.toLowerCase().includes(q));
		}
		if (typeFilter !== 'all') result = result.filter((t) => t.type === typeFilter);
		if (accountFilter !== 'all') {
			result = result.filter(
				(t) => t.source_account === accountFilter || t.destination_account === accountFilter
			);
		}
		if (categoryFilter !== 'all') result = result.filter((t) => t.category === categoryFilter);
		return result;
	});

	let totalPages = $derived(Math.max(1, Math.ceil(filtered().length / PER_PAGE)));
	let paginated = $derived(() => {
		const start = (currentPage - 1) * PER_PAGE;
		return filtered().slice(start, start + PER_PAGE);
	});
	let showingText = $derived(() => {
		const total = filtered().length;
		if (total === 0) return t('transactions.list.noTransactions');
		const start = (currentPage - 1) * PER_PAGE + 1;
		const end = Math.min(currentPage * PER_PAGE, total);
		return t('transactions.list.showing', { start, end, total });
	});

	$effect(() => {
		searchQuery;
		typeFilter;
		accountFilter;
		categoryFilter;
		periodFilter;
		currentPage = 1;
	});

	function acctName(tx: Transaction): string {
		return tx.source_account_name || tx.destination_account_name || '-';
	}
</script>

<div class="flex flex-col gap-4">
	{#if isLoading}
		<p class="text-sm text-muted-foreground py-8 text-center">Memuat...</p>
	{:else}
		<div class="flex items-center gap-3">
			<h2 class="text-base font-semibold text-foreground">{t('transactions.list.title')}</h2>
			<span class="inline-flex items-center rounded-2xl bg-primary px-2.5 py-0.5 text-xs font-medium text-primary-foreground">
				{filtered().length}
			</span>
		</div>
		<div class="flex items-center gap-3">
			<div class="relative">
				<Search class="absolute left-2.5 top-2.5 size-4 text-muted-foreground" />
				<Input placeholder={t('transactions.list.searchPlaceholder')} class="w-60 pl-9" bind:value={searchQuery} />
			</div>
			<Button size="sm" onclick={() => goto('/transactions/create')}>
				<Plus class="size-4" />
				{t('transactions.list.add')}
			</Button>
		</div>

	<div class="flex flex-wrap items-center gap-3">
		<div class="relative">
			<select
				bind:value={typeFilter}
				class="cn-input h-9 w-40 appearance-none bg-background pr-8 text-sm"
			>
				<option value="all">{t('transactions.list.allTypes')}</option>
				<option value="withdrawal">{t('transactions.list.expense')}</option>
				<option value="deposit">{t('transactions.list.income')}</option>
				<option value="transfer">{t('transactions.list.transfer')}</option>
			</select>
			<ChevronDown class="pointer-events-none absolute right-2.5 top-1/2 size-3.5 -translate-y-1/2 text-muted-foreground" />
		</div>
		<div class="relative">
			<select
				bind:value={accountFilter}
				class="cn-input h-9 w-44 appearance-none bg-background pr-8 text-sm"
			>
				<option value="all">{t('transactions.list.allWallets')}</option>
				{#each wallets as w}
					<option value={w.id}>{w.name}</option>
				{/each}
			</select>
			<ChevronDown class="pointer-events-none absolute right-2.5 top-1/2 size-3.5 -translate-y-1/2 text-muted-foreground" />
		</div>
		<div class="relative">
			<select
				bind:value={categoryFilter}
				class="cn-input h-9 w-44 appearance-none bg-background pr-8 text-sm"
			>
				<option value="all">{t('transactions.list.allCategories')}</option>
				{#each categories as cat}
					<option value={cat.id}>{cat.name}</option>
				{/each}
			</select>
			<ChevronDown class="pointer-events-none absolute right-2.5 top-1/2 size-3.5 -translate-y-1/2 text-muted-foreground" />
		</div>
		<div class="relative">
			<select
				bind:value={periodFilter}
				class="cn-input h-9 w-44 appearance-none bg-background pr-8 text-sm"
			>
				<option value="this_month">{t('transactions.list.thisMonth')}</option>
				<option value="last_month">{t('transactions.list.lastMonth')}</option>
				<option value="this_week">{t('transactions.list.thisWeek')}</option>
				<option value="this_year">{t('transactions.list.thisYear')}</option>
				<option value="all">{t('transactions.list.allPeriods')}</option>
			</select>
			<ChevronDown class="pointer-events-none absolute right-2.5 top-1/2 size-3.5 -translate-y-1/2 text-muted-foreground" />
		</div>
	</div>

	<Card>
		<CardHeader class="pb-0">
			<CardTitle class="text-base">{t('transactions.list.history')}</CardTitle>
		</CardHeader>
		<CardContent class="p-0">
			<div class="overflow-x-auto">
				<table class="w-full text-sm">
					<thead>
						<tr class="border-b bg-muted/50">
							<th class="w-[120px] p-3 text-left font-semibold text-muted-foreground">{t('transactions.list.colDate')}</th>
							<th class="p-3 text-left font-semibold text-muted-foreground">{t('transactions.list.colDescription')}</th>
							<th class="w-[140px] p-3 text-left font-semibold text-muted-foreground">{t('transactions.list.colAmount')}</th>
							<th class="w-[140px] p-3 text-left font-semibold text-muted-foreground">{t('transactions.list.colCategory')}</th>
							<th class="w-[140px] p-3 text-left font-semibold text-muted-foreground">{t('transactions.list.colWallet')}</th>
						</tr>
					</thead>
					<tbody>
						{#each paginated() as tx}
							<tr class="border-b hover:bg-muted/30">
								<td class="p-3 whitespace-nowrap text-foreground">{formatDate(tx.date)}</td>
								<td class="p-3 text-foreground">{tx.description}</td>
								<td class="p-3 whitespace-nowrap">
									<span class="font-semibold {formatAmount(tx.amount).color}">{formatAmount(tx.amount).text}</span>
								</td>
								<td class="p-3 text-muted-foreground">{tx.category_name || '-'}</td>
								<td class="p-3 text-muted-foreground">{acctName(tx)}</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
			<div class="flex items-center justify-between border-t px-4 py-3">
				<span class="text-sm text-muted-foreground">{showingText()}</span>
				<div class="flex items-center gap-1">
					<Button
						variant="outline"
						size="sm"
						disabled={currentPage <= 1}
						onclick={() => currentPage--}
					>
						<ChevronLeft class="size-4" />
					</Button>
					<Button
						variant="outline"
						size="sm"
						disabled={currentPage >= totalPages}
						onclick={() => currentPage++}
					>
						<ChevronRight class="size-4" />
					</Button>
				</div>
			</div>
		</CardContent>
	</Card>
	{/if}
</div>

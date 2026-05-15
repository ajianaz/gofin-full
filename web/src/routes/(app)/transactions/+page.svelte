<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '$lib/components/ui/table/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Search, Plus, ChevronLeft, ChevronRight, ChevronDown, Trash2 } from '@lucide/svelte';
	import { transactionService, walletService, categoryService } from '$lib/services/index.js';
	import { formatAmount, formatDate } from '$lib/utils/format.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import type { Transaction, Account, Category } from '$lib/types/domain.js';
	import { ConfirmDialog } from '$lib/components/shared/index.js';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import { Select, SelectTrigger, SelectContent, SelectItem } from '$lib/components/ui/select/index.js';
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	const t = localeStore.t;

	let isLoading = $state(true);
	let errorMsg = $state('');
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
	let deleteTarget = $state<string | null>(null);
	let deleteOpen = $derived(deleteTarget !== null);

	async function handleDelete(id: string) {
		await transactionService.delete(id);
		items = items.filter((t) => t.id !== id);
	}

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
			errorMsg = t('common.error');
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
{#each Array(8) as _}
				<TableRow>
					<TableCell class="whitespace-nowrap"><Skeleton class="h-4 w-20" /></TableCell>
					<TableCell><Skeleton class="h-4 w-40" /></TableCell>
					<TableCell class="whitespace-nowrap"><Skeleton class="h-4 w-16" /></TableCell>
					<TableCell class="hidden md:table-cell"><Skeleton class="h-4 w-24" /></TableCell>
					<TableCell><Skeleton class="h-4 w-24" /></TableCell>
					<TableCell><Skeleton class="size-4" /></TableCell>
				</TableRow>
			{/each}
	{:else if errorMsg}
		<p class="text-sm text-destructive py-8 text-center">{errorMsg}</p>
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
			<Select bind:value={typeFilter}>
		<SelectTrigger class="w-40">
		</SelectTrigger>
		<SelectContent>
		<SelectItem value="all">{t('transactions.list.allTypes')}</SelectItem>
		<SelectItem value="withdrawal">{t('transactions.list.expense')}</SelectItem>
		<SelectItem value="deposit">{t('transactions.list.income')}</SelectItem>
		<SelectItem value="transfer">{t('transactions.list.transfer')}</SelectItem>
		</SelectContent>
</Select>
			<ChevronDown class="pointer-events-none absolute right-2.5 top-1/2 size-3.5 -translate-y-1/2 text-muted-foreground" />
		</div>
		<div class="relative">
			<Select bind:value={accountFilter}>
		<SelectTrigger class="w-44">
		</SelectTrigger>
		<SelectContent>
		<SelectItem value="all">{t('transactions.list.allWallets')}</SelectItem>
		{#each wallets as w}
<SelectItem value={w.id}>{w.name}</SelectItem>
{/each}
		</SelectContent>
</Select>
			<ChevronDown class="pointer-events-none absolute right-2.5 top-1/2 size-3.5 -translate-y-1/2 text-muted-foreground" />
		</div>
		<div class="relative">
			<Select bind:value={categoryFilter}>
		<SelectTrigger class="w-44">
		</SelectTrigger>
		<SelectContent>
		<SelectItem value="all">{t('transactions.list.allCategories')}</SelectItem>
		{#each categories as cat}
<SelectItem value={cat.id}>{cat.name}</SelectItem>
{/each}
		</SelectContent>
</Select>
			<ChevronDown class="pointer-events-none absolute right-2.5 top-1/2 size-3.5 -translate-y-1/2 text-muted-foreground" />
		</div>
		<div class="relative">
			<Select bind:value={periodFilter}>
		<SelectTrigger class="w-44">
		</SelectTrigger>
		<SelectContent>
		<SelectItem value="this_month">{t('transactions.list.thisMonth')}</SelectItem>
		<SelectItem value="last_month">{t('transactions.list.lastMonth')}</SelectItem>
		<SelectItem value="this_week">{t('transactions.list.thisWeek')}</SelectItem>
		<SelectItem value="this_year">{t('transactions.list.thisYear')}</SelectItem>
		<SelectItem value="all">{t('transactions.list.allPeriods')}</SelectItem>
		</SelectContent>
</Select>
			<ChevronDown class="pointer-events-none absolute right-2.5 top-1/2 size-3.5 -translate-y-1/2 text-muted-foreground" />
		</div>
	</div>

	<Card>
		<CardHeader class="pb-0">
			<CardTitle class="text-base">{t('transactions.list.history')}</CardTitle>
		</CardHeader>
		<CardContent class="p-0">
			<Table>
					<TableHeader>
						<TableRow>
							<TableHead class="w-[120px]">{t('transactions.list.colDate')}</TableHead>
							<TableHead>{t('transactions.list.colDescription')}</TableHead>
							<TableHead class="w-[140px]">{t('transactions.list.colAmount')}</TableHead>
							<TableHead class="hidden md:table-cell w-[140px]">{t('transactions.list.colCategory')}</TableHead>
							<TableHead class="hidden md:table-cell w-[140px]">{t('transactions.list.colWallet')}</TableHead>
							<TableHead class="w-[50px]"></TableHead>
						</TableRow>
					</TableHeader>
					<TableBody>
						{#each paginated() as tx}
							<TableRow>
								<TableCell class="whitespace-nowrap text-foreground">{formatDate(tx.date)}</TableCell>
								<TableCell class="text-foreground">{tx.description}</TableCell>
								<TableCell class="whitespace-nowrap">
									<span class="font-semibold {formatAmount(tx.amount).color}">{formatAmount(tx.amount).text}</span>
								</TableCell>
								<TableCell class="hidden md:table-cell text-muted-foreground">{tx.category_name || '-'}</TableCell>
								<TableCell class="hidden md:table-cell text-muted-foreground">{acctName(tx)}</TableCell>
								<TableCell>
									<button type="button" aria-label="{t('common.delete')}" class="text-muted-foreground hover:text-destructive transition-colors" onclick={() => (deleteTarget = tx.id)}>
										<Trash2 class="size-4" />
									</button>
								</TableCell>
							</TableRow>
						{:else}
							<TableRow><TableCell colspan={6}><EmptyState /></TableCell></TableRow>
						{/each}
					</TableBody>
				</Table>
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
	<ConfirmDialog
		bind:open={deleteOpen}
		title={t('common.delete')}
		description={t('common.deleteConfirm')}
		onConfirm={async () => {
			if (deleteTarget) {
				await handleDelete(deleteTarget);
				deleteTarget = null;
			}
		}}
	/>
</div>

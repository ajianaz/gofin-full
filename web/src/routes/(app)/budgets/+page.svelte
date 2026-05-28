<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { Card, CardContent } from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import { Plus, Trash2 } from '@lucide/svelte';
	import { budgetService, walletService } from '$lib/services/index.js';
	import type { Budget, Account } from '$lib/types/domain.js';
	import { formatCurrency, getDefaultSymbol } from '$lib/utils/format.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import { ConfirmDialog } from '$lib/components/shared/index.js';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	const t = localeStore.t;

	let items = $state<Budget[]>([]);
	let wallets = $state<Account[]>([]);
	let isLoading = $state(true);
	let errorMsg = $state('');
	let deleteTarget = $state<string | null>(null);
	let deleteOpen = $derived(deleteTarget !== null);

	async function handleDelete(id: string) {
		await budgetService.delete(id);
		items = items.filter((b) => b.id !== id);
	}

	onMount(async () => {
		try {
			const [budgetList, walletList] = await Promise.all([
				budgetService.list(),
				walletService.list()
			]);
			items = budgetList;
			wallets = walletList;
		} catch (e) {
			errorMsg = t('common.error');
			console.error(e);
		} finally {
			isLoading = false;
		}
	});

	let currencySymbol = $derived(wallets.length > 0 ? (wallets[0].currency_symbol || getDefaultSymbol()) : getDefaultSymbol());
	let currencyDecimal = $derived(wallets.length > 0 ? (wallets[0].currency_decimal_places ?? 2) : 2);
</script>

<div class="flex flex-col gap-4">
	<div class="flex flex-wrap items-center justify-between gap-3">
		<div class="flex items-center gap-3">
			<h2 class="text-base font-semibold text-foreground">{t('budgets.list.title')}</h2>
			<span class="inline-flex items-center rounded-2xl bg-primary px-2.5 py-0.5 text-xs font-medium text-primary-foreground">
				{items.length}
			</span>
		</div>
		<Button size="sm" onclick={() => goto('/budgets/create')}>
			<Plus class="size-4" />
			{t('budgets.list.add')}
		</Button>
	</div>

	<Card>
		<CardContent class="p-0">
			{#if isLoading}
				{#each Array(6) as _}
					<div class="flex items-center gap-4 px-5 py-4 border-b">
						<div class="flex size-10 shrink-0 items-center justify-center rounded-lg"><Skeleton class="size-5 rounded" /></div>
						<div class="flex flex-col gap-2 min-w-0 flex-1">
							<Skeleton class="h-4 w-32" />
							<Skeleton class="h-3 w-48" />
						</div>
						<div class="ml-auto text-right">
							<Skeleton class="h-4 w-16" />
						</div>
					</div>
				{/each}
			{:else if errorMsg}
				<p class="px-5 py-8 text-center text-sm text-destructive">{errorMsg}</p>
			{:else}
				{#each items as budget}
					{@const pct = parseFloat(budget.budget_amount) > 0
						? (parseFloat(budget.spend_amount) / parseFloat(budget.budget_amount)) * 100
						: 0}
					{@const remaining = parseFloat(budget.budget_amount) - parseFloat(budget.spend_amount)}
					<div class="flex items-center gap-4 px-5 py-4 border-b last:border-b-0">
						<div class="flex size-10 shrink-0 items-center justify-center rounded-lg bg-muted">
							<span class="text-lg font-bold text-foreground">{budget.name.charAt(0).toUpperCase()}</span>
						</div>
						<div class="flex flex-col gap-1.5 min-w-0 flex-1">
							<p class="text-sm font-semibold text-foreground truncate">{budget.name}</p>
							<div class="flex items-center gap-3 text-xs text-muted-foreground">
								<span>{t('budgets.list.used')}: {formatCurrency(budget.spend_amount, currencySymbol, currencyDecimal)}</span>
								<span>/</span>
								<span>{formatCurrency(budget.budget_amount, currencySymbol, currencyDecimal)}</span>
							</div>
							<Progress value={Math.min(pct, 100)} class="h-1.5 w-[120px]" />
						</div>
						<div class="ml-auto shrink-0 flex items-center gap-3">
							<div class="text-right">
								<p class="text-sm font-semibold {remaining >= 0 ? 'text-green-600' : 'text-red-600'}">
									{formatCurrency(remaining.toString(), currencySymbol, currencyDecimal)}
								</p>
								<p class="text-xs text-muted-foreground">{t('budgets.list.usedPercent', { pct: Math.round(pct) })}</p>
							</div>
							<Button variant="ghost" size="icon-sm" aria-label={t('common.delete')} class="text-muted-foreground hover:text-destructive" onclick={() => (deleteTarget = budget.id)}>
								<Trash2 class="size-4" />
							</Button>
						</div>
					</div>
				{:else}
					<EmptyState />
				{/each}
			{/if}
		</CardContent>
	</Card>
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

<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Card, CardContent } from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Plus, Trash2, ChevronDown } from '@lucide/svelte';
	import { billService } from '$lib/services/index.js';
	import type { Bill } from '$lib/types/domain.js';
	import { formatCurrency, formatDate } from '$lib/utils/format.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import { ConfirmDialog } from '$lib/components/shared/index.js';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import { Select, SelectTrigger, SelectContent, SelectItem } from '$lib/components/ui/select/index.js';
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	const t = localeStore.t;

	let items = $state<Bill[]>([]);
	let isLoading = $state(true);
	let errorMsg = $state('');
	let accountFilter = $state('all');
	let deleteTarget = $state<string | null>(null);
	let deleteOpen = $derived(deleteTarget !== null);

	async function handleDelete(id: string) {
		await billService.delete(id);
		items = items.filter((b) => b.id !== id);
	}

	onMount(async () => {
		try {
			items = await billService.list();
		} catch (e) {
			errorMsg = t('common.error');
			console.error(e);
		} finally {
			isLoading = false;
		}
	});

	let filtered = $derived(
		accountFilter === 'all' ? items : items.filter((b) => b.active)
	);
</script>

<div class="flex flex-col gap-4">
	<div class="flex flex-wrap items-center justify-between gap-3">
		<div class="flex items-center gap-2">
			<h2 class="text-base font-semibold text-foreground">{t('bills.list.title')}</h2>
			<span class="inline-flex items-center rounded-2xl bg-primary px-2.5 py-0.5 text-xs font-medium text-primary-foreground">
				{filtered.length}
			</span>
		</div>
		<div class="flex items-center gap-3">
			<div class="relative">
				<Select bind:value={accountFilter}>
		<SelectTrigger class="w-44">
		</SelectTrigger>
		<SelectContent>
		<SelectItem value="all">{t('bills.list.allWallets')}</SelectItem>
		<SelectItem value="active">{t('bills.list.activeOnly')}</SelectItem>
		</SelectContent>
</Select>
				<ChevronDown class="pointer-events-none absolute right-2.5 top-1/2 size-3.5 -translate-y-1/2 text-muted-foreground" />
			</div>
			<Button size="sm" onclick={() => goto('/bills/create')}>
				<Plus class="size-4" />
				{t('bills.list.add')}
			</Button>
		</div>
	</div>

	<Card>
		<CardContent class="p-0">
			<p class="px-5 pt-4 text-[13px] font-normal text-muted-foreground">{t('bills.list.activeList')}</p>
			{#if isLoading}
<Card>
		<CardContent class="p-0">
			{#each Array(5) as _}
				<div class="flex items-center gap-4 px-5 py-4 border-b">
					<div class="flex size-10 shrink-0 items-center justify-center rounded-lg"><Skeleton class="size-5 rounded" /></div>
					<div class="flex flex-col gap-2 min-w-0 flex-1">
						<Skeleton class="h-4 w-40" />
						<Skeleton class="h-3 w-24" />
					</div>
					<div class="ml-auto text-right">
						<Skeleton class="h-4 w-20" />
					</div>
				</div>
			{/each}
		</CardContent>
	</Card>
	{:else if errorMsg}
				<p class="px-5 py-8 text-center text-sm text-destructive">{errorMsg}</p>
			{:else}
				{#each filtered as bill}
					<div class="flex items-center justify-between px-5 py-3 border-b last:border-b-0">
						<div class="flex flex-col gap-1">
							<p class="text-sm font-semibold text-foreground">{bill.name}</p>
							<p class="text-[13px] text-muted-foreground">
								{#if bill.amount_min === bill.amount_max}
									{formatCurrency(bill.amount_min)}{t('bills.list.perMonth')}
								{:else}
									{formatCurrency(bill.amount_min)} — {formatCurrency(bill.amount_max)}
								{/if}
							</p>
						</div>
						<div class="flex shrink-0 items-center gap-3">
							<span class="text-[13px] text-muted-foreground">{formatDate(bill.next_date)}</span>
							{#if bill.active}
								<Badge variant="secondary" class="text-xs">{t('bills.list.active')}</Badge>
							{:else}
								<Badge variant="outline" class="text-xs">{t('bills.list.inactive')}</Badge>
							{/if}
							<button type="button" aria-label="{t('common.delete')}" class="text-muted-foreground hover:text-destructive transition-colors" onclick={() => (deleteTarget = bill.id)}>
								<Trash2 class="size-4" />
							</button>
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

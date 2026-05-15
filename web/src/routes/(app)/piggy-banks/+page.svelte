<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { Card, CardContent } from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import { Plus, ChevronDown, Trash2, PiggyBank as PiggyBankIcon } from '@lucide/svelte';
	import { piggyBankService, walletService } from '$lib/services/index.js';
	import { formatCurrency, formatPercentage } from '$lib/utils/format.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import type { PiggyBank, Account } from '$lib/types/domain.js';
	import { ConfirmDialog } from '$lib/components/shared/index.js';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import { Select, SelectTrigger, SelectContent, SelectItem } from '$lib/components/ui/select/index.js';
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	const t = localeStore.t;

	let isLoading = $state(true);
	let errorMsg = $state('');
	let items = $state<PiggyBank[]>([]);
	let wallets = $state<Account[]>([]);

	let accountFilter = $state('all');
	let deleteTarget = $state<{walletId: string; id: string} | null>(null);
	let deleteOpen = $derived(deleteTarget !== null);

	async function handleDelete(walletId: string, id: string) {
		await piggyBankService.delete(walletId, id);
		items = items.filter((pb) => pb.id !== id);
	}

	onMount(async () => {
		try {
			const walletList = await walletService.list();
			wallets = walletList;
			if (walletList.length > 0) {
				items = await piggyBankService.list(walletList[0].id);
			}
		} catch (e) {
			errorMsg = t('common.error');
			console.error('Failed to load piggy banks:', e);
		} finally {
			isLoading = false;
		}
	});

	let filtered = $derived(
		accountFilter === 'all'
			? items
			: items.filter((pb) => pb.account_id === accountFilter)
	);

	const accounts = $derived([...new Set(items.map((pb) => pb.account_id))]);
</script>

<div class="flex flex-col gap-4">
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
		<p class="text-sm text-destructive py-8 text-center">{errorMsg}</p>
	{:else}
	<div class="flex flex-wrap items-center justify-between gap-3">
		<div class="flex items-center gap-3">
			<h2 class="text-lg font-semibold text-foreground">{t('piggyBanks.list.title')}</h2>
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
		<SelectItem value="all">{t('piggyBanks.list.allWallets')}</SelectItem>
		{#each accounts as id}
{@const w = items.find((pb) => pb.account_id === id)}
<SelectItem value={id}>{w?.account_name ?? id}</SelectItem>
{/each}
		</SelectContent>
</Select>
				<ChevronDown class="pointer-events-none absolute right-2.5 top-1/2 size-3.5 -translate-y-1/2 text-muted-foreground" />
			</div>
			<Button size="sm" onclick={() => goto('/piggy-banks/create')}>
				<Plus class="size-4" />
				{t('piggyBanks.list.add')}
			</Button>
		</div>
	</div>

	<Card>
		<CardContent class="p-0">
			<p class="px-5 pt-4 text-sm font-semibold text-muted-foreground">{t('piggyBanks.list.active')}</p>
			{#each filtered as pb}
				{@const pct = parseFloat(pb.target_amount) > 0
					? (parseFloat(pb.current_amount) / parseFloat(pb.target_amount)) * 100
					: 0}
				<div class="flex items-center gap-4 px-5 py-4 border-b last:border-b-0">
					<div class="flex size-10 shrink-0 items-center justify-center rounded-lg bg-muted">
						<PiggyBankIcon class="size-5 text-foreground" />
					</div>
					<div class="flex flex-col gap-1 min-w-0">
						<p class="text-sm font-semibold text-foreground truncate">{pb.name}</p>
						<p class="text-xs text-muted-foreground">{t('piggyBanks.list.walletPrefix')} {pb.account_name}</p>
						<Progress value={Math.min(pct, 100)} class="mt-1 h-1.5 w-[120px]" />
					</div>
					<div class="ml-auto shrink-0 flex items-center gap-3">
						<div class="text-right">
							<p class="text-sm font-semibold text-foreground">{formatCurrency(pb.current_amount)}</p>
							<p class="text-xs text-muted-foreground">{t('piggyBanks.list.of', { pct: Math.round(pct), target: formatCurrency(pb.target_amount) })}</p>
						</div>
						<button type="button" aria-label="{t('common.delete')}" class="text-muted-foreground hover:text-destructive transition-colors" onclick={() => (deleteTarget = { walletId: pb.account_id, id: pb.id })}>
							<Trash2 class="size-4" />
						</button>
					</div>
				</div>
			{:else}
				<EmptyState />
			{/each}
		</CardContent>
	</Card>
	{/if}
	<ConfirmDialog
		bind:open={deleteOpen}
		title={t('common.delete')}
		description={t('common.deleteConfirm')}
		onConfirm={async () => {
			if (deleteTarget) {
				await handleDelete(deleteTarget.walletId, deleteTarget.id);
				deleteTarget = null;
			}
		}}
	/>
</div>

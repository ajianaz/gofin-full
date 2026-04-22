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
	const t = localeStore.t;

	let isLoading = $state(true);
	let items = $state<PiggyBank[]>([]);
	let wallets = $state<Account[]>([]);

	let accountFilter = $state('all');

	async function handleDelete(walletId: string, id: string) {
		if (!confirm('Hapus piggy bank ini?')) return;
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
		<p class="text-sm text-muted-foreground py-8 text-center">Memuat...</p>
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
				<select
					bind:value={accountFilter}
					class="cn-input h-9 w-44 appearance-none bg-background pr-8 text-sm"
				>
					<option value="all">{t('piggyBanks.list.allWallets')}</option>
					{#each accounts as id}
						{@const w = items.find((pb) => pb.account_id === id)}
						<option value={id}>{w?.account_name ?? id}</option>
					{/each}
				</select>
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
						<button type="button" class="text-muted-foreground hover:text-destructive transition-colors" onclick={() => handleDelete(pb.account_id, pb.id)}>
							<Trash2 class="size-4" />
						</button>
					</div>
				</div>
			{/each}
		</CardContent>
	</Card>
	{/if}
</div>

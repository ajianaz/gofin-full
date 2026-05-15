<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Card, CardContent } from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Plus, ChevronDown, Landmark, Smartphone, CreditCard, Wallet, Trash2 } from '@lucide/svelte';
	import { walletService } from '$lib/services/index.js';
	import type { Account } from '$lib/types/domain.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import { ConfirmDialog } from '$lib/components/shared/index.js';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import { Select, SelectTrigger, SelectContent, SelectItem } from '$lib/components/ui/select/index.js';
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	const t = localeStore.t;

	let items = $state<Account[]>([]);
	let isLoading = $state(true);
	let errorMsg = $state('');
	let typeFilter = $state('all');
	let deleteTarget = $state<string | null>(null);
	let deleteOpen = $derived(deleteTarget !== null);

	async function handleDelete(id: string) {
		await walletService.delete(id);
		items = items.filter((w) => w.id !== id);
	}

	onMount(async () => {
		try {
			items = await walletService.list();
		} catch (e) {
			errorMsg = t('common.error');
			console.error(e);
		} finally {
			isLoading = false;
		}
	});

	let filtered = $derived(
		typeFilter === 'all'
			? items
			: items.filter((w) => w.type === typeFilter)
	);

	function walletIcon(w: Account) {
		if (w.type === 'liability') return CreditCard;
		if (w.type === 'cash') return Wallet;
		const name = w.name.toLowerCase();
		if (name.includes('gopay') || name.includes('ovo') || name.includes('dana') || name.includes('shopeepay') || name.includes('ewallet')) return Smartphone;
		if (name.includes('credit') || name.includes('kartu')) return CreditCard;
		return Landmark;
	}

	function walletLabel(w: Account): string {
		if (w.type === 'liability') return t('wallets.list.creditCard');
		if (w.type === 'cash') return t('wallets.list.cash');
		const name = w.name.toLowerCase();
		if (name.includes('gopay') || name.includes('ovo') || name.includes('dana') || name.includes('shopeepay') || name.includes('ewallet')) return t('wallets.list.ewallet');
		if (name.includes('credit') || name.includes('kartu')) return t('wallets.list.creditCard');
		return t('wallets.list.bankAccount');
	}

	function formatBalance(balance: string): string {
		const num = Math.abs(parseFloat(balance));
		if (isNaN(num)) return 'Rp 0';
		return `Rp ${num.toLocaleString(localeStore.localeCode)}`;
	}
</script>

<div class="flex flex-col gap-4">
	<div class="flex flex-wrap items-center justify-between gap-3">
		<div class="flex items-center gap-3">
			<h2 class="text-base font-semibold text-foreground">{t('wallets.list.title')}</h2>
			<span class="inline-flex items-center rounded-2xl bg-primary px-2.5 py-0.5 text-xs font-medium text-primary-foreground">
				{filtered.length}
			</span>
		</div>
		<div class="flex items-center gap-3">
			<div class="relative">
				<Select bind:value={typeFilter}>
		<SelectTrigger class="w-40">
		</SelectTrigger>
		<SelectContent>
		<SelectItem value="all">{t('wallets.list.allTypes')}</SelectItem>
		<SelectItem value="asset">{t('wallets.list.bankAccount')}</SelectItem>
		<SelectItem value="cash">{t('wallets.list.cash')}</SelectItem>
		<SelectItem value="liability">{t('wallets.list.creditCard')}</SelectItem>
		<SelectItem value="expense">{t('wallets.list.ewallet')}</SelectItem>
		</SelectContent>
</Select>
				<ChevronDown class="pointer-events-none absolute right-2.5 top-1/2 size-3.5 -translate-y-1/2 text-muted-foreground" />
			</div>
			<Button size="sm" onclick={() => goto('/wallets/create')}>
				<Plus class="size-4" />
				{t('wallets.list.add')}
			</Button>
		</div>
	</div>

	<div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
		{#if isLoading}
{#each Array(6) as _}
				<div class="rounded-lg border bg-card p-5">
					<div class="flex items-center justify-between mb-4">
						<div class="flex items-center gap-2">
							<Skeleton class="size-[18px] rounded" />
							<Skeleton class="h-5 w-24" />
						</div>
						<Skeleton class="size-4 rounded" />
					</div>
					<Skeleton class="mb-1 h-7 w-32" />
					<Skeleton class="h-3 w-20" />
				</div>
			{/each}
	{:else if errorMsg}
			<p class="col-span-full text-sm text-destructive py-8 text-center">{errorMsg}</p>
		{:else}
			{#each filtered as wallet}
			{@const Icon = walletIcon(wallet)}
			<Card>
				<CardContent class="p-5">
					<div class="flex items-center justify-between mb-4">
						<div class="flex items-center gap-2">
							<Icon class="size-[18px] text-primary" />
							<span class="text-base font-semibold text-foreground">{wallet.name}</span>
						</div>
						<button type="button" aria-label="{t('common.delete')}" class="text-muted-foreground hover:text-destructive transition-colors" onclick={() => (deleteTarget = wallet.id)}>
							<Trash2 class="size-4" />
						</button>
					</div>
					<p class="text-xl font-bold {parseFloat(wallet.balance) < 0 ? 'text-red-600' : 'text-foreground'}">
						{formatBalance(wallet.balance)}
					</p>
					<p class="mt-1 text-xs text-muted-foreground">{walletLabel(wallet)}</p>
				</CardContent>
			</Card>
			{:else}
				<EmptyState />
			{/each}
		{/if}
	</div>
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

<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Card, CardContent } from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Plus, ChevronDown, Landmark, Smartphone, CreditCard, Wallet, Trash2, Pencil } from '@lucide/svelte';
	import { walletService } from '$lib/services/index.js';
	import { currencyService } from '$lib/services/currencies.js';
	import type { Account, Currency } from '$lib/types/domain.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import { getDefaultSymbol } from '$lib/utils/format.js';
	import { ConfirmDialog } from '$lib/components/shared/index.js';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import { Select, SelectTrigger, SelectContent, SelectItem, SelectValue } from '$lib/components/ui/select/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	const t = localeStore.t;

	let items = $state<Account[]>([]);
	let isLoading = $state(true);
	let errorMsg = $state('');
	let typeFilter = $state('all');
	let deleteTarget = $state<string | null>(null);
	let deleteOpen = $derived(deleteTarget !== null);

	// Edit state
	let editOpen = $state(false);
	let editWallet = $state<Account | null>(null);
	let editName = $state('');
	let editCurrencyId = $state('');
	let editIsSaving = $state(false);
	let editError = $state('');
	let currencies = $state<Currency[]>([]);

	async function loadCurrencies() {
		try {
			currencies = await currencyService.list();
		} catch (e) {
			console.error('Failed to load currencies:', e);
		}
	}

	function openEdit(wallet: Account) {
		editWallet = wallet;
		editName = wallet.name;
		editCurrencyId = wallet.currency_code || '';
		editError = '';
		editOpen = true;
	}

	async function handleEditSave() {
		if (!editWallet) return;
		editIsSaving = true;
		editError = '';
		try {
			const currencyId = currencies.find(c => c.code === editCurrencyId)?.id || editCurrencyId;
			await walletService.update(editWallet.id, {
				name: editName,
				currency_id: currencyId || undefined
			});
			// Refresh list
			items = await walletService.list();
			editOpen = false;
		} catch (err: any) {
			editError = err?.detail || err?.message || t('common.errorSave');
		} finally {
			editIsSaving = false;
		}
	}

	async function handleDelete(id: string) {
		await walletService.delete(id);
		items = items.filter((w) => w.id !== id);
	}

	onMount(async () => {
		try {
			[items] = await Promise.all([walletService.list(), loadCurrencies()]);
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

	function formatBalance(balance: string, symbol?: string): string {
		const num = Math.abs(parseFloat(balance));
		if (isNaN(num)) {
			const sym = symbol || getDefaultSymbol();
			return `${sym} ${(0).toLocaleString(localeStore.localeCode)}`;
		}
		const sym = symbol || getDefaultSymbol();
		return `${sym} ${num.toLocaleString(localeStore.localeCode)}`;
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
						<SelectValue />
					</SelectTrigger>
					<SelectContent>
						<SelectItem value="all" label={t('wallets.list.allTypes')}>{t('wallets.list.allTypes')}</SelectItem>
						<SelectItem value="asset" label={t('wallets.list.bankAccount')}>{t('wallets.list.bankAccount')}</SelectItem>
						<SelectItem value="cash" label={t('wallets.list.cash')}>{t('wallets.list.cash')}</SelectItem>
						<SelectItem value="liability" label={t('wallets.list.creditCard')}>{t('wallets.list.creditCard')}</SelectItem>
						<SelectItem value="expense" label={t('wallets.list.ewallet')}>{t('wallets.list.ewallet')}</SelectItem>
					</SelectContent>
				</Select>
			</div>
			<Button size="sm" onclick={() => goto('/wallets/create')}>
				<Plus class="size-4" />
				{t('wallets.list.add')}
			</Button>
		</div>
	</div>

	<Card>
		<CardContent class="p-0">
			{#if isLoading}
				{#each Array(6) as _}
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
			{:else if errorMsg}
				<p class="px-5 py-8 text-center text-sm text-destructive">{errorMsg}</p>
			{:else}
				{#each filtered as wallet}
					{@const Icon = walletIcon(wallet)}
					<div class="flex items-center gap-4 px-5 py-4 border-b last:border-b-0">
						<div class="flex size-10 shrink-0 items-center justify-center rounded-lg bg-muted">
							<Icon class="size-5 text-foreground" />
						</div>
						<div class="flex flex-col gap-1 min-w-0 flex-1">
							<p class="text-sm font-semibold text-foreground truncate">{wallet.name}</p>
							<p class="text-xs text-muted-foreground">{walletLabel(wallet)}</p>
						</div>
						<div class="ml-auto shrink-0 flex items-center gap-3">
							<div class="text-right">
								<p class="text-sm font-semibold {parseFloat(wallet.balance) < 0 ? 'text-red-600' : 'text-foreground'}">
									{formatBalance(wallet.balance, wallet.currency_symbol)}
								</p>
							</div>
							<Button variant="ghost" size="icon-sm" aria-label={t('wallets.edit.title')} class="text-muted-foreground hover:text-primary" onclick={() => openEdit(wallet)}>
								<Pencil class="size-4" />
							</Button>
							<Button variant="ghost" size="icon-sm" aria-label={t('common.delete')} class="text-muted-foreground hover:text-destructive" onclick={() => (deleteTarget = wallet.id)}>
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

	<!-- Edit Dialog -->
	{#if editOpen && editWallet}
		<!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
		<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" role="dialog" aria-modal="true" aria-label={t('wallets.edit.title')} onclick={() => (editOpen = false)}>
			<!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
			<div class="bg-card rounded-lg border p-6 w-full max-w-md mx-4 shadow-lg" role="document" onclick={(e) => e.stopPropagation()}>
				<h3 class="text-lg font-semibold text-foreground mb-4">{t('wallets.edit.title')}</h3>
				<div class="flex flex-col gap-4">
					<div class="flex flex-col gap-2">
						<Label for="edit-name">{t('wallets.edit.name')}</Label>
						<Input id="edit-name" bind:value={editName} />
					</div>
					<div class="flex flex-col gap-2">
						<Label for="edit-currency">{t('wallets.edit.currency')}</Label>
						<Select bind:value={editCurrencyId} id="edit-currency">
							<SelectTrigger class="w-full"><SelectValue /></SelectTrigger>
							<SelectContent>
								{#each currencies as c}
									<SelectItem value={c.code} label={c.code} ({c.symbol}) — {c.name}>{c.code} ({c.symbol}) — {c.name}</SelectItem>
								{/each}
							</SelectContent>
						</Select>
					</div>
					{#if editError}
						<p class="text-destructive text-sm">{editError}</p>
					{/if}
					<div class="flex gap-2 mt-2">
						<Button class="flex-1" onclick={handleEditSave} disabled={editIsSaving}>{editIsSaving ? t('common.saving') : t('wallets.edit.save')}</Button>
						<Button variant="outline" class="flex-1" onclick={() => (editOpen = false)}>{t('common.cancel')}</Button>
					</div>
				</div>
			</div>
		</div>
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

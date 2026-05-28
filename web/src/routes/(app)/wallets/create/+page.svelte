<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { ChevronDown } from '@lucide/svelte';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import { walletService, currencyService } from '$lib/services/index.js';
	import { Select, SelectTrigger, SelectContent, SelectItem, SelectValue } from '$lib/components/ui/select/index.js';
	import FormCard from '$lib/components/shared/FormCard.svelte';
	import type { Currency } from '$lib/types/domain.js';
	const t = localeStore.t;

	let isLoading = $state(false);
	let errorMsg = $state('');
	let name = $state('');
	let type = $state('asset');
	let currencyId = $state('');
	let openingBalance = $state('');
	let currencies = $state<Currency[]>([]);
	let loadingCurrencies = $state(true);

	onMount(async () => {
		try {
			currencies = await currencyService.list();
			// Auto-select first currency if available
			if (currencies.length > 0) {
				currencyId = currencies[0].id;
			}
		} catch (e) {
			console.error('Failed to load currencies:', e);
		} finally {
			loadingCurrencies = false;
		}
	});
</script>

<div class="flex flex-col gap-4">
	<FormCard title={t('wallets.create.title')}>
		<form class="flex flex-col gap-6" onsubmit={async (e) => { e.preventDefault(); isLoading = true; errorMsg = ''; try { await walletService.create({ name, wallet_type: type, currency_id: currencyId, opening_balance: openingBalance ? String(openingBalance) : undefined }); goto('/wallets'); } catch (err: any) { errorMsg = err?.detail || err?.message || t('common.error'); } finally { isLoading = false; } }}>
			<div class="grid gap-6 md:grid-cols-2">
				<div class="flex flex-col gap-2">
					<Label for="name">{t('wallets.create.name')}</Label>
					<Input id="name" placeholder={t('wallets.create.namePlaceholder')} bind:value={name} required />
				</div>

				<div class="flex flex-col gap-2">
					<Label for="type">{t('wallets.create.type')}</Label>
					<div class="relative">
						<Select bind:value={type} id="type" items={[
							{ value: 'asset', label: t('wallets.create.bankAccount') },
							{ value: 'cash', label: t('wallets.create.cash') },
							{ value: 'liability', label: t('wallets.create.creditCard') },
							{ value: 'expense', label: t('wallets.create.ewallet') },
							{ value: 'revenue', label: t('wallets.create.investment') },
						]}>
<SelectTrigger class="w-full">
											<SelectValue />
										</SelectTrigger>
										<SelectContent>
											<SelectItem value="asset" label={t('wallets.create.bankAccount')}>{t('wallets.create.bankAccount')}</SelectItem>
											<SelectItem value="cash" label={t('wallets.create.cash')}>{t('wallets.create.cash')}</SelectItem>
											<SelectItem value="liability" label={t('wallets.create.creditCard')}>{t('wallets.create.creditCard')}</SelectItem>
											<SelectItem value="expense" label={t('wallets.create.ewallet')}>{t('wallets.create.ewallet')}</SelectItem>
											<SelectItem value="revenue" label={t('wallets.create.investment')}>{t('wallets.create.investment')}</SelectItem>
										</SelectContent>
</Select>
											</div>
										</div>

				<div class="flex flex-col gap-2">
					<Label for="currency">{t('wallets.create.currency')}</Label>
					<div class="relative">
						{#if loadingCurrencies}
							<Input disabled value={t('common.loading')} />
						{:else}
							<Select bind:value={currencyId} id="currency" items={[
								...currencies.map(c => ({ value: c.id, label: `${c.code} (${c.symbol}) — ${c.name}` }))
							]}>
								<SelectTrigger class="w-full"><SelectValue /></SelectTrigger>
								<SelectContent>
									{#each currencies as c}
										<SelectItem value={c.id} label={`${c.code} (${c.symbol}) — ${c.name}`}>{c.code} ({c.symbol}) — {c.name}</SelectItem>
									{/each}
								</SelectContent>
</Select>
										{/if}
					</div>
				</div>

				<div class="flex flex-col gap-2">
					<Label for="opening-balance">{t('wallets.create.balance')}</Label>
					<Input id="opening-balance" type="number" step="0.01" placeholder="0" bind:value={openingBalance} />
				</div>
			</div>

			{#if errorMsg}
				<p class="text-destructive text-sm">{errorMsg}</p>
			{/if}

			<div class="flex gap-2">
				<Button type="submit" class="flex-1" disabled={isLoading || !currencyId}>{isLoading ? t('common.saving') : t('wallets.create.submit')}</Button>
				<Button type="button" variant="outline" class="flex-1" onclick={() => goto('/wallets')}>{t('common.cancel')}</Button>
			</div>
		</form>
	</FormCard>
</div>

<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { mockCurrencies } from '$lib/data/mock-currencies.js';
	import { ChevronDown } from '@lucide/svelte';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import { walletService } from '$lib/services/index.js';
	const t = localeStore.t;

	let isLoading = $state(false);
	let errorMsg = $state('');
	let name = $state('');
	let type = $state('');
	let currencyCode = $state('IDR');
	let balance = $state('');
	let virtualBalance = $state('');
	let iban = $state('');
	let includeNetWorth = $state(true);
</script>

<div class="flex flex-col gap-4">
	<Card>
		<CardHeader>
			<CardTitle>{t('wallets.create.title')}</CardTitle>
		</CardHeader>
		<CardContent>
			<form class="flex flex-col gap-6" onsubmit={async (e) => { e.preventDefault(); isLoading = true; errorMsg = ''; try { await walletService.create({ name, wallet_type: type }); goto('/wallets'); } catch (err: any) { errorMsg = err?.detail || err?.message || 'Gagal menyimpan'; } finally { isLoading = false; } }}>
				<div class="grid gap-6 md:grid-cols-2">
					<div class="flex flex-col gap-4">
						<div class="flex flex-col gap-2">
							<Label for="name">{t('wallets.create.name')}</Label>
							<Input id="name" placeholder={t('wallets.create.namePlaceholder')} bind:value={name} required />
						</div>

						<div class="flex flex-col gap-2">
							<Label for="currency">{t('wallets.create.currency')}</Label>
							<div class="relative">
								<select
									id="currency"
									bind:value={currencyCode}
									class="cn-input w-full appearance-none bg-background pr-8"
								>
									{#each mockCurrencies.filter((c) => c.enabled) as c}
										<option value={c.code}>{c.code} ({c.symbol})</option>
									{/each}
								</select>
								<ChevronDown class="pointer-events-none absolute right-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
							</div>
						</div>

						<div class="flex flex-col gap-2">
							<Label for="iban">{t('wallets.create.iban')}</Label>
							<Input id="iban" placeholder={t('wallets.create.ibanPlaceholder')} bind:value={iban} />
						</div>
					</div>

					<div class="flex flex-col gap-4">
						<div class="flex flex-col gap-2">
							<Label for="type">{t('wallets.create.type')}</Label>
							<div class="relative">
								<select
									id="type"
									bind:value={type}
									class="cn-input w-full appearance-none bg-background pr-8"
								>
									<option value="">{t('common.all')}</option>
									<option value="asset">{t('wallets.create.bankAccount')}</option>
									<option value="cash">{t('wallets.create.cash')}</option>
									<option value="liability">{t('wallets.create.creditCard')}</option>
									<option value="expense">{t('wallets.create.ewallet')}</option>
									<option value="revenue">{t('wallets.create.investment')}</option>
								</select>
								<ChevronDown class="pointer-events-none absolute right-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
							</div>
						</div>

						<div class="flex flex-col gap-2">
							<Label for="balance">{t('wallets.create.balance')}</Label>
							<Input id="balance" type="number" placeholder="0" bind:value={balance} />
						</div>

						<div class="flex flex-col gap-2">
							<Label for="virtualBalance">{t('wallets.create.virtualBalance')}</Label>
							<Input id="virtualBalance" type="number" placeholder="0" bind:value={virtualBalance} />
						</div>
					</div>
				</div>

				<div class="flex items-center gap-2">
					<Checkbox id="netWorth" bind:checked={includeNetWorth} />
					<Label for="netWorth">{t('wallets.create.includeNetWorth')}</Label>
				</div>

				{#if errorMsg}
					<p class="text-destructive text-sm">{errorMsg}</p>
				{/if}

				<div class="flex gap-2">
					<Button type="submit" class="flex-1" disabled={isLoading}>{isLoading ? 'Saving...' : t('wallets.create.submit')}</Button>
					<Button type="button" variant="outline" class="flex-1" onclick={() => goto('/wallets')}>{t('common.cancel')}</Button>
				</div>
			</form>
		</CardContent>
	</Card>
</div>

<script lang="ts">
	import { goto } from '$app/navigation';
	import { BackButton, FormCard } from '$lib/components/shared/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { ChevronDown } from '@lucide/svelte';
	import { piggyBankService, walletService } from '$lib/services/index.js';
	import { onMount } from 'svelte';
	import type { Account } from '$lib/types/domain.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	let name = $state('');
	let accountId = $state('');
	let targetAmount = $state('');
	let initBalance = $state('');
	let startDate = $state(new Date().toISOString().split('T')[0]);
	let targetDate = $state('');
	let isLoading = $state(false);
	let errorMsg = $state('');
	let wallets = $state<Account[]>([]);

	onMount(async () => {
		try { wallets = await walletService.list(); } catch {}
	});
</script>

<BackButton href="/piggy-banks" />

<FormCard title={t('piggyBanks.create.title')}>
	<form class="flex flex-col gap-6" onsubmit={async (e) => {
			e.preventDefault();
			isLoading = true;
			errorMsg = '';
			try {
				await piggyBankService.create({ wallet_id: accountId, name, target_amount: targetAmount });
				goto('/piggy-banks');
			} catch (err: any) {
				errorMsg = err?.detail || err?.message || 'Gagal menyimpan';
			} finally {
				isLoading = false;
			}
		}}>
		<div class="grid gap-6 md:grid-cols-2">
			<div class="flex flex-col gap-4">
				<div class="flex flex-col gap-2">
					<Label for="name">{t('piggyBanks.create.name')}</Label>
					<Input id="name" placeholder={t('piggyBanks.create.namePlaceholder')} bind:value={name} required />
				</div>
				<div class="flex flex-col gap-2">
					<Label for="target">{t('piggyBanks.create.targetAmount')}</Label>
					<Input id="target" type="number" placeholder="0" bind:value={targetAmount} required />
				</div>
				<div class="flex flex-col gap-2">
					<Label for="start">{t('piggyBanks.create.startDate')}</Label>
					<Input id="start" type="date" bind:value={startDate} required />
				</div>
			</div>
			<div class="flex flex-col gap-4">
				<div class="flex flex-col gap-2">
					<Label for="account">{t('piggyBanks.create.relatedWallet')}</Label>
					<div class="relative">
						<select
							id="account"
							bind:value={accountId}
							class="cn-input w-full appearance-none bg-background pr-8"
						>
							<option value="">{t('common.selectWallet')}</option>
							{#each wallets as w}
								<option value={w.id}>{w.name}</option>
							{/each}
						</select>
						<ChevronDown class="pointer-events-none absolute right-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
					</div>
				</div>
				<div class="flex flex-col gap-2">
					<Label for="init">{t('piggyBanks.create.initialBalance')}</Label>
					<Input id="init" type="number" placeholder="0" bind:value={initBalance} />
				</div>
				<div class="flex flex-col gap-2">
					<Label for="target-date">{t('piggyBanks.create.targetDate')}</Label>
					<Input id="target-date" type="date" bind:value={targetDate} />
				</div>
			</div>
		</div>
		<div class="flex gap-2">
			{#if errorMsg}
				<p class="text-destructive text-sm">{errorMsg}</p>
			{/if}
			<Button type="submit" class="flex-1" disabled={isLoading}>{isLoading ? 'Menyimpan...' : t('piggyBanks.create.submit')}</Button>
			<Button type="button" variant="outline" class="flex-1" onclick={() => goto('/piggy-banks')}>{t('common.cancel')}</Button>
		</div>
	</form>
</FormCard>

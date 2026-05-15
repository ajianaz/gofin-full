<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { ChevronDown } from '@lucide/svelte';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import { walletService } from '$lib/services/index.js';
	import { Select, SelectTrigger, SelectContent, SelectItem } from '$lib/components/ui/select/index.js';
	import FormCard from '$lib/components/shared/FormCard.svelte';
	const t = localeStore.t;

	let isLoading = $state(false);
	let errorMsg = $state('');
	let name = $state('');
	let type = $state('asset');
</script>

<div class="flex flex-col gap-4">
	<FormCard title="{t('wallets.create.title')}">
			<form class="flex flex-col gap-6" onsubmit={async (e) => { e.preventDefault(); isLoading = true; errorMsg = ''; try { await walletService.create({ name, wallet_type: type }); goto('/wallets'); } catch (err: any) { errorMsg = err?.detail || err?.message || t('common.error'); } finally { isLoading = false; } }}>
				<div class="grid gap-6 md:grid-cols-2">
					<div class="flex flex-col gap-2">
						<Label for="name">{t('wallets.create.name')}</Label>
						<Input id="name" placeholder={t('wallets.create.namePlaceholder')} bind:value={name} required />
					</div>

					<div class="flex flex-col gap-2">
						<Label for="type">{t('wallets.create.type')}</Label>
						<div class="relative">
							<Select bind:value={type} id="type">
		<SelectTrigger class="w-full">
		</SelectTrigger>
		<SelectContent>
		<SelectItem value="asset">{t('wallets.create.bankAccount')}</SelectItem>
		<SelectItem value="cash">{t('wallets.create.cash')}</SelectItem>
		<SelectItem value="liability">{t('wallets.create.creditCard')}</SelectItem>
		<SelectItem value="expense">{t('wallets.create.ewallet')}</SelectItem>
		<SelectItem value="revenue">{t('wallets.create.investment')}</SelectItem>
		</SelectContent>
</Select>
							<ChevronDown class="pointer-events-none absolute right-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
						</div>
					</div>
				</div>

				{#if errorMsg}
					<p class="text-destructive text-sm">{errorMsg}</p>
				{/if}

				<div class="flex gap-2">
					<Button type="submit" class="flex-1" disabled={isLoading}>{isLoading ? t('common.saving') : t('wallets.create.submit')}</Button>
					<Button type="button" variant="outline" class="flex-1" onclick={() => goto('/wallets')}>{t('common.cancel')}</Button>
				</div>
			</form>
		</FormCard>
</div>

<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { ChevronDown } from '@lucide/svelte';
	import { transactionService, walletService, categoryService, tagService } from '$lib/services/index.js';
	import type { Account } from '$lib/types/domain.js';
	import type { Category } from '$lib/types/domain.js';
	import type { Tag } from '$lib/types/domain.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import { Select, SelectTrigger, SelectContent, SelectItem, SelectValue } from '$lib/components/ui/select/index.js';
	import FormCard from '$lib/components/shared/FormCard.svelte';
	const t = localeStore.t;

	let isLoading = $state(false);
	let errorMsg = $state('');
	let wallets = $state<Account[]>([]);
	let categories = $state<Category[]>([]);
	let tags = $state<Tag[]>([]);

	let type = $state('withdrawal');
	let date = $state(new Date().toISOString().split('T')[0]);
	let description = $state('');
	let amount = $state('');
	let sourceAccount = $state('');
	let destAccount = $state('');
	let category = $state('');
	let note = $state('');

	onMount(async () => {
		try {
			wallets = await walletService.list();
			categories = await categoryService.list();
			tags = await tagService.list();
		} catch (e) { console.error(e); }
	});

	async function handleSubmit(e: Event) {
		e.preventDefault();
		isLoading = true;
		errorMsg = '';
		try {
			const payload = {
				type,
				description,
				amount: String(amount),
				source_id: sourceAccount,
				destination_id: type === 'transfer' ? destAccount : sourceAccount,
				date,
				category_ids: category ? [category] : [],
			};
			await transactionService.create(payload);
			goto('/transactions');
		} catch (err: any) {
			errorMsg = err.detail || err.message || t('common.errorSave');
		} finally {
			isLoading = false;
		}
	}
</script>

<div class="flex flex-col gap-4">
	<FormCard title={t('transactions.create.title')}>
			<form class="flex flex-col gap-6" onsubmit={handleSubmit}>
				<div class="grid gap-6 md:grid-cols-2">
					<div class="flex flex-col gap-4">
						<div class="flex flex-col gap-2">
							<Label for="type">{t('transactions.create.type')}</Label>
							<div class="relative">
						<Select bind:value={type} id="type" items={[
							{ value: 'withdrawal', label: t('transactions.create.expense') },
							{ value: 'deposit', label: t('transactions.create.income') },
							{ value: 'transfer', label: t('transactions.create.transfer') },
						]}>
<SelectTrigger class="w-full">
	<SelectValue />
</SelectTrigger>
	<SelectContent>
	<SelectItem value="withdrawal" label={t('transactions.create.expense')}>{t('transactions.create.expense')}</SelectItem>
	<SelectItem value="deposit" label={t('transactions.create.income')}>{t('transactions.create.income')}</SelectItem>
	<SelectItem value="transfer" label={t('transactions.create.transfer')}>{t('transactions.create.transfer')}</SelectItem>
	</SelectContent>
</Select>
							</div>
						</div>

						<div class="flex flex-col gap-2">
							<Label for="date">{t('transactions.create.date')}</Label>
							<Input id="date" type="date" bind:value={date} required />
						</div>

						<div class="flex flex-col gap-2">
							<Label for="source">{t('transactions.create.sourceWallet')}</Label>
							<div class="relative">
						<Select bind:value={sourceAccount} id="source" items={[
							{ value: '', label: t('common.selectWallet') },
							...wallets.map(w => ({ value: w.id, label: `${w.name} (${w.currency_code})` }))
						]}>
<SelectTrigger class="w-full">
	<SelectValue />
</SelectTrigger>
	<SelectContent>
	<SelectItem value="" label={t('common.selectWallet')}>{t('common.selectWallet')}</SelectItem>
	{#each wallets as w}
<SelectItem value={w.id} label={'{w.name} ({w.currency_code})'}>{w.name} ({w.currency_code})</SelectItem>
{/each}
	</SelectContent>
</Select>
							</div>
						</div>

						{#if type === 'transfer'}
							<div class="flex flex-col gap-2">
								<Label for="dest">{t('transactions.create.destWallet')}</Label>
								<div class="relative">
							<Select bind:value={destAccount} id="dest" items={[
								{ value: '', label: t('common.selectWallet') },
								...wallets.map(w => ({ value: w.id, label: `${w.name} (${w.currency_code})` }))
							]}>
<SelectTrigger class="w-full">
	<SelectValue />
</SelectTrigger>
	<SelectContent>
	<SelectItem value="" label={t('common.selectWallet')}>{t('common.selectWallet')}</SelectItem>
	{#each wallets as w}
<SelectItem value={w.id} label={'{w.name} ({w.currency_code})'}>{w.name} ({w.currency_code})</SelectItem>
{/each}
	</SelectContent>
</Select>
								</div>
							</div>
						{/if}

						<div class="flex flex-col gap-2">
							<Label for="category">{t('transactions.create.category')}</Label>
							<div class="relative">
						<Select bind:value={category} id="category" items={[
							{ value: '', label: t('common.selectCategory') },
							...categories.map(cat => ({ value: cat.id, label: cat.name }))
						]}>
<SelectTrigger class="w-full">
	<SelectValue />
</SelectTrigger>
	<SelectContent>
	<SelectItem value="" label={t('common.selectCategory')}>{t('common.selectCategory')}</SelectItem>
	{#each categories as cat}
<SelectItem value={cat.id} label={cat.name}>{cat.name}</SelectItem>
{/each}
	</SelectContent>
</Select>
							</div>
						</div>
					</div>

					<div class="flex flex-col gap-4">
						<div class="flex flex-col gap-2">
							<Label for="description">{t('transactions.create.description')}</Label>
							<Input id="description" placeholder={t('transactions.create.descriptionPlaceholder')} bind:value={description} required />
						</div>

						<div class="flex flex-col gap-2">
							<Label for="amount">{t('transactions.create.amount')}</Label>
							<Input id="amount" type="number" placeholder="0" bind:value={amount} required />
						</div>

						<div class="flex flex-col gap-2">
							<Label for="note">{t('transactions.create.note')}</Label>
							<Input id="note" placeholder={t('transactions.create.notePlaceholder')} bind:value={note} />
						</div>

						<div class="flex flex-col gap-2">
							<Label for="tags">{t('transactions.create.tag')}</Label>
						<Select value="" items={[
							...tags.map(tag => ({ value: tag.tag, label: tag.tag }))
						]}>
<SelectTrigger class="w-full">
								<SelectValue />
							</SelectTrigger>
							<SelectContent>
								{#each tags as tag}
								<SelectItem value={tag.tag} label={tag.tag}>{tag.tag}</SelectItem>
								{/each}
							</SelectContent>
						</Select>
						</div>
					</div>
				</div>

				<div class="flex gap-2">
					{#if errorMsg}
						<p class="text-sm text-destructive">{errorMsg}</p>
					{/if}
					<Button type="submit" class="flex-1" disabled={isLoading}>{isLoading ? t('common.saving') : t('transactions.create.submit')}</Button>
					<Button type="button" variant="outline" class="flex-1" onclick={() => goto('/transactions')}>{t('common.cancel')}</Button>
				</div>
			</form>
		</FormCard>
</div>

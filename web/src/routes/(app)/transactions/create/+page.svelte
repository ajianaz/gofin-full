<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { ChevronDown } from '@lucide/svelte';
	import { transactionService, walletService, categoryService, tagService } from '$lib/services/index.js';
	import type { Account } from '$lib/types/domain.js';
	import type { Category } from '$lib/types/domain.js';
	import type { Tag } from '$lib/types/domain.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
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
	<Card>
		<CardHeader>
			<CardTitle>{t('transactions.create.title')}</CardTitle>
		</CardHeader>
		<CardContent>
			<form class="flex flex-col gap-6" onsubmit={handleSubmit}>
				<div class="grid gap-6 md:grid-cols-2">
					<div class="flex flex-col gap-4">
						<div class="flex flex-col gap-2">
							<Label for="type">{t('transactions.create.type')}</Label>
							<div class="relative">
								<select
									id="type"
									bind:value={type}
									class="cn-input w-full appearance-none bg-background pr-8"
								>
									<option value="withdrawal">{t('transactions.create.expense')}</option>
									<option value="deposit">{t('transactions.create.income')}</option>
									<option value="transfer">{t('transactions.create.transfer')}</option>
								</select>
								<ChevronDown class="pointer-events-none absolute right-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
							</div>
						</div>

						<div class="flex flex-col gap-2">
							<Label for="date">{t('transactions.create.date')}</Label>
							<Input id="date" type="date" bind:value={date} required />
						</div>

						<div class="flex flex-col gap-2">
							<Label for="source">{t('transactions.create.sourceWallet')}</Label>
							<div class="relative">
								<select
									id="source"
									bind:value={sourceAccount}
									class="cn-input w-full appearance-none bg-background pr-8"
								>
									<option value="">{t('common.selectWallet')}</option>
									{#each wallets as w}
										<option value={w.id}>{w.name} ({w.currency_code})</option>
									{/each}
								</select>
								<ChevronDown class="pointer-events-none absolute right-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
							</div>
						</div>

						{#if type === 'transfer'}
							<div class="flex flex-col gap-2">
								<Label for="dest">{t('transactions.create.destWallet')}</Label>
								<div class="relative">
									<select
										id="dest"
										bind:value={destAccount}
										class="cn-input w-full appearance-none bg-background pr-8"
									>
										<option value="">{t('common.selectWallet')}</option>
										{#each wallets as w}
											<option value={w.id}>{w.name} ({w.currency_code})</option>
										{/each}
									</select>
									<ChevronDown class="pointer-events-none absolute right-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
								</div>
							</div>
						{/if}

						<div class="flex flex-col gap-2">
							<Label for="category">{t('transactions.create.category')}</Label>
							<div class="relative">
								<select
									id="category"
									bind:value={category}
									class="cn-input w-full appearance-none bg-background pr-8"
								>
									<option value="">{t('common.selectCategory')}</option>
									{#each categories as cat}
										<option value={cat.id}>{cat.name}</option>
									{/each}
								</select>
								<ChevronDown class="pointer-events-none absolute right-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
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
							<div class="relative">
								<select
									id="tags"
									class="cn-input w-full appearance-none bg-background pr-8"
								>
									<option value="">{t('transactions.create.selectTag')}</option>
									{#each tags as tag}
										<option value={tag.tag}>{tag.tag}</option>
									{/each}
								</select>
								<ChevronDown class="pointer-events-none absolute right-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
							</div>
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
		</CardContent>
	</Card>
</div>

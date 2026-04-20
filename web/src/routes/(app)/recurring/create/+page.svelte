<script lang="ts">
	import { goto } from '$app/navigation';
	import { BackButton, FormCard } from '$lib/components/shared/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { ChevronDown } from '@lucide/svelte';
	import { mockWallets } from '$lib/data/mock-wallets.js';
	import { mockCategories } from '$lib/data/mock-categories.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	let title = $state('');
	let repeatFreq = $state('monthly');
	let sourceAccount = $state('');
	let destAccount = $state('');
	let amount = $state('');
	let categoryId = $state('');
	let startDate = $state(new Date().toISOString().split('T')[0]);
	let endDate = $state('');
	let active = $state(true);
</script>

<BackButton href="/recurring" />

<FormCard title={t('recurring.create.title')}>
	<form class="flex flex-col gap-6" onsubmit={(e) => { e.preventDefault(); goto('/recurring'); }}>
		<div class="grid gap-6 md:grid-cols-2">
			<div class="flex flex-col gap-4">
				<div class="flex flex-col gap-2">
					<Label for="title">{t('recurring.create.titleField')}</Label>
					<Input id="title" placeholder={t('recurring.create.titlePlaceholder')} bind:value={title} required />
				</div>
				<div class="flex flex-col gap-2">
					<Label for="source">{t('recurring.create.sourceWallet')}</Label>
					<div class="relative">
						<select
							id="source"
							bind:value={sourceAccount}
							class="cn-input w-full appearance-none bg-background pr-8"
						>
							<option value="">{t('common.selectWallet')}</option>
							{#each mockWallets as w}
								<option value={w.id}>{w.name}</option>
							{/each}
						</select>
						<ChevronDown class="pointer-events-none absolute right-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
					</div>
				</div>
				<div class="flex flex-col gap-2">
					<Label for="amount">{t('recurring.create.amount')}</Label>
					<Input id="amount" type="number" placeholder="0" bind:value={amount} required />
				</div>
				<div class="flex flex-col gap-2">
					<Label for="start">{t('recurring.create.startDate')}</Label>
					<Input id="start" type="date" bind:value={startDate} required />
				</div>
			</div>
			<div class="flex flex-col gap-4">
				<div class="flex flex-col gap-2">
					<Label for="repeat">{t('recurring.create.repeatType')}</Label>
					<div class="relative">
						<select
							id="repeat"
							bind:value={repeatFreq}
							class="cn-input w-full appearance-none bg-background pr-8"
						>
							<option value="daily">{t('recurring.list.freqDaily')}</option>
							<option value="weekly">{t('recurring.list.freqWeekly')}</option>
							<option value="monthly">{t('recurring.list.freqMonthly')}</option>
							<option value="quarterly">{t('recurring.list.freqQuarterly')}</option>
							<option value="yearly">{t('recurring.list.freqYearly')}</option>
						</select>
						<ChevronDown class="pointer-events-none absolute right-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
					</div>
				</div>
				<div class="flex flex-col gap-2">
					<Label for="dest">{t('recurring.create.destWallet')}</Label>
					<div class="relative">
						<select
							id="dest"
							bind:value={destAccount}
							class="cn-input w-full appearance-none bg-background pr-8"
						>
							<option value="">{t('recurring.create.selectDestWallet')}</option>
							{#each mockWallets as w}
								<option value={w.id}>{w.name}</option>
							{/each}
						</select>
						<ChevronDown class="pointer-events-none absolute right-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
					</div>
				</div>
				<div class="flex flex-col gap-2">
					<Label for="category">{t('recurring.create.category')}</Label>
					<div class="relative">
						<select
							id="category"
							bind:value={categoryId}
							class="cn-input w-full appearance-none bg-background pr-8"
						>
							<option value="">{t('common.selectCategory')}</option>
							{#each mockCategories as cat}
								<option value={cat.id}>{cat.name}</option>
							{/each}
						</select>
						<ChevronDown class="pointer-events-none absolute right-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
					</div>
				</div>
				<div class="flex flex-col gap-2">
					<Label for="end">{t('recurring.create.endDate')}</Label>
					<Input id="end" type="date" placeholder={t('recurring.create.endDatePlaceholder')} bind:value={endDate} />
				</div>
			</div>
		</div>
		<div class="flex items-center gap-2">
			<Checkbox id="active" bind:checked={active} />
			<Label for="active">{t('recurring.create.active')}</Label>
		</div>
		<div class="flex gap-2 pt-2">
			<Button type="submit" class="flex-1">{t('recurring.create.submit')}</Button>
			<Button type="button" variant="outline" class="flex-1" onclick={() => goto('/recurring')}>{t('common.cancel')}</Button>
		</div>
	</form>
</FormCard>

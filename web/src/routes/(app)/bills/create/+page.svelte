<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { ChevronDown } from '@lucide/svelte';
	import { mockWallets } from '$lib/data/mock-wallets.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	let name = $state('');
	let accountId = $state('');
	let amountMin = $state('');
	let amountMax = $state('');
	let startDate = $state(new Date().toISOString().split('T')[0]);
	let repeatFreq = $state('monthly');
	let endDate = $state('');
	let active = $state(true);
</script>

<div class="flex flex-col gap-4">
	<Card>
		<CardHeader>
			<CardTitle>{t('bills.create.title')}</CardTitle>
		</CardHeader>
		<CardContent>
			<form class="flex flex-col gap-6" onsubmit={(e) => { e.preventDefault(); goto('/bills'); }}>
				<div class="grid gap-6 md:grid-cols-2">
					<div class="flex flex-col gap-4">
						<div class="flex flex-col gap-2">
							<Label for="name">{t('bills.create.name')}</Label>
							<Input id="name" placeholder={t('bills.create.namePlaceholder')} bind:value={name} required />
						</div>
						<div class="flex flex-col gap-2">
							<Label for="min">{t('bills.create.minAmount')}</Label>
							<Input id="min" type="number" placeholder="0" bind:value={amountMin} required />
						</div>
						<div class="flex flex-col gap-2">
							<Label for="start">{t('bills.create.startDate')}</Label>
							<Input id="start" type="date" bind:value={startDate} required />
						</div>
						<div class="flex flex-col gap-2">
							<Label for="end">{t('bills.create.endDate')}</Label>
							<Input id="end" type="date" bind:value={endDate} />
						</div>
					</div>
					<div class="flex flex-col gap-4">
						<div class="flex flex-col gap-2">
							<Label for="account">{t('bills.create.relatedWallet')}</Label>
							<div class="relative">
								<select
									id="account"
									bind:value={accountId}
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
							<Label for="max">{t('bills.create.maxAmount')}</Label>
							<Input id="max" type="number" placeholder="0" bind:value={amountMax} />
						</div>
						<div class="flex flex-col gap-2">
							<Label for="freq">{t('bills.create.frequency')}</Label>
							<div class="relative">
								<select
									id="freq"
									bind:value={repeatFreq}
									class="cn-input w-full appearance-none bg-background pr-8"
								>
									<option value="weekly">{t('bills.create.freqWeekly')}</option>
									<option value="monthly">{t('bills.create.freqMonthly')}</option>
									<option value="quarterly">{t('bills.create.freqQuarterly')}</option>
									<option value="yearly">{t('bills.create.freqYearly')}</option>
								</select>
								<ChevronDown class="pointer-events-none absolute right-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
							</div>
						</div>
						<div class="flex items-center gap-2 pt-2">
							<Checkbox id="active" bind:checked={active} />
							<Label for="active">{t('bills.create.active')}</Label>
						</div>
					</div>
				</div>
				<div class="flex gap-2">
					<Button type="submit" class="flex-1">{t('bills.create.submit')}</Button>
					<Button type="button" variant="outline" class="flex-1" onclick={() => goto('/bills')}>{t('common.cancel')}</Button>
				</div>
			</form>
		</CardContent>
	</Card>
</div>

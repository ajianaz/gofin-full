<script lang="ts">
	import { BackButton } from '$lib/components/shared/index.js';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '$lib/components/ui/table/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Select, SelectTrigger, SelectContent, SelectItem } from '$lib/components/ui/select/index.js';
	import { currencyService } from '$lib/services/index.js';
	import { formatDate } from '$lib/utils/format.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import { onMount } from 'svelte';
	import type { ExchangeRate, Currency } from '$lib/types/domain.js';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import { ConfirmDialog } from '$lib/components/shared/index.js';
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	import { Plus, Trash2, ChevronDown } from '@lucide/svelte';

	const t = localeStore.t;

	let rates: ExchangeRate[] = $state([]);
	let currencies: Currency[] = $state([]);
	let isLoading: boolean = $state(true);
	let error: string | null = $state(null);

	// Add form state
	let showAddForm = $state(false);
	let addFrom = $state('');
	let addTo = $state('');
	let addRate = $state('');
	let addDate = $state('');
	let addLoading = $state(false);
	let addError = $state('');

	// Delete state
	let deleteTarget = $state<string | null>(null);
	let deleteOpen = $derived(deleteTarget !== null);

	async function loadData() {
		try {
			const [ratesRes, currenciesRes] = await Promise.all([
				currencyService.exchangeRates(),
				currencyService.list()
			]);
			rates = ratesRes;
			currencies = currenciesRes;
		} catch (e: any) {
			error = e?.message ?? t('common.error');
		} finally {
			isLoading = false;
		}
	}

	onMount(loadData);

	async function handleAdd(e: Event) {
		e.preventDefault();
		addLoading = true;
		addError = '';
		try {
			const newRate = await currencyService.createExchangeRate({
				from_currency_id: addFrom,
				to_currency_id: addTo,
				rate: addRate,
				date: addDate || undefined
			});
			rates.unshift(newRate);
			showAddForm = false;
			addFrom = '';
			addTo = '';
			addRate = '';
			addDate = '';
		} catch (err: any) {
			addError = err?.detail || err?.message || t('common.error');
		} finally {
			addLoading = false;
		}
	}

	async function handleDelete(id: string) {
		await currencyService.deleteExchangeRate(id);
		rates = rates.filter((r) => r.id !== id);
	}
</script>

<BackButton href="/currencies" label={t('currencies.title')} />

<div class="mb-6 flex items-center justify-between">
	<div>
		<h1 class="text-2xl font-bold text-foreground">{t('currencies.exchangeRates.title')}</h1>
		<p class="text-sm text-muted-foreground mt-0.5">{t('currencies.exchangeRates.description')}</p>
	</div>
	{#if !isLoading}
		<Button size="sm" onclick={() => (showAddForm = !showAddForm)}>
			<Plus class="size-4 mr-1" />
			{t('common.add')}
		</Button>
	{/if}
</div>

{#if showAddForm}
	<Card class="mb-6">
		<CardHeader>
			<CardTitle class="text-base">{t('currencies.exchangeRates.addTitle')}</CardTitle>
		</CardHeader>
		<CardContent>
			<form class="grid gap-4 md:grid-cols-5 items-end" onsubmit={handleAdd}>
				<div class="flex flex-col gap-2">
					<Label>{t('currencies.exchangeRates.from')}</Label>
					<div class="relative">
						<Select bind:value={addFrom}>
							<SelectTrigger class="w-full"></SelectTrigger>
							<SelectContent>
								{#each currencies as c}
									<SelectItem value={c.id}>{c.code}</SelectItem>
								{/each}
							</SelectContent>
						</Select>
						<ChevronDown class="pointer-events-none absolute right-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
					</div>
				</div>
				<div class="flex flex-col gap-2">
					<Label>{t('currencies.exchangeRates.to')}</Label>
					<div class="relative">
						<Select bind:value={addTo}>
							<SelectTrigger class="w-full"></SelectTrigger>
							<SelectContent>
								{#each currencies as c}
									<SelectItem value={c.id}>{c.code}</SelectItem>
								{/each}
							</SelectContent>
						</Select>
						<ChevronDown class="pointer-events-none absolute right-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
					</div>
				</div>
				<div class="flex flex-col gap-2">
					<Label>{t('currencies.exchangeRates.rate')}</Label>
					<Input type="number" step="0.000001" placeholder="0" bind:value={addRate} required />
				</div>
				<div class="flex flex-col gap-2">
					<Label>{t('currencies.exchangeRates.date')}</Label>
					<Input type="date" bind:value={addDate} />
				</div>
				<div class="flex gap-2">
					<Button type="submit" size="sm" disabled={addLoading || !addFrom || !addTo || !addRate}>
						{addLoading ? t('common.saving') : t('common.save')}
					</Button>
					<Button type="button" size="sm" variant="outline" onclick={() => (showAddForm = false)}>
						{t('common.cancel')}
					</Button>
				</div>
			</form>
			{#if addError}
				<p class="text-destructive text-sm mt-2">{addError}</p>
			{/if}
		</CardContent>
	</Card>
{/if}

<Card>
	<CardHeader>
		<CardTitle class="text-base">{t('currencies.exchangeRates.title')}</CardTitle>
	</CardHeader>
	<CardContent class="p-0">
		{#if isLoading}
{#each Array(8) as _}
<TableRow>
				<TableCell><Skeleton class="h-4 w-full" /></TableCell>
				<TableCell><Skeleton class="h-4 w-full" /></TableCell>
				<TableCell><Skeleton class="h-4 w-full" /></TableCell>
				<TableCell class="hidden md:table-cell"><Skeleton class="h-4 w-full" /></TableCell>
			</TableRow>
			{/each}
		{:else if error}
			<div class="p-8 text-center text-destructive">{error}</div>
		{:else if rates.length === 0}
			<EmptyState />
		{:else}
			<Table>
				<TableHeader>
					<TableRow>
						<TableHead>{t('currencies.exchangeRates.from')}</TableHead>
						<TableHead>{t('currencies.exchangeRates.to')}</TableHead>
						<TableHead class="text-right">{t('currencies.exchangeRates.rate')}</TableHead>
						<TableHead class="hidden md:table-cell">{t('currencies.exchangeRates.date')}</TableHead>
						<TableHead class="w-12"></TableHead>
					</TableRow>
				</TableHeader>
				<TableBody>
					{#each rates as rate}
						<TableRow>
							<TableCell class="font-mono font-medium text-foreground">{rate.from_code}</TableCell>
							<TableCell class="font-mono font-medium text-foreground">{rate.to_code}</TableCell>
							<TableCell class="text-right font-medium text-foreground">{rate.rate}</TableCell>
							<TableCell class="hidden md:table-cell text-muted-foreground">{formatDate(rate.date)}</TableCell>
							<TableCell>
								<button
									type="button"
									aria-label={t('common.delete')}
									class="text-muted-foreground hover:text-destructive transition-colors"
									onclick={() => (deleteTarget = rate.id)}
								>
									<Trash2 class="size-4" />
								</button>
							</TableCell>
						</TableRow>
					{/each}
				</TableBody>
			</Table>
		{/if}
	</CardContent>
</Card>

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

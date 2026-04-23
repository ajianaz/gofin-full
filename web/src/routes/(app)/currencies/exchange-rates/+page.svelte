<script lang="ts">
	import { BackButton } from '$lib/components/shared/index.js';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { currencyService } from '$lib/services/index.js';
	import { formatDate } from '$lib/utils/format.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import { onMount } from 'svelte';
	import type { ExchangeRate } from '$lib/types/domain.js';

	const t = localeStore.t;

	let rates: ExchangeRate[] = $state([]);
	let isLoading: boolean = $state(true);
	let error: string | null = $state(null);

	onMount(async () => {
		try {
			rates = await currencyService.exchangeRates();
		} catch (e: any) {
			error = e?.message ?? t('common.error');
		} finally {
			isLoading = false;
		}
	});
</script>

<BackButton href="/currencies" label={t('currencies.title')} />

<div class="mb-6">
	<h1 class="text-2xl font-bold text-foreground">{t('currencies.exchangeRates.title')}</h1>
	<p class="text-sm text-muted-foreground mt-0.5">{t('currencies.exchangeRates.description')}</p>
</div>

<Card>
	<CardHeader>
		<CardTitle class="text-base">{t('currencies.exchangeRates.title')}</CardTitle>
	</CardHeader>
	<CardContent class="p-0">
		{#if isLoading}
			<div class="p-8 text-center text-muted-foreground">{t('common.loading')}</div>
		{:else if error}
			<div class="p-8 text-center text-destructive">{error}</div>
		{:else if rates.length === 0}
			<div class="p-8 text-center text-muted-foreground">{t('common.noData')}</div>
		{:else}
			<div class="overflow-x-auto">
				<table class="w-full text-sm">
					<thead>
						<tr class="border-b bg-muted/50">
							<th class="text-left p-3 font-medium text-muted-foreground">{t('currencies.exchangeRates.from')}</th>
							<th class="text-left p-3 font-medium text-muted-foreground">{t('currencies.exchangeRates.to')}</th>
							<th class="text-right p-3 font-medium text-muted-foreground">{t('currencies.exchangeRates.rate')}</th>
							<th class="text-left p-3 font-medium text-muted-foreground">{t('currencies.exchangeRates.date')}</th>
						</tr>
					</thead>
					<tbody>
						{#each rates as rate}
							<tr class="border-b hover:bg-muted/30">
								<td class="p-3 font-mono font-medium text-foreground">{rate.from_code}</td>
								<td class="p-3 font-mono font-medium text-foreground">{rate.to_code}</td>
								<td class="p-3 text-right font-medium text-foreground">{rate.rate}</td>
								<td class="p-3 text-muted-foreground">{formatDate(rate.date)}</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{/if}
	</CardContent>
</Card>

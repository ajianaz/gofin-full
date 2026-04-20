<script lang="ts">
	import { BackButton } from '$lib/components/shared/index.js';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { mockExchangeRates } from '$lib/data/mock-exchange-rates.js';
	import { formatDate } from '$lib/utils/format.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;
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
					{#each mockExchangeRates as rate}
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
	</CardContent>
</Card>

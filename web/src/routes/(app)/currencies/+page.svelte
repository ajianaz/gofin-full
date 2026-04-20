<script lang="ts">
	import { PageHeader, StatusBadge } from '$lib/components/shared/index.js';
	import { Card, CardContent } from '$lib/components/ui/card/index.js';
	import { mockCurrencies } from '$lib/data/mock-currencies.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;
</script>

<PageHeader title={t('currencies.title')} description={t('currencies.description')} />

<Card>
	<CardContent class="p-0">
		<div class="overflow-x-auto">
			<table class="w-full text-sm">
				<thead>
					<tr class="border-b bg-muted/50">
						<th class="text-left p-3 font-medium text-muted-foreground">{t('currencies.code')}</th>
						<th class="text-left p-3 font-medium text-muted-foreground">{t('currencies.name')}</th>
						<th class="text-left p-3 font-medium text-muted-foreground">{t('currencies.symbol')}</th>
						<th class="text-left p-3 font-medium text-muted-foreground">{t('currencies.decimalPlaces')}</th>
						<th class="text-left p-3 font-medium text-muted-foreground">{t('common.status')}</th>
					</tr>
				</thead>
				<tbody>
					{#each mockCurrencies as currency}
						<tr class="border-b hover:bg-muted/30">
							<td class="p-3 font-mono font-medium text-foreground">{currency.code}</td>
							<td class="p-3 text-foreground">{currency.name}</td>
							<td class="p-3 text-muted-foreground">{currency.symbol}</td>
							<td class="p-3 text-muted-foreground">{currency.decimal_places}</td>
							<td class="p-3"><StatusBadge status={currency.enabled ? 'active' : 'inactive'} /></td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</CardContent>
</Card>

<div class="mt-6">
	<a href="/currencies/exchange-rates" class="text-sm text-primary font-medium hover:underline">{t('currencies.viewRates')}</a>
</div>

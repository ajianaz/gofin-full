<script lang="ts">
	import { BackButton } from '$lib/components/shared/index.js';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '$lib/components/ui/table/index.js';
	import { currencyService } from '$lib/services/index.js';
	import { formatDate } from '$lib/utils/format.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import { onMount } from 'svelte';
	import type { ExchangeRate } from '$lib/types/domain.js';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import EmptyState from '$lib/components/shared/EmptyState.svelte';

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
						</TableRow>
					</TableHeader>
					<TableBody>
						{#each rates as rate}
							<TableRow>
								<TableCell class="font-mono font-medium text-foreground">{rate.from_code}</TableCell>
								<TableCell class="font-mono font-medium text-foreground">{rate.to_code}</TableCell>
								<TableCell class="text-right font-medium text-foreground">{rate.rate}</TableCell>
								<TableCell class="hidden md:table-cell text-muted-foreground">{formatDate(rate.date)}</TableCell>
							</TableRow>
						{/each}
					</TableBody>
				</Table>
		{/if}
	</CardContent>
</Card>

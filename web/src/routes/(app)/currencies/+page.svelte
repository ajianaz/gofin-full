<script lang="ts">
	import { onMount } from 'svelte';
	import { PageHeader, StatusBadge } from '$lib/components/shared/index.js';
	import { Card, CardContent } from '$lib/components/ui/card/index.js';
	import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '$lib/components/ui/table/index.js';
	import { currencyService } from '$lib/services/index.js';
	import type { Currency } from '$lib/types/domain.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	const t = localeStore.t;

	let items = $state<Currency[]>([]);
	let isLoading = $state(true);
	let errorMsg = $state('');

	onMount(async () => {
		try {
			items = await currencyService.list();
		} catch (e) {
			errorMsg = t('common.error');
			console.error(e);
		} finally {
			isLoading = false;
		}
	});
</script>

<PageHeader title={t('currencies.title')} description={t('currencies.description')} />

<Card>
	<CardContent class="p-0">
		<Table>
				<TableHeader>
					<TableRow>
						<TableHead>{t('currencies.code')}</TableHead>
						<TableHead>{t('currencies.name')}</TableHead>
						<TableHead class="hidden md:table-cell">{t('currencies.symbol')}</TableHead>
						<TableHead>{t('currencies.decimalPlaces')}</TableHead>
						<TableHead class="hidden md:table-cell">{t('common.status')}</TableHead>
					</TableRow>
				</TableHeader>
				<TableBody>
					{#if isLoading}
				{#each Array(5) as _}
<TableRow>
					<TableCell><Skeleton class="h-4 w-full" /></TableCell>
					<TableCell><Skeleton class="h-4 w-full" /></TableCell>
					<TableCell class="hidden md:table-cell"><Skeleton class="h-4 w-full" /></TableCell>
					<TableCell class="hidden md:table-cell"><Skeleton class="h-4 w-full" /></TableCell>
					<TableCell><Skeleton class="h-4 w-full" /></TableCell>
				</TableRow>
				{/each}
	{:else if errorMsg}
						<TableRow>
							<TableCell colspan={5} class="py-8 text-center text-sm text-destructive">{errorMsg}</TableCell>
						</TableRow>
					{:else}
						{#each items as currency}
							<TableRow>
								<TableCell class="font-mono font-medium text-foreground">{currency.code}</TableCell>
								<TableCell class="text-foreground">{currency.name}</TableCell>
								<TableCell class="hidden md:table-cell text-muted-foreground">{currency.symbol}</TableCell>
								<TableCell class="text-muted-foreground">{currency.decimal_places}</TableCell>
								<TableCell><StatusBadge status={currency.enabled ? 'active' : 'inactive'} /></TableCell>
							</TableRow>
						{:else}
							<TableRow>
								<TableCell colspan={5}><EmptyState /></TableCell>
							</TableRow>
						{/each}
					{/if}
				</TableBody>
			</Table>
	</CardContent>
</Card>

<div class="mt-6">
	<a href="/currencies/exchange-rates" class="text-sm text-primary font-medium hover:underline">{t('currencies.viewRates')}</a>
</div>

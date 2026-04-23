<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { ArrowLeft } from '@lucide/svelte';
	import { walletService } from '$lib/services/index.js';
	import type { Account } from '$lib/types/domain.js';
	import { formatCurrency } from '$lib/utils/format.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	let loading = $state(true);
	let error = $state('');

	let wallets: Account[] = $state([]);

	let assets = $derived(wallets.filter((w) => w.type === 'asset' || w.type === 'cash'));
	let liabilities = $derived(wallets.filter((w) => w.type === 'liability'));
	let totalAssets = $derived(assets.reduce((s, w) => s + parseFloat(w.balance), 0));
	let totalLiabilities = $derived(liabilities.reduce((s, w) => s + Math.abs(parseFloat(w.balance)), 0));
	let netWorth = $derived(totalAssets - totalLiabilities);

	onMount(async () => {
		try {
			wallets = await walletService.list();
		} catch (e: any) {
			error = e?.detail || e?.message || 'Failed to load wallets';
		} finally {
			loading = false;
		}
	});
</script>

<div class="flex flex-col gap-4">
	<div class="flex items-center gap-3">
		<Button variant="ghost" size="sm" onclick={() => goto('/reports')}>
			<ArrowLeft class="size-4" />
			{t('common.back')}
		</Button>
		<h2 class="text-base font-semibold text-foreground">{t('reports.netWorth.title')}</h2>
	</div>

	{#if loading}
		<div class="flex items-center justify-center py-12">
			<p class="text-sm text-muted-foreground">{t('common.loading')}</p>
		</div>
	{:else if error}
		<div class="flex items-center justify-center py-12">
			<p class="text-sm text-destructive">{error}</p>
		</div>
	{:else}
		<Card>
			<CardHeader class="pb-1"><CardTitle class="text-sm font-semibold">{t('reports.netWorth.title')}</CardTitle></CardHeader>
			<CardContent>
				<p class="text-xl font-bold {netWorth >= 0 ? 'text-green-600' : 'text-destructive'}">{formatCurrency(netWorth.toString())}</p>
			</CardContent>
		</Card>

		<div class="grid gap-4 md:grid-cols-2">
			<Card>
				<CardHeader>
					<CardTitle class="text-base">{t('reports.netWorth.assets')}</CardTitle>
				</CardHeader>
				<CardContent class="p-0">
					{#if assets.length === 0}
						<p class="px-4 py-6 text-sm text-muted-foreground text-center">{t('common.noData')}</p>
					{:else}
						{#each assets as w}
							<div class="flex items-center justify-between px-4 py-3 border-b last:border-b-0">
								<span class="text-sm font-medium text-foreground">{w.name}</span>
								<span class="text-sm font-medium text-green-600">{formatCurrency(w.balance)}</span>
							</div>
						{/each}
					{/if}
					<div class="flex items-center justify-between bg-muted/50 px-4 py-3 font-semibold">
						<span class="text-sm text-foreground">{t('reports.netWorth.totalAssets')}</span>
						<span class="text-sm text-foreground">{formatCurrency(totalAssets.toString())}</span>
					</div>
				</CardContent>
			</Card>

			<Card>
				<CardHeader>
					<CardTitle class="text-base">{t('reports.netWorth.liabilities')}</CardTitle>
				</CardHeader>
				<CardContent class="p-0">
					{#if liabilities.length === 0}
						<p class="px-4 py-6 text-sm text-muted-foreground text-center">{t('common.noData')}</p>
					{:else}
						{#each liabilities as w}
							<div class="flex items-center justify-between px-4 py-3 border-b last:border-b-0">
								<span class="text-sm font-medium text-foreground">{w.name}</span>
								<span class="text-sm font-medium text-destructive">{formatCurrency(w.balance)}</span>
							</div>
						{/each}
					{/if}
					<div class="flex items-center justify-between bg-muted/50 px-4 py-3 font-semibold">
						<span class="text-sm text-foreground">{t('reports.netWorth.totalLiabilities')}</span>
						<span class="text-sm text-foreground">{formatCurrency(totalLiabilities.toString())}</span>
					</div>
				</CardContent>
			</Card>
		</div>
	{/if}
</div>

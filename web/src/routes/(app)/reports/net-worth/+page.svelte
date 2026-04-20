<script lang="ts">
	import { goto } from '$app/navigation';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { ArrowLeft } from '@lucide/svelte';
	import { mockWallets } from '$lib/data/mock-wallets.js';
	import { formatCurrency } from '$lib/utils/format.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	const assets = mockWallets.filter((w) => w.type === 'asset' || w.type === 'cash');
	const liabilities = mockWallets.filter((w) => w.type === 'liability');
	const totalAssets = assets.reduce((s, w) => s + parseFloat(w.balance), 0);
	const totalLiabilities = liabilities.reduce((s, w) => s + Math.abs(parseFloat(w.balance)), 0);
	const netWorth = totalAssets - totalLiabilities;
</script>

<div class="flex flex-col gap-4">
	<div class="flex items-center gap-3">
		<Button variant="ghost" size="sm" onclick={() => goto('/reports')}>
			<ArrowLeft class="size-4" />
			{t('common.back')}
		</Button>
		<h2 class="text-base font-semibold text-foreground">{t('reports.netWorth.title')}</h2>
	</div>

	<Card>
		<CardHeader class="pb-1"><CardTitle class="text-sm font-semibold">{t('reports.netWorth.title')}</CardTitle></CardHeader>
		<CardContent>
			<p class="text-xl font-bold {netWorth >= 0 ? 'text-green-600' : 'text-destructive'}">{formatCurrency(netWorth.toString())}</p>
			<p class="text-xs text-muted-foreground">+8% {t('reports.netWorth.fromLastMonth')}</p>
		</CardContent>
	</Card>

	<div class="grid gap-4 md:grid-cols-2">
		<Card>
			<CardHeader>
				<CardTitle class="text-base">{t('reports.netWorth.assets')}</CardTitle>
			</CardHeader>
			<CardContent class="p-0">
				{#each assets as w}
					<div class="flex items-center justify-between px-4 py-3 border-b last:border-b-0">
						<span class="text-sm font-medium text-foreground">{w.name}</span>
						<span class="text-sm font-medium text-green-600">{formatCurrency(w.balance)}</span>
					</div>
				{/each}
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
				{#each liabilities as w}
					<div class="flex items-center justify-between px-4 py-3 border-b last:border-b-0">
						<span class="text-sm font-medium text-foreground">{w.name}</span>
						<span class="text-sm font-medium text-destructive">{formatCurrency(w.balance)}</span>
					</div>
				{/each}
				<div class="flex items-center justify-between bg-muted/50 px-4 py-3 font-semibold">
					<span class="text-sm text-foreground">{t('reports.netWorth.totalLiabilities')}</span>
					<span class="text-sm text-foreground">{formatCurrency(totalLiabilities.toString())}</span>
				</div>
			</CardContent>
		</Card>
	</div>
</div>

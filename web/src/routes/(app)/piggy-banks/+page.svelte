<script lang="ts">
	import { goto } from '$app/navigation';
	import { Card, CardContent } from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import { Plus, ChevronDown, PiggyBank } from '@lucide/svelte';
	import { mockPiggyBanks } from '$lib/data/mock-piggy-banks.js';
	import { formatCurrency, formatPercentage } from '$lib/utils/format.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	let accountFilter = $state('all');

	let filtered = $derived(
		accountFilter === 'all'
			? mockPiggyBanks
			: mockPiggyBanks.filter((pb) => pb.account_id === accountFilter)
	);

	const accounts = [...new Set(mockPiggyBanks.map((pb) => pb.account_id))];
</script>

<div class="flex flex-col gap-4">
	<div class="flex flex-wrap items-center justify-between gap-3">
		<div class="flex items-center gap-3">
			<h2 class="text-lg font-semibold text-foreground">{t('piggyBanks.list.title')}</h2>
			<span class="inline-flex items-center rounded-2xl bg-primary px-2.5 py-0.5 text-xs font-medium text-primary-foreground">
				{filtered.length}
			</span>
		</div>
		<div class="flex items-center gap-3">
			<div class="relative">
				<select
					bind:value={accountFilter}
					class="cn-input h-9 w-44 appearance-none bg-background pr-8 text-sm"
				>
					<option value="all">{t('piggyBanks.list.allWallets')}</option>
					{#each accounts as id}
						{@const w = mockPiggyBanks.find((pb) => pb.account_id === id)}
						<option value={id}>{w?.account_name ?? id}</option>
					{/each}
				</select>
				<ChevronDown class="pointer-events-none absolute right-2.5 top-1/2 size-3.5 -translate-y-1/2 text-muted-foreground" />
			</div>
			<Button size="sm" onclick={() => goto('/piggy-banks/create')}>
				<Plus class="size-4" />
				{t('piggyBanks.list.add')}
			</Button>
		</div>
	</div>

	<Card>
		<CardContent class="p-0">
			<p class="px-5 pt-4 text-sm font-semibold text-muted-foreground">{t('piggyBanks.list.active')}</p>
			{#each filtered as pb}
				{@const pct = parseFloat(pb.target_amount) > 0
					? (parseFloat(pb.current_amount) / parseFloat(pb.target_amount)) * 100
					: 0}
				<div class="flex items-center gap-4 px-5 py-4 border-b last:border-b-0">
					<div class="flex size-10 shrink-0 items-center justify-center rounded-lg bg-muted">
						<PiggyBank class="size-5 text-foreground" />
					</div>
					<div class="flex flex-col gap-1 min-w-0">
						<p class="text-sm font-semibold text-foreground truncate">{pb.name}</p>
						<p class="text-xs text-muted-foreground">{t('piggyBanks.list.walletPrefix')} {pb.account_name}</p>
						<Progress value={Math.min(pct, 100)} class="mt-1 h-1.5 w-[120px]" />
					</div>
					<div class="ml-auto shrink-0 text-right">
						<p class="text-sm font-semibold text-foreground">{formatCurrency(pb.current_amount)}</p>
						<p class="text-xs text-muted-foreground">{t('piggyBanks.list.of', { pct: Math.round(pct), target: formatCurrency(pb.target_amount) })}</p>
					</div>
				</div>
			{/each}
		</CardContent>
	</Card>
</div>

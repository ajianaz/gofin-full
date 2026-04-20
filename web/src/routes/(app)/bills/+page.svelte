<script lang="ts">
	import { goto } from '$app/navigation';
	import { Card, CardContent } from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Plus, ChevronDown } from '@lucide/svelte';
	import { mockBills } from '$lib/data/mock-bills.js';
	import { formatCurrency, formatDate } from '$lib/utils/format.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	let accountFilter = $state('all');

	let filtered = $derived(
		accountFilter === 'all' ? mockBills : mockBills.filter((b) => b.active)
	);
</script>

<div class="flex flex-col gap-4">
	<div class="flex flex-wrap items-center justify-between gap-3">
		<div class="flex items-center gap-2">
			<h2 class="text-base font-semibold text-foreground">{t('bills.list.title')}</h2>
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
					<option value="all">{t('bills.list.allWallets')}</option>
					<option value="active">{t('bills.list.activeOnly')}</option>
				</select>
				<ChevronDown class="pointer-events-none absolute right-2.5 top-1/2 size-3.5 -translate-y-1/2 text-muted-foreground" />
			</div>
			<Button size="sm" onclick={() => goto('/bills/create')}>
				<Plus class="size-4" />
				{t('bills.list.add')}
			</Button>
		</div>
	</div>

	<Card>
		<CardContent class="p-0">
			<p class="px-5 pt-4 text-[13px] font-normal text-muted-foreground">{t('bills.list.activeList')}</p>
			{#each filtered as bill}
				<div class="flex items-center justify-between px-5 py-3 border-b last:border-b-0">
					<div class="flex flex-col gap-1">
						<p class="text-sm font-semibold text-foreground">{bill.name}</p>
						<p class="text-[13px] text-muted-foreground">
							{#if bill.amount_min === bill.amount_max}
								{formatCurrency(bill.amount_min)}{t('bills.list.perMonth')}
							{:else}
								{formatCurrency(bill.amount_min)} — {formatCurrency(bill.amount_max)}
							{/if}
						</p>
					</div>
					<div class="flex shrink-0 items-center gap-3">
						<span class="text-[13px] text-muted-foreground">{formatDate(bill.next_date)}</span>
						{#if bill.active}
							<Badge variant="secondary" class="text-xs">{t('bills.list.active')}</Badge>
						{:else}
							<Badge variant="outline" class="text-xs">{t('bills.list.inactive')}</Badge>
						{/if}
					</div>
				</div>
			{/each}
		</CardContent>
	</Card>
</div>

<script lang="ts">
	import { goto } from '$app/navigation';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { ArrowLeft, ChevronDown } from '@lucide/svelte';
	import { mockTransactions } from '$lib/data/mock-transactions.js';
	import { formatCurrency } from '$lib/utils/format.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	let period = $state('monthly');

	let periodData = $derived(() => {
		const monthly: Record<string, { income: number; expense: number }> = {};
		for (const tx of mockTransactions) {
			const month = tx.date.substring(0, 7);
			if (!monthly[month]) monthly[month] = { income: 0, expense: 0 };
			if (tx.type === 'deposit') monthly[month].income += parseFloat(tx.amount);
			else if (tx.type === 'withdrawal') monthly[month].expense += Math.abs(parseFloat(tx.amount));
		}
		return Object.entries(monthly)
			.sort(([a], [b]) => b.localeCompare(a))
			.map(([month, data]) => {
				const diff = data.income - data.expense;
				return { period: month, ...data, diff };
			});
	});
</script>

<div class="flex flex-col gap-4">
	<div class="flex flex-wrap items-center justify-between gap-3">
		<div class="flex items-center gap-3">
			<Button variant="ghost" size="sm" onclick={() => goto('/reports')}>
				<ArrowLeft class="size-4" />
				{t('common.back')}
			</Button>
			<h2 class="text-base font-semibold text-foreground">{t('reports.spendingByPeriod.title')}</h2>
		</div>
		<div class="relative">
			<select bind:value={period} class="cn-input h-9 w-36 appearance-none bg-background pr-8 text-sm">
				<option value="daily">{t('reports.period.daily')}</option>
				<option value="weekly">{t('reports.period.weekly')}</option>
				<option value="monthly">{t('reports.period.monthly')}</option>
				<option value="yearly">{t('reports.period.yearly')}</option>
			</select>
			<ChevronDown class="pointer-events-none absolute right-2.5 top-1/2 size-3.5 -translate-y-1/2 text-muted-foreground" />
		</div>
	</div>

	<Card>
		<CardHeader>
			<CardTitle class="text-base">{t('reports.spendingByPeriod.history')}</CardTitle>
		</CardHeader>
		<CardContent class="p-0">
			<div class="overflow-x-auto">
				<table class="w-full text-sm">
					<thead>
						<tr class="border-b bg-muted/50">
							<th class="text-left px-4 py-3 font-medium text-muted-foreground">{t('reports.spendingByPeriod.period')}</th>
							<th class="text-right px-4 py-3 font-medium text-muted-foreground">{t('reports.spendingByPeriod.income')}</th>
							<th class="text-right px-4 py-3 font-medium text-muted-foreground">{t('reports.spendingByPeriod.expense')}</th>
							<th class="text-right px-4 py-3 font-medium text-muted-foreground">{t('reports.spendingByPeriod.diff')}</th>
						</tr>
					</thead>
					<tbody>
						{#each periodData() as row}
							<tr class="border-b last:border-b-0 hover:bg-muted/30">
								<td class="px-4 py-3 font-medium text-foreground">{row.period}</td>
								<td class="px-4 py-3 text-right text-green-600">{formatCurrency(row.income.toString())}</td>
								<td class="px-4 py-3 text-right text-destructive">{formatCurrency(row.expense.toString())}</td>
								<td class="px-4 py-3 text-right font-medium {row.diff >= 0 ? 'text-green-600' : 'text-destructive'}">
									{row.diff >= 0 ? '+' : '-'}{formatCurrency(Math.abs(row.diff).toString())}
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</CardContent>
	</Card>
</div>

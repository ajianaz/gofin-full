<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { ChevronDown } from '@lucide/svelte';
	import { budgetService } from '$lib/services/index.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	let name = $state('');
	let active = $state(true);
	let autoBudgetType = $state('none');
	let period = $state('monthly');
	let amount = $state('');
	let isLoading = $state(false);
	let errorMsg = $state('');
</script>

<div class="flex flex-col gap-4">
	<Card>
		<CardHeader>
			<CardTitle>{t('budgets.create.title')}</CardTitle>
		</CardHeader>
		<CardContent>
			<form class="flex flex-col gap-4" onsubmit={async (e) => {
					e.preventDefault();
					isLoading = true;
					errorMsg = '';
					try {
						await budgetService.create({ name });
						goto('/budgets');
					} catch (err: any) {
						errorMsg = err.detail || err.message || t('common.errorSave');
					} finally {
						isLoading = false;
					}
				}}>
				<div class="grid gap-2">
					<Label for="name">{t('budgets.create.name')}</Label>
					<Input id="name" placeholder={t('budgets.create.namePlaceholder')} bind:value={name} required />
				</div>

				<div class="flex items-center justify-between">
					<Label for="active">{t('budgets.create.active')}</Label>
					<Checkbox id="active" bind:checked={active} />
				</div>

				<div class="grid gap-2">
					<Label for="auto-type">{t('budgets.create.autoBudgetType')}</Label>
					<div class="relative">
						<select
							id="auto-type"
							bind:value={autoBudgetType}
							class="cn-input w-full appearance-none bg-background pr-8"
						>
							<option value="none">{t('budgets.create.autoBudgetNone')}</option>
							<option value="reset">{t('budgets.create.autoBudgetReset')}</option>
							<option value="rollover">{t('budgets.create.autoBudgetRollover')}</option>
							<option value="fixed">{t('budgets.create.autoBudgetFixed')}</option>
							<option value="adjust">{t('budgets.create.autoBudgetAdjust')}</option>
						</select>
						<ChevronDown class="pointer-events-none absolute right-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
					</div>
				</div>

				<div class="grid gap-2">
					<Label for="amount">{t('budgets.create.amount')}</Label>
					<Input id="amount" type="number" placeholder="3000000" bind:value={amount} required />
				</div>

				<div class="grid gap-2">
					<Label for="period">{t('budgets.create.period')}</Label>
					<div class="relative">
						<select
							id="period"
							bind:value={period}
							class="cn-input w-full appearance-none bg-background pr-8"
						>
							<option value="daily">{t('budgets.create.periodDaily')}</option>
							<option value="weekly">{t('budgets.create.periodWeekly')}</option>
							<option value="monthly">{t('budgets.create.periodMonthly')}</option>
							<option value="quarterly">{t('budgets.create.periodQuarterly')}</option>
							<option value="yearly">{t('budgets.create.periodYearly')}</option>
						</select>
						<ChevronDown class="pointer-events-none absolute right-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
					</div>
				</div>

				<div class="flex gap-2 pt-2">
					{#if errorMsg}
						<p class="text-destructive text-sm">{errorMsg}</p>
					{/if}
					<Button type="submit" class="flex-1" disabled={isLoading}>{isLoading ? t('common.saving') : t('common.save')}</Button>
					<Button type="button" variant="outline" class="flex-1" onclick={() => goto('/budgets')}>{t('common.cancel')}</Button>
				</div>
			</form>
		</CardContent>
	</Card>
</div>

<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { ChevronDown } from '@lucide/svelte';
	import { budgetService } from '$lib/services/index.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import { Select, SelectTrigger, SelectContent, SelectItem } from '$lib/components/ui/select/index.js';
	import FormCard from '$lib/components/shared/FormCard.svelte';
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
	<FormCard title="{t('budgets.create.title')}">
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
						<Select bind:value={autoBudgetType} id="auto-type">
		<SelectTrigger class="w-full">
		</SelectTrigger>
		<SelectContent>
		<SelectItem value="none">{t('budgets.create.autoBudgetNone')}</SelectItem>
		<SelectItem value="reset">{t('budgets.create.autoBudgetReset')}</SelectItem>
		<SelectItem value="rollover">{t('budgets.create.autoBudgetRollover')}</SelectItem>
		<SelectItem value="fixed">{t('budgets.create.autoBudgetFixed')}</SelectItem>
		<SelectItem value="adjust">{t('budgets.create.autoBudgetAdjust')}</SelectItem>
		</SelectContent>
</Select>
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
						<Select bind:value={period} id="period">
		<SelectTrigger class="w-full">
		</SelectTrigger>
		<SelectContent>
		<SelectItem value="daily">{t('budgets.create.periodDaily')}</SelectItem>
		<SelectItem value="weekly">{t('budgets.create.periodWeekly')}</SelectItem>
		<SelectItem value="monthly">{t('budgets.create.periodMonthly')}</SelectItem>
		<SelectItem value="quarterly">{t('budgets.create.periodQuarterly')}</SelectItem>
		<SelectItem value="yearly">{t('budgets.create.periodYearly')}</SelectItem>
		</SelectContent>
</Select>
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
		</FormCard>
</div>

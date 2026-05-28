<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { ChevronDown } from '@lucide/svelte';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import { Select, SelectTrigger, SelectContent, SelectItem, SelectValue } from '$lib/components/ui/select/index.js';
	const t = localeStore.t;

	let groupId = $derived($page.params.groupId);

	let title = $state('');
	let triggerType = $state('description_contains');
	let triggerOperator = $state('contains');
	let triggerValue = $state('');
	let actionType = $state('set_category');
	let actionValue = $state('');
	let strict = $state(false);
	let stopProcessing = $state(false);
	let active = $state(true);
</script>

<div class="flex flex-col gap-4">
	<Card>
		<CardHeader>
			<div class="flex items-center justify-between">
				<CardTitle class="text-base">{t('rules.createRule.title')}</CardTitle>
				<div class="flex gap-2">
					<Button size="sm" onclick={() => goto('/rules/{groupId}')}>{t('common.save')}</Button>
					<Button size="sm" variant="outline" onclick={() => goto('/rules/{groupId}')}>{t('common.cancel')}</Button>
				</div>
			</div>
		</CardHeader>
		<CardContent>
			<form class="flex flex-col gap-6" onsubmit={(e) => { e.preventDefault(); goto('/rules/{groupId}'); }}>
				<div class="grid gap-6 md:grid-cols-2">
					<div class="flex flex-col gap-2">
						<Label for="title">{t('rules.createRule.ruleName')}</Label>
						<Input id="title" placeholder={t('rules.createRule.ruleNamePlaceholder')} bind:value={title} required />
					</div>
				</div>

				<div class="flex flex-col gap-4">
					<div class="border-b"></div>
					<h3 class="text-sm font-semibold text-foreground">{t('rules.createRule.trigger')}</h3>
					<div class="grid gap-4 md:grid-cols-3">
						<div class="flex flex-col gap-2">
							<Label for="trigger-type">{t('rules.createRule.triggerType')}</Label>
							<div class="relative">
								<Select bind:value={triggerType} id="trigger-type">
<SelectTrigger class="w-full">
		<SelectValue />
	</SelectTrigger>
		<SelectContent>
		<SelectItem value="description_contains" label={t('rules.createRule.triggerDescContains')}>{t('rules.createRule.triggerDescContains')}</SelectItem>
		<SelectItem value="amount_less_than" label={t('rules.createRule.triggerAmountLess')}>{t('rules.createRule.triggerAmountLess')}</SelectItem>
		<SelectItem value="amount_greater_than" label={t('rules.createRule.triggerAmountGreater')}>{t('rules.createRule.triggerAmountGreater')}</SelectItem>
		<SelectItem value="deposit" label={t('rules.createRule.triggerDeposit')}>{t('rules.createRule.triggerDeposit')}</SelectItem>
		</SelectContent>
</Select>
								
							</div>
						</div>
						<div class="flex flex-col gap-2">
							<Label for="trigger-op">{t('rules.createRule.operator')}</Label>
							<div class="relative">
								<Select bind:value={triggerOperator} id="trigger-op">
<SelectTrigger class="w-full">
		<SelectValue />
	</SelectTrigger>
		<SelectContent>
		<SelectItem value="contains" label={t('rules.createRule.opContains')}>{t('rules.createRule.opContains')}</SelectItem>
		<SelectItem value="equals" label={t('rules.createRule.opEquals')}>{t('rules.createRule.opEquals')}</SelectItem>
		<SelectItem value="starts_with" label={t('rules.createRule.opStartsWith')}>{t('rules.createRule.opStartsWith')}</SelectItem>
		<SelectItem value="ends_with" label={t('rules.createRule.opEndsWith')}>{t('rules.createRule.opEndsWith')}</SelectItem>
		</SelectContent>
</Select>
								
							</div>
						</div>
						<div class="flex flex-col gap-2">
							<Label for="trigger-val">{t('rules.createRule.value')}</Label>
							<Input id="trigger-val" placeholder={t('rules.createRule.valuePlaceholder')} bind:value={triggerValue} />
						</div>
					</div>
				</div>

				<div class="flex flex-col gap-4">
					<div class="border-b"></div>
					<h3 class="text-sm font-semibold text-foreground">{t('rules.createRule.action')}</h3>
					<div class="grid gap-4 md:grid-cols-2">
						<div class="flex flex-col gap-2">
							<Label for="action-type">{t('rules.createRule.actionType')}</Label>
							<div class="relative">
								<Select bind:value={actionType} id="action-type">
<SelectTrigger class="w-full">
		<SelectValue />
	</SelectTrigger>
		<SelectContent>
		<SelectItem value="set_category" label={t('rules.createRule.actionSetCategory')}>{t('rules.createRule.actionSetCategory')}</SelectItem>
		<SelectItem value="add_tag" label={t('rules.createRule.actionAddTag')}>{t('rules.createRule.actionAddTag')}</SelectItem>
		<SelectItem value="move_to_account" label={t('rules.createRule.actionMoveToAccount')}>{t('rules.createRule.actionMoveToAccount')}</SelectItem>
		<SelectItem value="set_budget" label={t('rules.createRule.actionSetBudget')}>{t('rules.createRule.actionSetBudget')}</SelectItem>
		</SelectContent>
</Select>
								
							</div>
						</div>
						<div class="flex flex-col gap-2">
							<Label for="action-val">{t('rules.createRule.value')}</Label>
							<Input id="action-val" placeholder={t('rules.createRule.actionValuePlaceholder')} bind:value={actionValue} />
						</div>
					</div>
				</div>

				<div class="flex items-center gap-6">
					<div class="flex items-center gap-2">
						<Checkbox id="active" bind:checked={active} />
						<Label for="active">{t('rules.createRule.active')}</Label>
					</div>
					<div class="flex items-center gap-2">
						<Checkbox id="strict" bind:checked={strict} />
						<Label for="strict">{t('rules.createRule.strictMatching')}</Label>
					</div>
					<div class="flex items-center gap-2">
						<Checkbox id="stop" bind:checked={stopProcessing} />
						<Label for="stop">{t('rules.createRule.stopProcessing')}</Label>
					</div>
				</div>

				<div class="flex gap-2 pt-2">
					<Button type="submit" class="flex-1">{t('common.save')}</Button>
					<Button type="button" variant="outline" class="flex-1" onclick={() => goto('/rules/{groupId}')}>{t('common.cancel')}</Button>
				</div>
			</form>
		</CardContent>
	</Card>
</div>

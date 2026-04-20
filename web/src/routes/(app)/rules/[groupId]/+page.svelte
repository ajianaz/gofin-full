<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { Card, CardContent } from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Plus, ArrowLeft, ChevronRight } from '@lucide/svelte';
	import { mockRuleGroups, mockRules } from '$lib/data/mock-rules.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	let groupId = $derived($page.params.groupId);
	let group = $derived(mockRuleGroups.find((g) => g.id === groupId));
	let rules = $derived(mockRules.filter((r) => r.rule_group_id === groupId));

	const triggerLabels: Record<string, string> = {
		description_contains: t('rules.group.triggerDescContains'),
		deposit: t('rules.group.triggerDeposit'),
		amount_less_than: t('rules.group.triggerAmountLess'),
		amount_greater_than: t('rules.group.triggerAmountGreater')
	};

	const actionLabels: Record<string, string> = {
		set_category: t('rules.group.actionSetCategory'),
		add_tag: t('rules.group.actionAddTag'),
		move_to_account: t('rules.group.actionMoveToAccount')
	};
</script>

<div class="flex flex-col gap-4">
	<div class="flex flex-wrap items-center justify-between gap-3">
		<div class="flex items-center gap-3">
			<Button variant="ghost" size="sm" onclick={() => goto('/rules')}>
				<ArrowLeft class="size-4" />
				{t('rules.group.back')}
			</Button>
			{#if group}
				<h2 class="text-base font-semibold text-foreground">{group.title}</h2>
				<span class="inline-flex items-center rounded-full bg-muted px-3 py-1 text-xs font-medium text-muted-foreground">{t('rules.group.ruleCount', { count: group.rule_count })}</span>
			{/if}
		</div>
		<Button size="sm" onclick={() => goto('/rules/{groupId}/create')}>
			<Plus class="size-4" />
			{t('rules.group.add')}
		</Button>
	</div>

	<Card>
		<CardContent class="p-0">
			{#each rules as rule}
				<button
					type="button"
					class="flex w-full items-center justify-between px-5 py-4 border-b last:border-b-0 hover:bg-muted/50 transition-colors text-left"
				>
					<div class="flex flex-col gap-1">
						<p class="text-sm font-semibold text-foreground">{rule.title}</p>
						<p class="text-[13px] text-muted-foreground">
							{t('rules.group.ifPrefix')} {triggerLabels[rule.trigger_type] ?? rule.trigger_type}
							{#if rule.trigger_value} '{rule.trigger_value}'{/if}
							→ {actionLabels[rule.action_type] ?? rule.action_type}
							{#if rule.action_value} '{rule.action_value}'{/if}
						</p>
					</div>
					<div class="flex shrink-0 items-center gap-3">
						{#if rule.active}
							<span class="inline-flex items-center rounded-full bg-primary px-3 py-1 text-xs font-medium text-primary-foreground">{t('rules.list.active')}</span>
						{:else}
							<span class="inline-flex items-center rounded-full bg-muted px-3 py-1 text-xs font-medium text-muted-foreground">{t('rules.list.inactive')}</span>
						{/if}
						<ChevronRight class="size-4 text-muted-foreground" />
					</div>
				</button>
			{/each}
			{#if rules.length === 0}
				<p class="px-5 py-8 text-center text-sm text-muted-foreground">{t('rules.group.empty')}</p>
			{/if}
		</CardContent>
	</Card>
</div>

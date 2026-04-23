<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { Card, CardContent } from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Plus, ArrowLeft, ChevronRight } from '@lucide/svelte';
	import { ruleService } from '$lib/services/rules.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import type { RuleGroup, Rule } from '$lib/types/domain.js';
	const t = localeStore.t;

	let groupId = $derived($page.params.groupId);
	let isLoading = $state(true);
	let errorMsg = $state('');
	let group = $state<RuleGroup | undefined>(undefined);
	let rules = $state<Rule[]>([]);

	onMount(async () => {
		try {
			group = await ruleService.getGroup(groupId ?? '');
			rules = await ruleService.listRules();
		} catch (e) {
			errorMsg = t('common.error');
			console.error('Failed to load rule group:', e);
		} finally {
			isLoading = false;
		}
	});
</script>

<div class="flex flex-col gap-4">
	{#if isLoading}
		<p class="text-sm text-muted-foreground py-8 text-center">{t('common.loading')}</p>
	{:else if errorMsg}
		<p class="text-sm text-destructive py-8 text-center">{errorMsg}</p>
	{:else}
	<div class="flex flex-wrap items-center justify-between gap-3">
		<div class="flex items-center gap-3">
			<Button variant="ghost" size="sm" onclick={() => goto('/rules')}>
				<ArrowLeft class="size-4" />
				{t('rules.group.back')}
			</Button>
			{#if group}
				<h2 class="text-base font-semibold text-foreground">{group.title}</h2>
				<span class="inline-flex items-center rounded-full bg-muted px-3 py-1 text-xs font-medium text-muted-foreground">{t('rules.group.ruleCount', { count: rules.length })}</span>
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
							{rule.active ? t('rules.list.active') : t('rules.list.inactive')} · Priority #{rule.priority}
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
	{/if}
</div>

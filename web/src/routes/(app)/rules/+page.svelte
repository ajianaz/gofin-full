<script lang="ts">
	import { goto } from '$app/navigation';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Plus, ChevronRight } from '@lucide/svelte';
	import { mockRuleGroups } from '$lib/data/mock-rules.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;
</script>

<div class="flex flex-col gap-4">
	<div class="flex flex-wrap items-center justify-between gap-3">
		<div class="flex items-center gap-3">
			<h2 class="text-lg font-semibold text-foreground">{t('rules.list.title')}</h2>
			<span class="inline-flex items-center rounded-2xl bg-primary px-2.5 py-0.5 text-xs font-medium text-primary-foreground">
				{mockRuleGroups.length}
			</span>
		</div>
		<Button size="sm" onclick={() => goto('/rules/create')}>
			<Plus class="size-4" />
			{t('rules.list.add')}
		</Button>
	</div>

	<Card>
		<CardHeader>
			<CardTitle class="text-base">{t('rules.list.title')}</CardTitle>
			<p class="text-[13px] text-muted-foreground">{t('rules.list.subtitle')}</p>
		</CardHeader>
		<CardContent class="p-0">
			{#each mockRuleGroups as group}
				<button
					type="button"
					class="flex w-full items-center justify-between px-5 py-4 border-b last:border-b-0 hover:bg-muted/50 transition-colors text-left"
					onclick={() => goto('/rules/{group.id}')}
				>
					<div class="flex flex-col gap-1">
						<p class="text-sm font-semibold text-foreground">{group.title}</p>
						<p class="text-[13px] text-muted-foreground">
							{t('rules.list.ruleCount', { count: group.rule_count })} · {group.active ? t('rules.list.active') : t('rules.list.inactive')}
						</p>
					</div>
					<div class="flex shrink-0 items-center gap-3">
						{#if group.active}
							<span class="inline-flex items-center rounded-full bg-primary px-3 py-1 text-xs font-medium text-primary-foreground">{t('rules.list.active')}</span>
						{:else}
							<span class="inline-flex items-center rounded-full bg-muted px-3 py-1 text-xs font-medium text-muted-foreground">{t('rules.list.inactive')}</span>
						{/if}
						<span class="flex size-8 items-center justify-center rounded-md bg-muted text-sm font-semibold text-foreground">{group.order}</span>
						<ChevronRight class="size-4 text-muted-foreground" />
					</div>
				</button>
			{/each}
		</CardContent>
	</Card>
</div>

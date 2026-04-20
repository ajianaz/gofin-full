<script lang="ts">
	import { PageHeader, StatusBadge } from '$lib/components/shared/index.js';
	import { Card, CardContent } from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Plus } from '@lucide/svelte';
	import { mockGroups } from '$lib/data/mock-groups.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;
</script>

<PageHeader title={t('groups.title')} description={t('groups.description')}>
	{#snippet actions()}
		<Button size="sm">
			<Plus class="size-4" />
			{t('groups.newGroup')}
		</Button>
	{/snippet}
</PageHeader>

<Card>
	<CardContent class="p-0">
		<div class="overflow-x-auto">
			<table class="w-full text-sm">
				<thead>
					<tr class="border-b bg-muted/50">
						<th class="text-left p-3 font-medium text-muted-foreground">{t('groups.name')}</th>
						<th class="text-left p-3 font-medium text-muted-foreground">{t('groups.members')}</th>
						<th class="text-left p-3 font-medium text-muted-foreground">{t('groups.activeGroup')}</th>
					</tr>
				</thead>
				<tbody>
					{#each mockGroups as group}
						<tr class="border-b hover:bg-muted/30">
							<td class="p-3 font-medium text-foreground">{group.title}</td>
							<td class="p-3 text-muted-foreground">{t('groups.memberCount', { count: group.member_count })}</td>
							<td class="p-3">
								<StatusBadge status={group.is_current ? 'active' : 'inactive'} />
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</CardContent>
</Card>

<script lang="ts">
	import { PageHeader } from '$lib/components/shared/index.js';
	import { Card, CardContent } from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Plus } from '@lucide/svelte';
	import { mockTags } from '$lib/data/mock-tags.js';
	import { formatDate } from '$lib/utils/format.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;
</script>

<PageHeader title={t('tags.list.title')} description={t('tags.list.description')}>
	{#snippet actions()}
		<Button asChild size="sm">
			<a href="/tags/create">
				<Plus class="size-4" />
				{t('tags.list.add')}
			</a>
		</Button>
	{/snippet}
</PageHeader>

<div class="flex flex-wrap gap-2 mb-4">
	{#each mockTags as tag}
		<Badge variant="outline" class="px-3 py-1.5 text-sm">
			#{tag.tag}
		</Badge>
	{/each}
</div>

<Card>
	<CardContent class="p-0">
		<div class="overflow-x-auto">
			<table class="w-full text-sm">
				<thead>
					<tr class="border-b bg-muted/50">
						<th class="text-left p-3 font-medium text-muted-foreground">{t('tags.list.colTag')}</th>
						<th class="text-left p-3 font-medium text-muted-foreground">{t('tags.list.colDescription')}</th>
						<th class="text-left p-3 font-medium text-muted-foreground">{t('tags.list.colDate')}</th>
					</tr>
				</thead>
				<tbody>
					{#each mockTags as tag}
						<tr class="border-b hover:bg-muted/30">
							<td class="p-3"><Badge variant="secondary">#{tag.tag}</Badge></td>
							<td class="p-3 text-muted-foreground">{tag.description ?? '-'}</td>
							<td class="p-3 text-muted-foreground">{formatDate(tag.date)}</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</CardContent>
</Card>

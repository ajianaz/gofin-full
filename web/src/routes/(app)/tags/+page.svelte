<script lang="ts">
	import { onMount } from 'svelte';
	import { PageHeader } from '$lib/components/shared/index.js';
	import { Card, CardContent } from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Plus, Trash2 } from '@lucide/svelte';
	import { tagService } from '$lib/services/index.js';
	import type { Tag } from '$lib/types/domain.js';
	import { formatDate } from '$lib/utils/format.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import { ConfirmDialog } from '$lib/components/shared/index.js';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	const t = localeStore.t;

	let items = $state<Tag[]>([]);
	let isLoading = $state(true);
	let errorMsg = $state('');
	let deleteTarget = $state<string | null>(null);
	let deleteOpen = $derived(deleteTarget !== null);

	async function handleDelete(id: string) {
		await tagService.delete(id);
		items = items.filter((t) => t.id !== id);
	}

	onMount(async () => {
		try {
			items = await tagService.list();
		} catch (e) {
			errorMsg = t('common.error');
			console.error(e);
		} finally {
			isLoading = false;
		}
	});
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

<Card>
	<CardContent class="p-0">
		<div class="overflow-x-auto">
			<table class="w-full text-sm">
				<thead>
					<tr class="border-b bg-muted/50">
						<th class="text-left p-3 font-medium text-muted-foreground">{t('tags.list.colTag')}</th>
						<th class="text-left p-3 font-medium text-muted-foreground">{t('tags.list.colDescription')}</th>
						<th class="text-left p-3 font-medium text-muted-foreground">{t('tags.list.colDate')}</th>
						<th class="w-[50px] p-3"></th>
					</tr>
				</thead>
				<tbody>
					{#if isLoading}
				{#each Array(5) as _}
				<tr class="border-b">
					<td class="p-3"><Skeleton class="h-4 w-full" /></td>
					<td class="p-3"><Skeleton class="h-4 w-full" /></td>
					<td class="p-3"><Skeleton class="h-4 w-full" /></td>
					<td class="p-3"><Skeleton class="h-4 w-full" /></td>
					<td class="p-3"><Skeleton class="h-4 w-full" /></td>
				</tr>
				<tr class="border-b">
					<td class="p-3"><Skeleton class="h-4 w-full" /></td>
					<td class="p-3"><Skeleton class="h-4 w-full" /></td>
					<td class="p-3"><Skeleton class="h-4 w-full" /></td>
					<td class="p-3"><Skeleton class="h-4 w-full" /></td>
					<td class="p-3"><Skeleton class="h-4 w-full" /></td>
				</tr>
				<tr class="border-b">
					<td class="p-3"><Skeleton class="h-4 w-full" /></td>
					<td class="p-3"><Skeleton class="h-4 w-full" /></td>
					<td class="p-3"><Skeleton class="h-4 w-full" /></td>
					<td class="p-3"><Skeleton class="h-4 w-full" /></td>
					<td class="p-3"><Skeleton class="h-4 w-full" /></td>
				</tr>
				<tr class="border-b">
					<td class="p-3"><Skeleton class="h-4 w-full" /></td>
					<td class="p-3"><Skeleton class="h-4 w-full" /></td>
					<td class="p-3"><Skeleton class="h-4 w-full" /></td>
					<td class="p-3"><Skeleton class="h-4 w-full" /></td>
					<td class="p-3"><Skeleton class="h-4 w-full" /></td>
				</tr>
				<tr class="border-b">
					<td class="p-3"><Skeleton class="h-4 w-full" /></td>
					<td class="p-3"><Skeleton class="h-4 w-full" /></td>
					<td class="p-3"><Skeleton class="h-4 w-full" /></td>
					<td class="p-3"><Skeleton class="h-4 w-full" /></td>
					<td class="p-3"><Skeleton class="h-4 w-full" /></td>
				</tr>
				{/each}
	{:else if errorMsg}
						<tr>
							<td colspan="4" class="p-8 text-center text-sm text-destructive">{errorMsg}</td>
						</tr>
					{:else}
					{#each items as tag}
						<tr class="border-b hover:bg-muted/30">
							<td class="p-3"><Badge variant="secondary">#{tag.tag}</Badge></td>
							<td class="p-3 text-muted-foreground">{tag.description ?? '-'}</td>
							<td class="p-3 text-muted-foreground">{formatDate(tag.date)}</td>
							<td class="p-3">
								<button type="button" aria-label="{t('common.delete')}" class="text-muted-foreground hover:text-destructive transition-colors" onclick={() => (deleteTarget = tag.id)}>
									<Trash2 class="size-4" />
								</button>
							</td>
						</tr>
					{:else}
						<tr><td colspan="4"><EmptyState /></td></tr>
					{/each}
					{/if}
				</tbody>
			</table>
	<ConfirmDialog
		bind:open={deleteOpen}
		title={t('common.delete')}
		description={t('common.deleteConfirm')}
		onConfirm={async () => {
			if (deleteTarget) {
				await handleDelete(deleteTarget);
				deleteTarget = null;
			}
		}}
	/>
</div>
	</CardContent>
</Card>

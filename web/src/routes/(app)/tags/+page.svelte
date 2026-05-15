<script lang="ts">
	import { onMount } from 'svelte';
	import { PageHeader } from '$lib/components/shared/index.js';
	import { Card, CardContent } from '$lib/components/ui/card/index.js';
	import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '$lib/components/ui/table/index.js';
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
		<Table>
				<TableHeader>
					<TableRow>
						<TableHead>{t('tags.list.colTag')}</TableHead>
						<TableHead>{t('tags.list.colDescription')}</TableHead>
						<TableHead class="hidden md:table-cell">{t('tags.list.colDate')}</TableHead>
						<TableHead class="w-[50px]"></TableHead>
					</TableRow>
				</TableHeader>
				<TableBody>
					{#if isLoading}
				{#each Array(5) as _}
<TableRow>
					<TableCell><Skeleton class="h-4 w-full" /></TableCell>
					<TableCell><Skeleton class="h-4 w-full" /></TableCell>
					<TableCell class="hidden md:table-cell"><Skeleton class="h-4 w-full" /></TableCell>
					<TableCell><Skeleton class="h-4 w-full" /></TableCell>
				</TableRow>
				{/each}
	{:else if errorMsg}
						<TableRow>
							<TableCell colspan="4" class="p-8 text-center text-sm text-destructive">{errorMsg}</TableCell>
						</TableRow>
					{:else}
					{#each items as tag}
						<TableRow>
							<TableCell><Badge variant="secondary">#{tag.tag}</Badge></TableCell>
							<TableCell class="text-muted-foreground">{tag.description ?? '-'}</TableCell>
							<TableCell class="hidden md:table-cell text-muted-foreground">{formatDate(tag.date)}</TableCell>
							<TableCell>
								<button type="button" aria-label="{t('common.delete')}" class="text-muted-foreground hover:text-destructive transition-colors" onclick={() => (deleteTarget = tag.id)}>
									<Trash2 class="size-4" />
								</button>
							</TableCell>
						</TableRow>
					{:else}
						<TableRow><TableCell colspan="4"><EmptyState /></TableCell></TableRow>
					{/each}
					{/if}
				</TableBody>
			</Table>
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
	</CardContent>
</Card>

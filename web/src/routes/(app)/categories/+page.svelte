<script lang="ts">
	import { onMount } from 'svelte';
	import { PageHeader } from '$lib/components/shared/index.js';
	import { Card, CardContent } from '$lib/components/ui/card/index.js';
	import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '$lib/components/ui/table/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Plus, Trash2 } from '@lucide/svelte';
	import { categoryService } from '$lib/services/index.js';
	import type { Category } from '$lib/types/domain.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import { ConfirmDialog } from '$lib/components/shared/index.js';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	const t = localeStore.t;

	let items = $state<Category[]>([]);
	let isLoading = $state(true);
	let errorMsg = $state('');
	let deleteTarget = $state<string | null>(null);
	let deleteOpen = $derived(deleteTarget !== null);

	async function handleDelete(id: string) {
		await categoryService.delete(id);
		items = items.filter((c) => c.id !== id);
	}

	onMount(async () => {
		try {
			items = await categoryService.list();
		} catch (e) {
			errorMsg = t('common.error');
			console.error(e);
		} finally {
			isLoading = false;
		}
	});
</script>

<PageHeader title={t('categories.list.title')} description={t('categories.list.description')}>
	{#snippet actions()}
		<Button asChild size="sm">
			<a href="/categories/create">
				<Plus class="size-4" />
				{t('categories.list.add')}
			</a>
		</Button>
	{/snippet}
</PageHeader>

<Card>
	<CardContent class="p-0">
		<Table>
				<TableHeader>
					<TableRow>
						<TableHead>{t('categories.list.colName')}</TableHead>
						<TableHead>{t('categories.list.colType')}</TableHead>
						<TableHead class="hidden md:table-cell">{t('categories.list.colTransactions')}</TableHead>
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
					{#each items as cat}
						<TableRow>
							<TableCell class="font-medium text-foreground">{cat.name}</TableCell>
							<TableCell>
								<Badge variant={cat.type === 'income' ? 'default' : 'secondary'}>
									{cat.type === 'expense' ? t('categories.list.expense') : cat.type === 'income' ? t('categories.list.income') : t('categories.list.transfer')}
								</Badge>
							</TableCell>
							<TableCell class="hidden md:table-cell text-muted-foreground">{cat.transaction_count}</TableCell>
							<TableCell>
								<button type="button" aria-label="{t('common.delete')}" class="text-muted-foreground hover:text-destructive transition-colors" onclick={() => (deleteTarget = cat.id)}>
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

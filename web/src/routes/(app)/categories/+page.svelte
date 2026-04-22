<script lang="ts">
	import { onMount } from 'svelte';
	import { PageHeader } from '$lib/components/shared/index.js';
	import { Card, CardContent } from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Plus, Trash2 } from '@lucide/svelte';
	import { categoryService } from '$lib/services/index.js';
	import type { Category } from '$lib/types/domain.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	let items = $state<Category[]>([]);
	let isLoading = $state(true);

	async function handleDelete(id: string) {
		if (!confirm('Hapus kategori ini?')) return;
		await categoryService.delete(id);
		items = items.filter((c) => c.id !== id);
	}

	onMount(async () => {
		try {
			items = await categoryService.list();
		} catch (e) {
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
		<div class="overflow-x-auto">
			<table class="w-full text-sm">
				<thead>
					<tr class="border-b bg-muted/50">
						<th class="text-left p-3 font-medium text-muted-foreground">{t('categories.list.colName')}</th>
						<th class="text-left p-3 font-medium text-muted-foreground">{t('categories.list.colType')}</th>
						<th class="text-left p-3 font-medium text-muted-foreground">{t('categories.list.colTransactions')}</th>
						<th class="w-[50px] p-3"></th>
					</tr>
				</thead>
				<tbody>
					{#if isLoading}
					<tr>
						<td colspan="4" class="p-8 text-center text-sm text-muted-foreground">Memuat...</td>
					</tr>
					{:else}
					{#each items as cat}
						<tr class="border-b hover:bg-muted/30">
							<td class="p-3 font-medium text-foreground">{cat.name}</td>
							<td class="p-3">
								<Badge variant={cat.type === 'income' ? 'default' : 'secondary'}>
									{cat.type === 'expense' ? t('categories.list.expense') : cat.type === 'income' ? t('categories.list.income') : t('categories.list.transfer')}
								</Badge>
							</td>
							<td class="p-3 text-muted-foreground">{cat.transaction_count}</td>
							<td class="p-3">
								<button type="button" class="text-muted-foreground hover:text-destructive transition-colors" onclick={() => handleDelete(cat.id)}>
									<Trash2 class="size-4" />
								</button>
							</td>
						</tr>
					{/each}
					{/if}
				</tbody>
			</table>
		</div>
	</CardContent>
</Card>

<script lang="ts">
	import { onMount } from 'svelte';
	import { Card, CardContent } from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Plus, Copy, Trash2, AlertTriangle } from '@lucide/svelte';
	import { apiKeyService } from '$lib/services/api-keys.js';
	import { formatDate } from '$lib/utils/format.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import type { ApiKeyListItem } from '$lib/types/domain.js';
	import { ConfirmDialog } from '$lib/components/shared/index.js';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	const t = localeStore.t;

	let apiKeys = $state<ApiKeyListItem[]>([]);
	let isLoading = $state(true);
	let errorMsg = $state('');
	let isCreating = $state(false);
	let deleteTarget = $state<string | null>(null);
	let deleteOpen = $derived(deleteTarget !== null);
	let isDeleting = $state<string | null>(null);

	async function loadApiKeys() {
		isLoading = true;
		errorMsg = '';
		try {
			apiKeys = await apiKeyService.list();
		} catch (e) {
			errorMsg = t('common.error');
			console.error(e);
		} finally {
			isLoading = false;
		}
	}

	async function handleCreate() {
		const name = prompt(t('settings.apiKeys.name'));
		if (!name || !name.trim()) return;

		isCreating = true;
		try {
			const result = await apiKeyService.create(name.trim());
			alert(result.key || t('settings.apiKeys.key'));
			await loadApiKeys();
		} catch (e) {
			errorMsg = t('common.errorSave');
			console.error(e);
		} finally {
			isCreating = false;
		}
	}

	async function handleDelete(id: string) {
	
		isDeleting = id;
		try {
		await apiKeyService.delete(id);
		apiKeys = apiKeys.filter((k) => k.id !== id);
		await loadApiKeys();
		} catch (e) {
		errorMsg = t('common.errorSave');
		console.error(e);
		} finally {
		isDeleting = null;
		}
	}

	async function handleCopy(prefix: string) {
		try {
			await navigator.clipboard.writeText(prefix);
		} catch {
			console.error('Failed to copy');
		}
	}

	onMount(loadApiKeys);
</script>

<div class="flex flex-col gap-4">
	<div class="flex flex-wrap items-center justify-between gap-3">
		<h2 class="text-lg font-semibold text-foreground">{t('settings.apiKeys.title')}</h2>
		<Button size="sm" onclick={handleCreate} disabled={isCreating}>
			<Plus class="size-4" />
			{isCreating ? t('common.saving') : t('settings.apiKeys.createNew')}
		</Button>
	</div>

	<div class="flex items-start gap-3 rounded-lg border border-destructive/50 bg-destructive/5 p-4">
		<AlertTriangle class="size-5 shrink-0 text-destructive mt-0.5" />
		<div class="flex flex-col gap-1">
			<p class="text-sm font-semibold text-destructive">{t('settings.apiKeys.securityWarning')}</p>
			<p class="text-sm text-muted-foreground">{t('settings.apiKeys.securityWarningDesc')}</p>
		</div>
	</div>

	{#if errorMsg}
		<div class="rounded-lg border border-destructive/50 bg-destructive/5 px-4 py-3">
			<p class="text-sm text-destructive">{errorMsg}</p>
		</div>
	{/if}

	<Card>
		<CardContent class="p-0">
			<div class="overflow-x-auto">
				<table class="w-full text-sm">
					<thead>
						<tr class="border-b bg-muted/50">
							<th class="text-left px-4 py-3 font-medium text-muted-foreground">{t('settings.apiKeys.name')}</th>
							<th class="text-left px-4 py-3 font-medium text-muted-foreground">{t('settings.apiKeys.key')}</th>
							<th class="text-left px-4 py-3 font-medium text-muted-foreground w-40">{t('settings.apiKeys.created')}</th>
							<th class="text-right px-4 py-3 font-medium text-muted-foreground w-48">{t('common.actions')}</th>
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
	{:else if apiKeys.length === 0}
							<tr>
								<td colspan="4"><EmptyState /></td>
							</tr>
						{:else}
							{#each apiKeys as key (key.id)}
								<tr class="border-b last:border-b-0 hover:bg-muted/30">
									<td class="px-4 py-3 font-medium text-foreground">{key.name}</td>
									<td class="px-4 py-3 font-mono text-sm text-muted-foreground">{key.key_prefix}...</td>
									<td class="px-4 py-3 text-muted-foreground">{formatDate(key.created_at)}</td>
									<td class="px-4 py-3 text-right">
										<div class="inline-flex items-center gap-2">
											<Button variant="ghost" size="sm" onclick={() => handleCopy(key.key_prefix)}>
												<Copy class="size-4" />
											</Button>
											<Button variant="ghost" size="sm" class="text-destructive" onclick={() => (deleteTarget = key.id)} disabled={isDeleting === key.id}>
												<Trash2 class="size-4" />
											</Button>
										</div>
									</td>
								</tr>
							{/each}
						{/if}
					</tbody>
				</table>
			</div>
		</CardContent>
	</Card>
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

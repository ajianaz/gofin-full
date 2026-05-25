<script lang="ts">
	import { onMount } from 'svelte';
	import { authStore } from '$lib/stores/auth.svelte.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import { groupService } from '$lib/services/index.js';
	import PageHeader from '$lib/components/shared/PageHeader.svelte';
	import { Card, CardContent } from '$lib/components/ui/card/index.js';
	import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '$lib/components/ui/table/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Plus, Users, Check } from '@lucide/svelte';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	import type { UserGroup } from '$lib/types/domain.js';
	const t = localeStore.t;

	let groups = $state<UserGroup[]>([]);
	let isLoading = $state(true);
	let errorMsg = $state('');
	let newTitle = $state('');
	let isCreating = $state(false);

	async function loadGroups() {
		isLoading = true;
		errorMsg = '';
		try {
			groups = await groupService.list();
		} catch (e) {
			errorMsg = t('common.error');
			console.error(e);
		} finally {
			isLoading = false;
		}
	}

	async function handleCreate() {
		if (!newTitle.trim()) return;
		isCreating = true;
		try {
			await groupService.create({ title: newTitle.trim() });
			newTitle = '';
			await loadGroups();
		} catch (e) {
			errorMsg = t('common.errorSave');
			console.error(e);
		} finally {
			isCreating = false;
		}
	}

	async function handleSwitch(groupId: string) {
		try {
			const tokens = await groupService.switch(groupId);
			if (tokens) {
				localStorage.setItem('access_token', tokens.access_token);
				if (tokens.refresh_token) {
					localStorage.setItem('refresh_token', tokens.refresh_token);
				}
				await authStore.fetchUser();
				await loadGroups();
			}
		} catch (e) {
			errorMsg = t('common.errorSave');
			console.error(e);
		}
	}

	onMount(loadGroups);
</script>

<div class="flex flex-col gap-6">
	<PageHeader title={t('settings.groups.title')} description={t('settings.groups.description')} icon={Users} />

	{#if errorMsg}
		<div class="rounded-lg border border-destructive/50 bg-destructive/5 px-4 py-3">
			<p class="text-sm text-destructive">{errorMsg}</p>
		</div>
	{/if}

	<div class="flex items-end gap-3 max-w-md">
		<div class="flex-1 grid gap-2">
			<label for="new-group" class="text-sm font-medium text-foreground">{t('settings.groups.createLabel')}</label>
			<Input id="new-group" placeholder={t('settings.groups.createPlaceholder')} bind:value={newTitle} />
		</div>
		<Button size="sm" onclick={handleCreate} disabled={isCreating || !newTitle.trim()}>
			<Plus class="size-4" />
			{t('common.create')}
		</Button>
	</div>

	<Card>
		<CardContent class="p-0">
			<Table>
				<TableHeader>
					<TableRow>
						<TableHead>{t('settings.groups.colName')}</TableHead>
						<TableHead class="text-center">{t('settings.groups.colMembers')}</TableHead>
						<TableHead class="text-right">{t('common.actions')}</TableHead>
					</TableRow>
				</TableHeader>
				<TableBody>
					{#if isLoading}
						{#each Array(3) as _}
							<TableRow>
								<TableCell><Skeleton class="h-4 w-40" /></TableCell>
								<TableCell><Skeleton class="h-4 w-12 mx-auto" /></TableCell>
								<TableCell><Skeleton class="h-4 w-20 ml-auto" /></TableCell>
							</TableRow>
						{/each}
					{:else if groups.length === 0}
						<TableRow>
							<TableCell colspan={3}><EmptyState /></TableCell>
						</TableRow>
					{:else}
						{#each groups as group (group.id)}
							<TableRow>
								<TableCell class="font-medium text-foreground">
									<div class="flex items-center gap-2">
										{group.title}
										{#if group.is_current}
											<span class="inline-flex items-center rounded-full bg-primary px-2 py-0.5 text-xs font-medium text-primary-foreground">
												<Check class="size-3 mr-0.5" />
												{t('settings.groups.active')}
											</span>
										{/if}
									</div>
								</TableCell>
								<TableCell class="text-center text-muted-foreground">{group.member_count ?? 0}</TableCell>
								<TableCell class="text-right">
									{#if !group.is_current}
										<Button variant="outline" size="sm" onclick={() => handleSwitch(group.id)}>
											{t('settings.groups.switchTo')}
										</Button>
									{/if}
								</TableCell>
							</TableRow>
						{/each}
					{/if}
				</TableBody>
			</Table>
		</CardContent>
	</Card>
</div>

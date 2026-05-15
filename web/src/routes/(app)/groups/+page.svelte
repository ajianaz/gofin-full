<script lang="ts">
	import { onMount } from 'svelte';
	import { PageHeader, StatusBadge } from '$lib/components/shared/index.js';
	import { Card, CardContent } from '$lib/components/ui/card/index.js';
	import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '$lib/components/ui/table/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Plus } from '@lucide/svelte';
	import { groupService } from '$lib/services/index.js';
	import { authStore } from '$lib/stores/auth.svelte.js';
	import type { UserGroup } from '$lib/types/domain.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	const t = localeStore.t;

	let groups = $state<UserGroup[]>([]);
	let isLoading = $state(true);
	let errorMsg = $state('');

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

	async function handleSwitch(id: string) {
		try {
			const tokens = await groupService.switch(id);
			if (tokens) {
				authStore.setTokens(tokens);
			}
			await loadGroups();
		} catch (e) {
			errorMsg = t('common.error');
			console.error(e);
		}
	}

	onMount(loadGroups);
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
		<Table>
				<TableHeader>
					<TableRow>
						<TableHead>{t('groups.name')}</TableHead>
						<TableHead>{t('groups.members')}</TableHead>
						<TableHead>{t('groups.activeGroup')}</TableHead>
						<TableHead class="w-[100px]"></TableHead>
					</TableRow>
				</TableHeader>
				<TableBody>
					{#if isLoading}
				{#each Array(5) as _}
<TableRow>
					<TableCell><Skeleton class="h-4 w-full" /></TableCell>
					<TableCell><Skeleton class="h-4 w-full" /></TableCell>
					<TableCell><Skeleton class="h-4 w-full" /></TableCell>
					<TableCell><Skeleton class="h-4 w-full" /></TableCell>
				</TableRow>
				{/each}
	{:else if errorMsg}
						<TableRow>
							<TableCell colspan="4" class="p-8 text-center text-sm text-destructive">{errorMsg}</TableCell>
						</TableRow>
					{:else}
						{#each groups as group}
							<TableRow>
								<TableCell class="font-medium text-foreground">{group.title}</TableCell>
								<TableCell class="text-muted-foreground">{t('groups.memberCount', { count: group.member_count })}</TableCell>
								<TableCell>
									<StatusBadge status={group.is_current ? 'active' : 'inactive'} />
								</TableCell>
								<TableCell>
									{#if !group.is_current}
										<Button variant="outline" size="sm" onclick={() => handleSwitch(group.id)}>
											{t('groups.switch')}
										</Button>
									{/if}
								</TableCell>
							</TableRow>
						{:else}
							<TableRow>
								<TableCell colspan="4"><EmptyState /></TableCell>
							</TableRow>
						{/each}
					{/if}
				</TableBody>
			</Table>
	</CardContent>
</Card>

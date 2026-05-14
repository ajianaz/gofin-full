<script lang="ts">
	import { onMount } from 'svelte';
	import { PageHeader, StatusBadge } from '$lib/components/shared/index.js';
	import { Card, CardContent } from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Plus } from '@lucide/svelte';
	import { groupService } from '$lib/services/index.js';
	import { authStore } from '$lib/stores/auth.svelte.js';
	import type { UserGroup } from '$lib/types/domain.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
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
		<div class="overflow-x-auto">
			<table class="w-full text-sm">
				<thead>
					<tr class="border-b bg-muted/50">
						<th class="text-left p-3 font-medium text-muted-foreground">{t('groups.name')}</th>
						<th class="text-left p-3 font-medium text-muted-foreground">{t('groups.members')}</th>
						<th class="text-left p-3 font-medium text-muted-foreground">{t('groups.activeGroup')}</th>
						<th class="w-[100px] p-3"></th>
					</tr>
				</thead>
				<tbody>
					{#if isLoading}
						<tr>
							<td colspan="4" class="p-8 text-center text-sm text-muted-foreground">{t('common.loading')}</td>
						</tr>
					{:else if errorMsg}
						<tr>
							<td colspan="4" class="p-8 text-center text-sm text-destructive">{errorMsg}</td>
						</tr>
					{:else}
						{#each groups as group}
							<tr class="border-b hover:bg-muted/30">
								<td class="p-3 font-medium text-foreground">{group.title}</td>
								<td class="p-3 text-muted-foreground">{t('groups.memberCount', { count: group.member_count })}</td>
								<td class="p-3">
									<StatusBadge status={group.is_current ? 'active' : 'inactive'} />
								</td>
								<td class="p-3">
									{#if !group.is_current}
										<Button variant="outline" size="sm" onclick={() => handleSwitch(group.id)}>
											{t('groups.switch')}
										</Button>
									{/if}
								</td>
							</tr>
						{:else}
							<tr>
								<td colspan="4" class="p-8 text-center text-sm text-muted-foreground">{t('common.noData')}</td>
							</tr>
						{/each}
					{/if}
				</tbody>
			</table>
		</div>
	</CardContent>
</Card>

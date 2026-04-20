<script lang="ts">
	import { PageHeader, StatusBadge } from '$lib/components/shared/index.js';
	import { Card, CardContent } from '$lib/components/ui/card/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { mockUsers } from '$lib/data/mock-users.js';
	import { formatDate } from '$lib/utils/format.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	let roleLabels = $derived<Record<string, string>>({
		admin: t('admin.users.roleAdmin'),
		manager: t('admin.users.roleManager'),
		user: t('admin.users.roleUser')
	});
</script>

<PageHeader title={t('admin.users.title')} description={t('admin.users.description')} />

<Card>
	<CardContent class="p-0">
		<div class="overflow-x-auto">
			<table class="w-full text-sm">
				<thead>
					<tr class="border-b bg-muted/50">
						<th class="text-left p-3 font-medium text-muted-foreground">{t('admin.users.email')}</th>
						<th class="text-left p-3 font-medium text-muted-foreground">{t('admin.users.name')}</th>
						<th class="text-left p-3 font-medium text-muted-foreground">{t('admin.users.role')}</th>
						<th class="text-left p-3 font-medium text-muted-foreground">{t('common.status')}</th>
						<th class="text-left p-3 font-medium text-muted-foreground">{t('admin.users.joined')}</th>
					</tr>
				</thead>
				<tbody>
					{#each mockUsers as user}
						<tr class="border-b hover:bg-muted/30">
							<td class="p-3 text-foreground">{user.email}</td>
							<td class="p-3 font-medium text-foreground">{user.name}</td>
							<td class="p-3">
								<Badge variant="outline">{roleLabels[user.role] ?? user.role}</Badge>
							</td>
							<td class="p-3"><StatusBadge status={user.is_active ? 'active' : 'inactive'} /></td>
							<td class="p-3 text-muted-foreground">{formatDate(user.created_at)}</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</CardContent>
</Card>

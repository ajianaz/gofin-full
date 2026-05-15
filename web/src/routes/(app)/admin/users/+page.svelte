<script lang="ts">
	import { onMount } from 'svelte';
	import { PageHeader, StatusBadge } from '$lib/components/shared/index.js';
	import { Card, CardContent } from '$lib/components/ui/card/index.js';
	import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '$lib/components/ui/table/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { adminService } from '$lib/services/index.js';
	import { formatDate } from '$lib/utils/format.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	const t = localeStore.t;

	interface AdminUser {
		id: string;
		email: string;
		name: string;
		role: string;
		is_active: boolean;
		created_at: string;
	}

	let users = $state<AdminUser[]>([]);
	let isLoading = $state(true);
	let errorMsg = $state('');

	let roleLabels = $derived<Record<string, string>>({
		admin: t('admin.users.roleAdmin'),
		manager: t('admin.users.roleManager'),
		user: t('admin.users.roleUser')
	});

	onMount(async () => {
		try {
			users = await adminService.listUsers();
		} catch (e) {
			errorMsg = t('common.error');
			console.error(e);
		} finally {
			isLoading = false;
		}
	});
</script>

<PageHeader title={t('admin.users.title')} description={t('admin.users.description')} />

<Card>
	<CardContent class="p-0">
		<Table>
				<TableHeader>
					<TableRow>
						<TableHead>{t('admin.users.email')}</TableHead>
						<TableHead>{t('admin.users.name')}</TableHead>
						<TableHead>{t('admin.users.role')}</TableHead>
						<TableHead class="hidden md:table-cell">{t('common.status')}</TableHead>
						<TableHead>{t('admin.users.joined')}</TableHead>
					</TableRow>
				</TableHeader>
				<TableBody>
					{#if isLoading}
				{#each Array(5) as _}
<TableRow>
					<TableCell><Skeleton class="h-4 w-full" /></TableCell>
					<TableCell><Skeleton class="h-4 w-full" /></TableCell>
					<TableCell><Skeleton class="h-4 w-full" /></TableCell>
					<TableCell class="hidden md:table-cell"><Skeleton class="h-4 w-full" /></TableCell>
					<TableCell><Skeleton class="h-4 w-full" /></TableCell>
				</TableRow>
				{/each}
	{:else if errorMsg}
						<TableRow>
							<TableCell colspan="5" class="py-8 text-center text-sm text-destructive">{errorMsg}</TableCell>
						</TableRow>
					{:else}
						{#each users as user}
							<TableRow>
								<TableCell class="text-foreground">{user.email}</TableCell>
								<TableCell class="font-medium text-foreground">{user.name}</TableCell>
								<TableCell>
									<Badge variant="outline">{roleLabels[user.role] ?? user.role}</Badge>
								</TableCell>
								<TableCell class="hidden md:table-cell"><StatusBadge status={user.is_active ? 'active' : 'inactive'} /></TableCell>
								<TableCell class="text-muted-foreground">{formatDate(user.created_at)}</TableCell>
							</TableRow>
						{:else}
							<TableRow>
								<TableCell colspan="5"><EmptyState /></TableCell>
							</TableRow>
						{/each}
					{/if}
				</TableBody>
			</Table>
	</CardContent>
</Card>

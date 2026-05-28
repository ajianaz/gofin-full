<script lang="ts">
	import { onMount } from 'svelte';
	import { PageHeader, FilterBar } from '$lib/components/shared/index.js';
	import { Card, CardContent } from '$lib/components/ui/card/index.js';
	import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '$lib/components/ui/table/index.js';
	import { Select, SelectTrigger, SelectContent, SelectItem, SelectValue } from '$lib/components/ui/select/index.js';
	import { adminService } from '$lib/services/index.js';
	import { formatDate } from '$lib/utils/format.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	const t = localeStore.t;

	interface AuditLogEntry {
		id: string;
		action: string;
		user_email: string;
		entity_type: string;
		entity_id: string;
		changes: string;
		created_at: string;
	}

	let logs = $state<AuditLogEntry[]>([]);
	let isLoading = $state(true);
	let errorMsg = $state('');

	let actionFilter = $state('all');
	let entityFilter = $state('all');

	let filtered = $derived(() => {
		if (actionFilter === 'all') return logs;
		return logs.filter((l) => l.action.startsWith(actionFilter));
	});

	async function fetchLogs(entityType?: string) {
		isLoading = true;
		errorMsg = '';
		try {
			logs = await adminService.listAuditLogs(entityType);
		} catch (e) {
			errorMsg = t('common.error');
			console.error(e);
		} finally {
			isLoading = false;
		}
	}

	onMount(() => fetchLogs());

	$effect(() => {
		if (entityFilter === 'all') {
			fetchLogs();
		} else {
			fetchLogs(entityFilter);
		}
	});
</script>

<PageHeader title={t('admin.auditLog.title')} description={t('admin.auditLog.description')} />

<FilterBar>
	<Select bind:value={actionFilter} items={[
		{ value: 'all', label: t('admin.auditLog.allActions') },
		{ value: 'user.login', label: t('admin.auditLog.login') },
		{ value: 'transaction', label: t('admin.auditLog.transaction') },
		{ value: 'budget', label: t('admin.auditLog.budget') },
		{ value: 'api_key', label: t('admin.auditLog.apiKey') },
		{ value: 'user.update', label: t('admin.auditLog.userUpdate') },
		{ value: 'group.create', label: t('admin.auditLog.group') },
		{ value: 'piggy_bank', label: t('admin.auditLog.piggyBank') },
		{ value: 'recurring', label: t('admin.auditLog.recurring') },
		{ value: 'rule', label: t('admin.auditLog.rule') },
		{ value: 'currency', label: t('admin.auditLog.currency') },
	]}>
		<SelectTrigger class="w-40">
			<SelectValue />
		</SelectTrigger>
		<SelectContent>
			<SelectItem value="all" label={t('admin.auditLog.allActions')}>{t('admin.auditLog.allActions')}</SelectItem>
			<SelectItem value="user.login" label={t('admin.auditLog.login')}>{t('admin.auditLog.login')}</SelectItem>
			<SelectItem value="transaction" label={t('admin.auditLog.transaction')}>{t('admin.auditLog.transaction')}</SelectItem>
			<SelectItem value="budget" label={t('admin.auditLog.budget')}>{t('admin.auditLog.budget')}</SelectItem>
			<SelectItem value="api_key" label={t('admin.auditLog.apiKey')}>{t('admin.auditLog.apiKey')}</SelectItem>
			<SelectItem value="user.update" label={t('admin.auditLog.userUpdate')}>{t('admin.auditLog.userUpdate')}</SelectItem>
			<SelectItem value="group.create" label={t('admin.auditLog.group')}>{t('admin.auditLog.group')}</SelectItem>
			<SelectItem value="piggy_bank" label={t('admin.auditLog.piggyBank')}>{t('admin.auditLog.piggyBank')}</SelectItem>
			<SelectItem value="recurring" label={t('admin.auditLog.recurring')}>{t('admin.auditLog.recurring')}</SelectItem>
			<SelectItem value="rule" label={t('admin.auditLog.rule')}>{t('admin.auditLog.rule')}</SelectItem>
			<SelectItem value="currency" label={t('admin.auditLog.currency')}>{t('admin.auditLog.currency')}</SelectItem>
		</SelectContent>
	</Select>
	<Select bind:value={entityFilter} items={[
		{ value: 'all', label: t('admin.auditLog.allEntities') },
		{ value: 'user', label: t('admin.auditLog.user') },
		{ value: 'transaction', label: t('admin.auditLog.transaction') },
		{ value: 'budget', label: t('admin.auditLog.budget') },
		{ value: 'api_key', label: t('admin.auditLog.apiKey') },
		{ value: 'group', label: t('admin.auditLog.group') },
		{ value: 'piggy_bank', label: t('admin.auditLog.piggyBank') },
		{ value: 'recurring', label: t('admin.auditLog.recurring') },
		{ value: 'rule', label: t('admin.auditLog.rule') },
		{ value: 'currency', label: t('admin.auditLog.currency') },
	]}>
		<SelectTrigger class="w-40">
			<SelectValue />
		</SelectTrigger>
		<SelectContent>
			<SelectItem value="all" label={t('admin.auditLog.allEntities')}>{t('admin.auditLog.allEntities')}</SelectItem>
			<SelectItem value="user" label={t('admin.auditLog.user')}>{t('admin.auditLog.user')}</SelectItem>
			<SelectItem value="transaction" label={t('admin.auditLog.transaction')}>{t('admin.auditLog.transaction')}</SelectItem>
			<SelectItem value="budget" label={t('admin.auditLog.budget')}>{t('admin.auditLog.budget')}</SelectItem>
			<SelectItem value="api_key" label={t('admin.auditLog.apiKey')}>{t('admin.auditLog.apiKey')}</SelectItem>
			<SelectItem value="group" label={t('admin.auditLog.group')}>{t('admin.auditLog.group')}</SelectItem>
			<SelectItem value="piggy_bank" label={t('admin.auditLog.piggyBank')}>{t('admin.auditLog.piggyBank')}</SelectItem>
			<SelectItem value="recurring" label={t('admin.auditLog.recurring')}>{t('admin.auditLog.recurring')}</SelectItem>
			<SelectItem value="rule" label={t('admin.auditLog.rule')}>{t('admin.auditLog.rule')}</SelectItem>
			<SelectItem value="currency" label={t('admin.auditLog.currency')}>{t('admin.auditLog.currency')}</SelectItem>
		</SelectContent>
	</Select>
</FilterBar>

<Card>
	<CardContent class="p-0">
		<Table>
				<TableHeader>
					<TableRow>
						<TableHead>{t('admin.auditLog.time')}</TableHead>
						<TableHead>{t('admin.auditLog.user')}</TableHead>
						<TableHead>{t('admin.auditLog.action')}</TableHead>
						<TableHead class="hidden md:table-cell">{t('admin.auditLog.entity')}</TableHead>
						<TableHead>{t('admin.auditLog.changes')}</TableHead>
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
							<TableCell colspan={5} class="py-8 text-center text-sm text-destructive">{errorMsg}</TableCell>
						</TableRow>
					{:else}
						{#each filtered() as log}
							<TableRow>
								<TableCell class="text-muted-foreground whitespace-nowrap">{formatDate(log.created_at)}</TableCell>
								<TableCell class="text-foreground">{log.user_email}</TableCell>
								<TableCell>
									<span class="inline-flex items-center rounded-md bg-muted px-1.5 py-0.5 text-xs">{log.action}</span>
								</TableCell>
								<TableCell class="hidden md:table-cell text-muted-foreground">{log.entity_type} ({log.entity_id})</TableCell>
								<TableCell class="text-foreground max-w-xs truncate">{log.changes}</TableCell>
							</TableRow>
						{:else}
							<TableRow>
								<TableCell colspan={5}><EmptyState /></TableCell>
							</TableRow>
						{/each}
					{/if}
				</TableBody>
			</Table>
	</CardContent>
</Card>

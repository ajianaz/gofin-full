<script lang="ts">
	import { PageHeader, FilterBar } from '$lib/components/shared/index.js';
	import { Card, CardContent } from '$lib/components/ui/card/index.js';
	import { Select } from '$lib/components/ui/select/index.js';
	import { mockAuditLog } from '$lib/data/mock-audit-log.js';
	import { formatDate } from '$lib/utils/format.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	let actionFilter = $state('all');
	let entityFilter = $state('all');

	let filtered = $derived(() => {
		let result = [...mockAuditLog];
		if (actionFilter !== 'all') result = result.filter((l) => l.action.startsWith(actionFilter));
		if (entityFilter !== 'all') result = result.filter((l) => l.entity_type === entityFilter);
		return result;
	});
</script>

<PageHeader title={t('admin.auditLog.title')} description={t('admin.auditLog.description')} />

<FilterBar>
	<Select bind:value={actionFilter} class="w-40">
		<option value="all">{t('admin.auditLog.allActions')}</option>
		<option value="user.login">{t('admin.auditLog.login')}</option>
		<option value="transaction">{t('admin.auditLog.transaction')}</option>
		<option value="budget">{t('admin.auditLog.budget')}</option>
		<option value="api_key">{t('admin.auditLog.apiKey')}</option>
		<option value="user.update">{t('admin.auditLog.userUpdate')}</option>
		<option value="group.create">{t('admin.auditLog.group')}</option>
		<option value="piggy_bank">{t('admin.auditLog.piggyBank')}</option>
		<option value="recurring">{t('admin.auditLog.recurring')}</option>
		<option value="rule">{t('admin.auditLog.rule')}</option>
		<option value="currency">{t('admin.auditLog.currency')}</option>
	</Select>
	<Select bind:value={entityFilter} class="w-40">
		<option value="all">{t('admin.auditLog.allEntities')}</option>
		<option value="user">{t('admin.auditLog.user')}</option>
		<option value="transaction">{t('admin.auditLog.transaction')}</option>
		<option value="budget">{t('admin.auditLog.budget')}</option>
		<option value="api_key">{t('admin.auditLog.apiKey')}</option>
		<option value="group">{t('admin.auditLog.group')}</option>
		<option value="piggy_bank">{t('admin.auditLog.piggyBank')}</option>
		<option value="recurring">{t('admin.auditLog.recurring')}</option>
		<option value="rule">{t('admin.auditLog.rule')}</option>
		<option value="currency">{t('admin.auditLog.currency')}</option>
	</Select>
</FilterBar>

<Card>
	<CardContent class="p-0">
		<div class="overflow-x-auto">
			<table class="w-full text-sm">
				<thead>
					<tr class="border-b bg-muted/50">
						<th class="text-left p-3 font-medium text-muted-foreground">{t('admin.auditLog.time')}</th>
						<th class="text-left p-3 font-medium text-muted-foreground">{t('admin.auditLog.user')}</th>
						<th class="text-left p-3 font-medium text-muted-foreground">{t('admin.auditLog.action')}</th>
						<th class="text-left p-3 font-medium text-muted-foreground">{t('admin.auditLog.entity')}</th>
						<th class="text-left p-3 font-medium text-muted-foreground">{t('admin.auditLog.changes')}</th>
					</tr>
				</thead>
				<tbody>
					{#each filtered() as log}
						<tr class="border-b hover:bg-muted/30">
							<td class="p-3 text-muted-foreground whitespace-nowrap">{formatDate(log.created_at)}</td>
							<td class="p-3 text-foreground">{log.user_email}</td>
							<td class="p-3">
								<span class="inline-flex items-center rounded-md bg-muted px-1.5 py-0.5 text-xs">{log.action}</span>
							</td>
							<td class="p-3 text-muted-foreground">{log.entity_type} ({log.entity_id})</td>
							<td class="p-3 text-foreground max-w-xs truncate">{log.changes}</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</CardContent>
</Card>

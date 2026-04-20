<script lang="ts">
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Plus, Copy, Trash2, AlertTriangle } from '@lucide/svelte';
	import { mockApiKeys } from '$lib/data/mock-api-keys.js';
	import { formatDate } from '$lib/utils/format.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;
</script>

<div class="flex flex-col gap-4">
	<div class="flex flex-wrap items-center justify-between gap-3">
		<h2 class="text-lg font-semibold text-foreground">{t('settings.apiKeys.title')}</h2>
		<Button size="sm">
			<Plus class="size-4" />
			{t('settings.apiKeys.createNew')}
		</Button>
	</div>

	<div class="flex items-start gap-3 rounded-lg border border-destructive/50 bg-destructive/5 p-4">
		<AlertTriangle class="size-5 shrink-0 text-destructive mt-0.5" />
		<div class="flex flex-col gap-1">
			<p class="text-sm font-semibold text-destructive">{t('settings.apiKeys.securityWarning')}</p>
			<p class="text-sm text-muted-foreground">{t('settings.apiKeys.securityWarningDesc')}</p>
		</div>
	</div>

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
						{#each mockApiKeys as key}
							<tr class="border-b last:border-b-0 hover:bg-muted/30">
								<td class="px-4 py-3 font-medium text-foreground">{key.title}</td>
								<td class="px-4 py-3 font-mono text-sm text-muted-foreground">{key.key}</td>
								<td class="px-4 py-3 text-muted-foreground">{formatDate(key.created_at)}</td>
								<td class="px-4 py-3 text-right">
									<div class="inline-flex items-center gap-2">
										<Button variant="ghost" size="sm">
											<Copy class="size-4" />
										</Button>
										<Button variant="ghost" size="sm" class="text-destructive hover:text-destructive">
											<Trash2 class="size-4" />
										</Button>
									</div>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</CardContent>
	</Card>
</div>

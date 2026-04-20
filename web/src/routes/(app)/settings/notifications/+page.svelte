<script lang="ts">
	import { Card, CardContent } from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Bell } from '@lucide/svelte';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	let notifications = $derived([
		{ id: 'n1', title: t('settings.notifications.newTransaction'), time: t('settings.notifications.time5min'), description: t('settings.notifications.newTransactionDesc'), isNew: true },
		{ id: 'n2', title: t('settings.notifications.budgetExceeded'), time: t('settings.notifications.time1hr'), description: t('settings.notifications.budgetExceededDesc'), isNew: true },
		{ id: 'n3', title: t('settings.notifications.billDue'), time: t('settings.notifications.time3hr'), description: t('settings.notifications.billDueDesc'), isNew: false },
		{ id: 'n4', title: t('settings.notifications.piggyTarget'), time: t('settings.notifications.time1day'), description: t('settings.notifications.piggyTargetDesc'), isNew: false },
		{ id: 'n5', title: t('settings.notifications.newMember'), time: t('settings.notifications.time2day'), description: t('settings.notifications.newMemberDesc'), isNew: false },
		{ id: 'n6', title: t('settings.notifications.weeklyReport'), time: t('settings.notifications.time3day'), description: t('settings.notifications.weeklyReportDesc'), isNew: false }
	]);
</script>

<div class="flex flex-col gap-4">
	<div class="flex flex-wrap items-center justify-between gap-3">
		<h2 class="text-lg font-semibold text-foreground">{t('settings.notifications.title')}</h2>
		<Button variant="outline" size="sm">
			<Bell class="size-4" />
			{t('settings.notifications.markAllRead')}
		</Button>
	</div>

	<div class="flex flex-col gap-3">
		{#each notifications as notif}
			<Card>
				<CardContent class="p-0">
					<div class="flex items-start justify-between px-5 py-4">
						<div class="flex flex-col gap-2 min-w-0">
							<div class="flex items-center gap-2">
								<p class="text-sm font-semibold text-foreground">{notif.title}</p>
								{#if notif.isNew}
									<span class="inline-flex items-center rounded-full bg-primary px-2 py-0.5 text-[11px] font-medium text-primary-foreground">{t('settings.notifications.new')}</span>
								{/if}
							</div>
							<p class="text-sm text-muted-foreground">{notif.description}</p>
						</div>
						<span class="shrink-0 text-xs text-muted-foreground">{notif.time}</span>
					</div>
				</CardContent>
			</Card>
		{/each}
	</div>
</div>

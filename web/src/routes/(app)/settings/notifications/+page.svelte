<script lang="ts">
	import { onMount } from 'svelte';
	import { Card, CardContent } from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Bell } from '@lucide/svelte';
	import { notificationService } from '$lib/services/index.js';
	import type { Notification } from '$lib/types/domain.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	let notifications = $state<Notification[]>([]);
	let isLoading = $state(true);
	let errorMsg = $state('');

	async function fetchNotifications() {
		try {
			notifications = await notificationService.list();
		} catch (e) {
			errorMsg = t('common.error');
			console.error(e);
		} finally {
			isLoading = false;
		}
	}

	async function handleMarkAllRead() {
		try {
			await notificationService.markAllRead();
			notifications = notifications.map((n) => ({ ...n, read: true }));
		} catch (e) {
			errorMsg = t('common.error');
			console.error(e);
		}
	}

	onMount(fetchNotifications);
</script>

<div class="flex flex-col gap-4">
	<div class="flex flex-wrap items-center justify-between gap-3">
		<h2 class="text-lg font-semibold text-foreground">{t('settings.notifications.title')}</h2>
		<Button variant="outline" size="sm" onclick={handleMarkAllRead}>
			<Bell class="size-4" />
			{t('settings.notifications.markAllRead')}
		</Button>
	</div>

	{#if isLoading}
		<div class="flex flex-col gap-3">
			{#each Array(3) as _}
				<Card>
					<CardContent class="px-5 py-4">
						<div class="flex flex-col gap-2 animate-pulse">
							<div class="h-4 w-1/3 rounded bg-muted"></div>
							<div class="h-3 w-2/3 rounded bg-muted"></div>
						</div>
					</CardContent>
				</Card>
			{/each}
		</div>
	{:else if errorMsg}
		<p class="py-8 text-center text-sm text-destructive">{errorMsg}</p>
	{:else if notifications.length === 0}
		<p class="py-8 text-center text-sm text-muted-foreground">{t('common.noData')}</p>
	{:else}
		<div class="flex flex-col gap-3">
			{#each notifications as notif}
				<Card>
					<CardContent class="p-0">
						<div class="flex items-start justify-between px-5 py-4">
							<div class="flex flex-col gap-2 min-w-0">
								<div class="flex items-center gap-2">
									<p class="text-sm font-semibold text-foreground">{notif.title}</p>
									{#if !notif.read}
										<span class="inline-flex items-center rounded-full bg-primary px-2 py-0.5 text-[11px] font-medium text-primary-foreground">{t('settings.notifications.new')}</span>
									{/if}
								</div>
								<p class="text-sm text-muted-foreground">{notif.message}</p>
							</div>
						</div>
					</CardContent>
				</Card>
			{/each}
		</div>
	{/if}
</div>

<script lang="ts">
	import { onMount } from 'svelte';
	import { authStore } from '$lib/stores/auth.svelte.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import PageHeader from '$lib/components/shared/PageHeader.svelte';
	import { Card, CardContent } from '$lib/components/ui/card/index.js';
	import { Shield, Mail, KeyRound } from '@lucide/svelte';
	const t = localeStore.t;

	let userEmail = $state('');

	onMount(() => {
		userEmail = authStore.user?.email ?? '';
	});
</script>

<div class="flex flex-col gap-6">
	<PageHeader title={t('settings.account.title')} description={t('settings.account.description')} icon={Shield} />

	<Card>
		<CardContent class="p-6">
			<h3 class="text-base font-semibold text-foreground mb-4 flex items-center gap-2">
				<Mail class="size-4" />
				{t('settings.account.changeEmail')}
			</h3>
			<p class="text-sm text-muted-foreground">{t('settings.account.email')}: <span class="font-medium text-foreground">{userEmail}</span></p>
			<p class="text-xs text-muted-foreground mt-2">{t('settings.account.emailHint')}</p>
		</CardContent>
	</Card>

	<Card>
		<CardContent class="p-6">
			<h3 class="text-base font-semibold text-foreground mb-4 flex items-center gap-2">
				<KeyRound class="size-4" />
				{t('settings.account.changePassword')}
			</h3>
			<p class="text-sm text-muted-foreground">{t('settings.account.passwordHint')}</p>
		</CardContent>
	</Card>
</div>

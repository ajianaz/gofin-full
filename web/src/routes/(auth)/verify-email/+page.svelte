<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Alert, AlertDescription } from '$lib/components/ui/alert/index.js';

	const t = localeStore.t;

	let token = $derived($page.url.searchParams.get('token') || '');
	let error = $state<string | null>(null);
	let success = $state(false);
	let submitting = $state(false);

	onMount(async () => {
		if (!token) {
			error = t('auth.verifyEmail.error');
			return;
		}

		submitting = true;
		try {
			const res = await fetch('/api/v1/auth/verify-email', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ token })
			});

			if (!res.ok) {
				const data = await res.json().catch(() => ({}));
				error = data.message || t('auth.verifyEmail.error');
				return;
			}

			success = true;
		} catch {
			error = t('auth.verifyEmail.error');
		} finally {
			submitting = false;
		}
	});
</script>

<Card>
	<CardHeader class="text-center px-8 pt-8">
		<CardTitle class="text-xl font-bold">{t('auth.verifyEmail.title')}</CardTitle>
	</CardHeader>

	<CardContent class="px-8 pb-8">
		{#if submitting}
			<div class="flex items-center justify-center py-8">
				<div class="animate-spin h-8 w-8 border-4 border-primary border-t-transparent rounded-full"></div>
			</div>
		{:else if success}
			<div class="text-center space-y-4 py-4">
				<p class="text-sm text-foreground">{t('auth.verifyEmail.success')}</p>
				<Button type="button" class="w-full" size="lg" onclick={() => goto('/login')}>
					{t('auth.resetPassword.login')}
				</Button>
			</div>
		{:else}
			<Alert variant="destructive">
				<AlertDescription>{error}</AlertDescription>
			</Alert>
			<div class="mt-4 text-center">
				<Button type="button" variant="outline" onclick={() => goto('/login')}>
					{t('common.back')}
				</Button>
			</div>
		{/if}
	</CardContent>
</Card>

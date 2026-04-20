<script lang="ts">
	import { goto } from '$app/navigation';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Alert, AlertDescription } from '$lib/components/ui/alert/index.js';

	const t = localeStore.t;

	let password = $state('');
	let confirmPassword = $state('');
	let error = $state<string | null>(null);
	let success = $state(false);

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		error = null;

		if (!password || !confirmPassword) {
			error = t('auth.resetPassword.errorRequired');
			return;
		}

		if (password.length < 8) {
			error = t('auth.resetPassword.errorPasswordLength');
			return;
		}

		if (password !== confirmPassword) {
			error = t('auth.resetPassword.errorPasswordMismatch');
			return;
		}

		success = true;
	}
</script>

<Card>
	<CardHeader class="text-center px-8 pt-8">
		<CardTitle class="text-xl font-bold">{t('auth.resetPassword.title')}</CardTitle>
		<p class="text-sm text-muted-foreground mt-1">{t('auth.resetPassword.subtitle')}</p>
	</CardHeader>

	<CardContent class="px-8 pb-8">
		<form onsubmit={handleSubmit} class="grid gap-4">
			{#if error}
				<Alert variant="destructive">
					<AlertDescription>{error}</AlertDescription>
				</Alert>
			{/if}

			{#if success}
				<div class="text-center space-y-4 py-4">
					<p class="text-sm text-foreground">{t('auth.resetPassword.success')}</p>
					<Button type="button" class="w-full" size="lg" onclick={() => goto('/login')}>
						{t('auth.resetPassword.login')}
					</Button>
				</div>
			{:else}
				<div class="grid gap-4">
					<div class="grid gap-2">
						<Label for="password">{t('auth.resetPassword.newPassword')}</Label>
						<Input
							id="password"
							type="password"
							placeholder={t('auth.resetPassword.newPasswordPlaceholder')}
							bind:value={password}
							required
							autocomplete="new-password"
						/>
					</div>

					<div class="grid gap-2">
						<Label for="confirm-password">{t('auth.resetPassword.confirmPassword')}</Label>
						<Input
							id="confirm-password"
							type="password"
							placeholder={t('auth.resetPassword.confirmPasswordPlaceholder')}
							bind:value={confirmPassword}
							required
							autocomplete="new-password"
						/>
					</div>
				</div>

				<Button type="submit" class="w-full" size="lg">
					{t('auth.resetPassword.submit')}
				</Button>
			{/if}
		</form>
	</CardContent>
</Card>

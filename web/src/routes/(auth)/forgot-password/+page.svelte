<script lang="ts">
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Alert, AlertDescription } from '$lib/components/ui/alert/index.js';

	const t = localeStore.t;

	let email = $state('');
	let error = $state<string | null>(null);
	let success = $state(false);

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		error = null;

		if (!email) {
			error = t('auth.forgotPassword.errorRequired');
			return;
		}

		success = true;
	}
</script>

<Card>
	<CardHeader class="text-center px-8 pt-8">
		<CardTitle class="text-xl font-bold">{t('auth.forgotPassword.title')}</CardTitle>
		<p class="text-sm text-muted-foreground mt-1">{t('auth.forgotPassword.subtitle')}</p>
	</CardHeader>

	<CardContent class="px-8 pb-8">
		<form onsubmit={handleSubmit} class="grid gap-4">
			{#if error}
				<Alert variant="destructive">
					<AlertDescription>{error}</AlertDescription>
				</Alert>
			{/if}

			{#if success}
				<div class="text-center space-y-2 py-4">
					<p class="text-sm text-foreground">
						{t('auth.forgotPassword.successMessage')} <strong>{email}</strong>.
					</p>
					<p class="text-xs text-muted-foreground">
						{t('auth.forgotPassword.checkInbox')}
					</p>
				</div>
			{:else}
				<div class="grid gap-2">
					<Label for="email">{t('auth.forgotPassword.email')}</Label>
					<Input
						id="email"
						type="email"
						placeholder={t('auth.forgotPassword.emailPlaceholder')}
						bind:value={email}
						required
						autocomplete="email"
					/>
				</div>
			{/if}

			<Button type="submit" class="w-full" size="lg" disabled={success}>
				{t('auth.forgotPassword.submit')}
			</Button>

			<p class="text-center text-sm text-muted-foreground">
				{t('auth.forgotPassword.rememberPassword')}
				<a href="/login" class="text-primary font-medium hover:underline">{t('auth.forgotPassword.login')}</a>
			</p>
		</form>
	</CardContent>
</Card>

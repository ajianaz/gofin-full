<script lang="ts">
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/auth.svelte.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import { Alert, AlertDescription } from '$lib/components/ui/alert/index.js';
	import type { ApiError } from '$lib/types/index.js';

	const t = localeStore.t;

	let email = $state('');
	let password = $state('');
	let confirmPassword = $state('');
	let error = $state<string | null>(null);

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		error = null;

		if (!email || !password || !confirmPassword) {
			error = t('auth.register.errorRequired');
			return;
		}

		if (password.length < 8) {
			error = t('auth.register.errorPasswordLength');
			return;
		}

		if (password !== confirmPassword) {
			error = t('auth.register.errorPasswordMismatch');
			return;
		}

		try {
			await authStore.register(email, password);
			goto('/dashboard');
		} catch (err) {
			const apiErr = err as ApiError;
			error = apiErr.detail || apiErr.message || t('auth.register.errorFailed');
		}
	}
</script>

<Card>
	<CardHeader class="text-center px-8 pt-8">
		<CardTitle class="text-xl font-bold">{t('auth.register.title')}</CardTitle>
		<p class="text-sm text-muted-foreground mt-1">{t('auth.register.subtitle')}</p>
	</CardHeader>

	<CardContent class="px-8 pb-8">
		<form onsubmit={handleSubmit} class="grid gap-4">
			{#if error}
				<Alert variant="destructive">
					<AlertDescription>{error}</AlertDescription>
				</Alert>
			{/if}

			<div class="grid gap-4">
				<div class="grid gap-2">
					<Label for="email">{t('auth.register.email')}</Label>
					<Input
						id="email"
						type="email"
						placeholder={t('auth.register.emailPlaceholder')}
						bind:value={email}
						required
						autocomplete="email"
					/>
				</div>

				<div class="grid gap-2">
					<Label for="password">{t('auth.register.password')}</Label>
					<Input
						id="password"
						type="password"
						placeholder={t('auth.register.passwordPlaceholder')}
						bind:value={password}
						required
						autocomplete="new-password"
					/>
				</div>

				<div class="grid gap-2">
					<Label for="confirm-password">{t('auth.register.confirmPassword')}</Label>
					<Input
						id="confirm-password"
						type="password"
						placeholder={t('auth.register.confirmPasswordPlaceholder')}
						bind:value={confirmPassword}
						required
						autocomplete="new-password"
					/>
				</div>
			</div>

			<Button type="submit" class="w-full" size="lg" disabled={authStore.isLoading}>
				{authStore.isLoading ? t('auth.register.submitting') : t('auth.register.submit')}
			</Button>

			<div class="flex items-center gap-3">
				<Separator class="flex-1" />
				<span class="text-xs text-muted-foreground">{t('common.or')}</span>
				<Separator class="flex-1" />
			</div>

			<Button variant="outline" class="w-full" size="lg" type="button">
				{t('auth.register.google')}
			</Button>

			<p class="text-center text-sm text-muted-foreground">
				{t('auth.register.hasAccount')}
				<a href="/login" class="text-primary font-medium hover:underline">{t('auth.register.login')}</a>
			</p>
		</form>
	</CardContent>
</Card>

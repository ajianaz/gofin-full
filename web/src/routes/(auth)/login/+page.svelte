<script lang="ts">
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/auth.svelte.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import { Alert, AlertDescription } from '$lib/components/ui/alert/index.js';
	import type { ApiError } from '$lib/types/index.js';

	const t = localeStore.t;

	let email = $state('');
	let password = $state('');
	let remember = $state(false);
	let error = $state<string | null>(null);

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		error = null;

		if (!email || !password) {
			error = t('auth.login.errorRequired');
			return;
		}

		try {
			await authStore.login(email, password);
			goto('/dashboard');
		} catch (err) {
			const apiErr = err as ApiError;
			error = apiErr.detail || apiErr.message || t('auth.login.errorFailed');
		}
	}
</script>

<Card>
	<CardHeader class="text-center px-8 pt-8">
		<CardTitle class="text-xl font-bold">{t('auth.login.title')}</CardTitle>
		<p class="text-sm text-muted-foreground mt-1">{t('auth.login.subtitle')}</p>
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
					<Label for="email">{t('auth.login.email')}</Label>
					<Input
						id="email"
						type="email"
						placeholder={t('auth.login.emailPlaceholder')}
						bind:value={email}
						required
						autocomplete="email"
					/>
				</div>

				<div class="grid gap-2">
					<Label for="password">{t('auth.login.password')}</Label>
					<Input
						id="password"
						type="password"
						placeholder={t('auth.login.passwordPlaceholder')}
						bind:value={password}
						required
						autocomplete="current-password"
					/>
				</div>

				<div class="flex items-center justify-between">
					<label class="flex items-center gap-2 text-sm">
						<Checkbox bind:checked={remember} />
						{t('auth.login.rememberMe')}
					</label>
					<a href="/forgot-password" class="text-sm text-primary font-medium hover:underline">
						{t('auth.login.forgotPassword')}
					</a>
				</div>
			</div>

			<Button type="submit" class="w-full" size="lg" disabled={authStore.isLoading}>
				{authStore.isLoading ? t('auth.login.submitting') : t('auth.login.submit')}
			</Button>

			<div class="flex items-center gap-3">
				<Separator class="flex-1" />
				<span class="text-xs text-muted-foreground">{t('common.or')}</span>
				<Separator class="flex-1" />
			</div>

			<Button variant="outline" class="w-full" size="lg" type="button">
				{t('auth.login.google')}
			</Button>

			<p class="text-center text-sm text-muted-foreground">
				{t('auth.login.noAccount')}
				<a href="/register" class="text-primary font-medium hover:underline">{t('auth.login.registerNow')}</a>
			</p>
		</form>
	</CardContent>
</Card>

<script lang="ts">
	import { onMount } from 'svelte';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import { authService } from '$lib/services/auth.js';
	import type { User } from '$lib/types/index.js';
	const t = localeStore.t;

	let user = $state<User | null>(null);
	let isLoading = $state(true);
	let isSaving = $state(false);
	let isChangingPassword = $state(false);
	let errorMsg = $state('');
	let successMsg = $state('');

	let name = $state('');
	let email = $state('');
	let currentPassword = $state('');
	let newPassword = $state('');
	let confirmPassword = $state('');

	let profileFormEl: HTMLFormElement | undefined = $state();
	let passwordFormEl: HTMLFormElement | undefined = $state();

	let initials = $derived(
		user
			? user.name
					.split(' ')
					.filter((w) => w.length > 0)
					.map((w) => w[0].toUpperCase())
					.slice(0, 2)
					.join('')
			: ''
	);

	async function fetchUser() {
		try {
			user = await authService.getMe();
			name = user.name;
			email = user.email;
		} catch (err: any) {
			errorMsg = err?.detail || err?.message || t('common.error');
		} finally {
			isLoading = false;
		}
	}

	async function handleSaveProfile(e: SubmitEvent) {
		e.preventDefault();
		isSaving = true;
		errorMsg = '';
		successMsg = '';
		try {
			user = await authService.updateProfile({ name });
			name = user.name;
			email = user.email;
			successMsg = t('settings.profile.saveSuccess');
		} catch (err: any) {
			errorMsg = err?.detail || err?.message || t('common.error');
		} finally {
			isSaving = false;
		}
	}

	function handleCancel() {
		errorMsg = '';
		successMsg = '';
		currentPassword = '';
		newPassword = '';
		confirmPassword = '';
		if (user) {
			name = user.name;
			email = user.email;
		}
	}

	async function handleChangePassword(e: SubmitEvent) {
		e.preventDefault();
		errorMsg = '';
		successMsg = '';
		if (newPassword !== confirmPassword) {
			errorMsg = t('settings.profile.passwordMismatch');
			return;
		}
		if (!currentPassword || !newPassword) {
			errorMsg = t('settings.profile.passwordRequired');
			return;
		}
		isChangingPassword = true;
		try {
			await authService.changePassword({ current_password: currentPassword, new_password: newPassword });
			currentPassword = '';
			newPassword = '';
			confirmPassword = '';
			successMsg = t('settings.profile.passwordChanged');
		} catch (err: any) {
			errorMsg = err?.detail || err?.message || t('common.error');
		} finally {
			isChangingPassword = false;
		}
	}

	onMount(fetchUser);
</script>

{#if isLoading}
	<div class="flex flex-col gap-4">
		<Card>
			<CardContent class="py-8">
				<div class="flex items-center justify-center text-muted-foreground">{t('common.loading')}</div>
			</CardContent>
		</Card>
	</div>
{:else}
	<div class="flex flex-col gap-4">
		<Card>
			<CardHeader>
				<CardTitle class="text-base">{t('settings.profile.profileInfo')}</CardTitle>
			</CardHeader>
			<CardContent>
				<div class="flex items-center gap-4">
					<div class="flex size-16 items-center justify-center rounded-full bg-primary text-lg font-semibold text-primary-foreground">
						{initials}
					</div>
					<div class="flex flex-col gap-1">
						<p class="text-base font-semibold text-foreground">{name}</p>
						<p class="text-sm text-muted-foreground">{email}</p>
					</div>
				</div>
			</CardContent>
		</Card>

		<Card>
			<CardHeader>
				<div class="flex items-center justify-between">
					<CardTitle class="text-base">{t('settings.profile.editProfile')}</CardTitle>
					<div class="flex gap-2">
						<Button size="sm" onclick={() => profileFormEl?.requestSubmit()} disabled={isSaving}>{isSaving ? t('common.saving') : t('settings.profile.saveChanges')}</Button>
						<Button size="sm" variant="outline" onclick={handleCancel}>{t('common.cancel')}</Button>
					</div>
				</div>
			</CardHeader>
			<CardContent>
				{#if errorMsg}
					<p class="text-destructive text-sm mb-4">{errorMsg}</p>
				{/if}
				{#if successMsg}
					<p class="text-green-600 text-sm mb-4">{successMsg}</p>
				{/if}
				<div class="flex flex-col gap-4">
					<form bind:this={profileFormEl} class="flex flex-col gap-4" onsubmit={handleSaveProfile}>
						<div class="grid gap-4 md:grid-cols-2">
							<div class="flex flex-col gap-2">
								<Label for="name">{t('settings.profile.fullName')}</Label>
								<Input id="name" bind:value={name} />
							</div>
							<div class="flex flex-col gap-2">
								<Label for="email">{t('settings.profile.email')}</Label>
								<Input id="email" value={email} disabled class="opacity-60" />
							</div>
						</div>
						<div class="flex gap-2 pt-4">
							<Button type="submit" class="flex-1" disabled={isSaving}>{isSaving ? t('common.saving') : t('settings.profile.saveChanges')}</Button>
							<Button type="button" variant="outline" class="flex-1" onclick={handleCancel}>{t('common.cancel')}</Button>
						</div>
					</form>

					<div class="border-t pt-4">
						<p class="text-sm font-medium text-foreground mb-3">{t('settings.profile.changePassword')}</p>
						<form bind:this={passwordFormEl} class="flex flex-col gap-4" onsubmit={handleChangePassword}>
							<div class="grid gap-4 md:grid-cols-2">
								<div class="flex flex-col gap-2">
									<Label for="current-pw">{t('settings.profile.currentPassword')}</Label>
									<Input id="current-pw" type="password" placeholder="••••••••" bind:value={currentPassword} />
								</div>
								<div class="flex flex-col gap-2">
									<Label for="new-pw">{t('settings.profile.newPassword')}</Label>
									<Input id="new-pw" type="password" placeholder={t('settings.profile.newPasswordPlaceholder')} bind:value={newPassword} />
								</div>
							</div>
							<div class="flex flex-col gap-2 max-w-md mt-4">
								<Label for="confirm-pw">{t('settings.profile.confirmPassword')}</Label>
								<Input id="confirm-pw" type="password" placeholder={t('settings.profile.confirmPasswordPlaceholder')} bind:value={confirmPassword} />
							</div>
							<div class="flex gap-2 pt-4">
								<Button type="submit" class="flex-1" disabled={isChangingPassword}>{isChangingPassword ? t('common.saving') : t('settings.profile.changePassword')}</Button>
							</div>
						</form>
					</div>
				</div>
			</CardContent>
		</Card>
	</div>
{/if}

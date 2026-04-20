<script lang="ts">
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	let name = $state('Anaz S. Aji');
	let currentPassword = $state('');
	let newPassword = $state('');
	let confirmPassword = $state('');
</script>

<div class="flex flex-col gap-4">
	<Card>
		<CardHeader>
			<CardTitle class="text-base">{t('settings.profile.profileInfo')}</CardTitle>
		</CardHeader>
		<CardContent>
			<div class="flex items-center gap-4">
				<div class="flex size-16 items-center justify-center rounded-full bg-primary text-lg font-semibold text-primary-foreground">
					AS
				</div>
				<div class="flex flex-col gap-1">
					<p class="text-base font-semibold text-foreground">{name}</p>
					<p class="text-sm text-muted-foreground">anaz@gofin.app</p>
				</div>
			</div>
		</CardContent>
	</Card>

	<Card>
		<CardHeader>
			<div class="flex items-center justify-between">
				<CardTitle class="text-base">{t('settings.profile.editProfile')}</CardTitle>
				<div class="flex gap-2">
					<Button size="sm">{t('settings.profile.saveChanges')}</Button>
					<Button size="sm" variant="outline">{t('common.cancel')}</Button>
				</div>
			</div>
		</CardHeader>
		<CardContent>
			<form class="flex flex-col gap-4" onsubmit={(e) => { e.preventDefault(); }}>
				<div class="grid gap-4 md:grid-cols-2">
					<div class="flex flex-col gap-2">
						<Label for="name">{t('settings.profile.fullName')}</Label>
						<Input id="name" bind:value={name} />
					</div>
					<div class="flex flex-col gap-2">
						<Label for="email">{t('settings.profile.email')}</Label>
						<Input id="email" value="anaz@gofin.app" disabled class="opacity-60" />
					</div>
				</div>
				<div class="border-t pt-4">
					<p class="text-sm font-medium text-foreground mb-3">{t('settings.profile.changePassword')}</p>
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
				</div>
				<div class="flex gap-2 pt-4">
					<Button type="submit" class="flex-1">{t('settings.profile.saveChanges')}</Button>
					<Button type="button" variant="outline" class="flex-1">{t('common.cancel')}</Button>
				</div>
			</form>
		</CardContent>
	</Card>
</div>

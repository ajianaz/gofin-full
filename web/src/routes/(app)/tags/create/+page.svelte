<script lang="ts">
	import { goto } from '$app/navigation';
	import { BackButton, FormCard } from '$lib/components/shared/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	let tag = $state('');
	let date = $state(new Date().toISOString().split('T')[0]);
	let description = $state('');
</script>

<BackButton href="/tags" />

<FormCard title={t('tags.create.title')} description={t('tags.create.description')}>
	<form class="grid gap-4" onsubmit={(e) => { e.preventDefault(); goto('/tags'); }}>
		<div class="grid gap-2">
			<Label for="tag">{t('tags.create.tag')}</Label>
			<Input id="tag" placeholder={t('tags.create.tagPlaceholder')} bind:value={tag} required />
		</div>

		<div class="grid gap-2">
			<Label for="date">{t('tags.create.date')}</Label>
			<Input id="date" type="date" bind:value={date} required />
		</div>

		<div class="grid gap-2">
			<Label for="desc">{t('tags.create.descriptionField')}</Label>
			<Input id="desc" placeholder={t('tags.create.descriptionPlaceholder')} bind:value={description} />
		</div>

		<div class="flex gap-3 pt-2">
			<Button type="submit">{t('common.save')}</Button>
			<Button type="button" variant="outline" onclick={() => goto('/tags')}>{t('common.cancel')}</Button>
		</div>
	</form>
</FormCard>

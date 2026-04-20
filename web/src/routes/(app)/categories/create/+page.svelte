<script lang="ts">
	import { goto } from '$app/navigation';
	import { BackButton, FormCard } from '$lib/components/shared/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Select } from '$lib/components/ui/select/index.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	let name = $state('');
	let type = $state('expense');
	let parent = $state('');
</script>

<BackButton href="/categories" />

<FormCard title={t('categories.create.title')} description={t('categories.create.description')}>
	<form class="grid gap-4" onsubmit={(e) => { e.preventDefault(); goto('/categories'); }}>
		<div class="grid gap-2">
			<Label for="name">{t('categories.create.name')}</Label>
			<Input id="name" placeholder={t('categories.create.namePlaceholder')} bind:value={name} required />
		</div>

		<div class="grid gap-2">
			<Label for="type">{t('categories.create.type')}</Label>
			<Select bind:value={type} id="type">
				<option value="expense">{t('categories.create.expense')}</option>
				<option value="income">{t('categories.create.income')}</option>
				<option value="transfer">{t('categories.create.transfer')}</option>
			</Select>
		</div>

		<div class="flex gap-3 pt-2">
			<Button type="submit">{t('common.save')}</Button>
			<Button type="button" variant="outline" onclick={() => goto('/categories')}>{t('common.cancel')}</Button>
		</div>
	</form>
</FormCard>

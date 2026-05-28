<script lang="ts">
	import { goto } from '$app/navigation';
	import { BackButton, FormCard } from '$lib/components/shared/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Select, SelectTrigger, SelectContent, SelectItem, SelectValue } from '$lib/components/ui/select/index.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import { categoryService } from '$lib/services/index.js';
	const t = localeStore.t;

	let isLoading = $state(false);
	let errorMsg = $state('');
	let name = $state('');
	let type = $state('expense');
	let parent = $state('');
</script>

<BackButton href="/categories" />

<FormCard title={t('categories.create.title')} description={t('categories.create.description')}>
	<form class="grid gap-4" onsubmit={async (e) => { e.preventDefault(); isLoading = true; errorMsg = ''; try { await categoryService.create({ name }); goto('/categories'); } catch (err: any) { errorMsg = err?.detail || err?.message || t('common.errorSave'); } finally { isLoading = false; } }}>
		<div class="grid gap-2">
			<Label for="name">{t('categories.create.name')}</Label>
			<Input id="name" placeholder={t('categories.create.namePlaceholder')} bind:value={name} required />
		</div>

		<div class="grid gap-2">
			<Label for="type">{t('categories.create.type')}</Label>
		<Select bind:value={type} items={[
			{ value: 'expense', label: t('categories.create.expense') },
			{ value: 'income', label: t('categories.create.income') },
			{ value: 'transfer', label: t('categories.create.transfer') },
		]}>
			<SelectTrigger class="w-full">
				<SelectValue />
			</SelectTrigger>
			<SelectContent>
				<SelectItem value="expense" label={t('categories.create.expense')}>{t('categories.create.expense')}</SelectItem>
				<SelectItem value="income" label={t('categories.create.income')}>{t('categories.create.income')}</SelectItem>
				<SelectItem value="transfer" label={t('categories.create.transfer')}>{t('categories.create.transfer')}</SelectItem>
			</SelectContent>
		</Select>
		</div>

		{#if errorMsg}
			<p class="text-destructive text-sm">{errorMsg}</p>
		{/if}

		<div class="flex flex-col sm:flex-row gap-3 pt-2">
			<Button type="submit" disabled={isLoading}>{isLoading ? t('common.saving') : t('common.save')}</Button>
			<Button type="button" variant="outline" onclick={() => goto('/categories')}>{t('common.cancel')}</Button>
		</div>
	</form>
</FormCard>

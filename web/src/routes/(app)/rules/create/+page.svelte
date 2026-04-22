<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Textarea } from '$lib/components/ui/textarea/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { ruleService } from '$lib/services/index.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	let isLoading = $state(false);
	let errorMsg = $state('');

	let title = $state('');
	let order = $state('1');
	let description = $state('');
	let active = $state(true);
	let stopProcessing = $state(false);

	async function handleSubmit(e: Event) {
		e.preventDefault();
		isLoading = true;
		errorMsg = '';
		try {
			await ruleService.createGroup({ title });
			goto('/rules');
		} catch (err: any) {
			errorMsg = err.detail || err.message || 'Gagal menyimpan';
		} finally {
			isLoading = false;
		}
	}
</script>

<div class="flex flex-col gap-4">
	<Card>
		<CardHeader>
			<div class="flex items-center justify-between">
				<CardTitle class="text-base">{t('rules.createGroup.title')}</CardTitle>
				<div class="flex gap-2">
					<Button size="sm" onclick={() => goto('/rules')}>{t('common.save')}</Button>
					<Button size="sm" variant="outline" onclick={() => goto('/rules')}>{t('common.cancel')}</Button>
				</div>
			</div>
		</CardHeader>
		<CardContent>
			<form class="flex flex-col gap-6" onsubmit={handleSubmit}>
				<div class="grid gap-6 md:grid-cols-2">
					<div class="flex flex-col gap-2">
						<Label for="title">{t('rules.createGroup.name')}</Label>
						<Input id="title" placeholder={t('rules.createGroup.namePlaceholder')} bind:value={title} required />
					</div>
					<div class="flex flex-col gap-2">
						<Label for="order">{t('rules.createGroup.order')}</Label>
						<Input id="order" type="number" placeholder="1" bind:value={order} />
					</div>
				</div>
				<div class="flex items-center gap-6">
					<div class="flex items-center gap-2">
						<Checkbox id="active" bind:checked={active} />
						<Label for="active">{t('rules.createGroup.activate')}</Label>
					</div>
					<div class="flex items-center gap-2">
						<Checkbox id="stop" bind:checked={stopProcessing} />
						<Label for="stop">{t('rules.createGroup.stopOnMatch')}</Label>
					</div>
				</div>
				<div class="flex flex-col gap-2">
					<Label for="desc">{t('rules.createGroup.description')}</Label>
					<Textarea id="desc" placeholder={t('rules.createGroup.descriptionPlaceholder')} bind:value={description} rows={3} />
				</div>
				{#if errorMsg}
					<p class="text-sm text-destructive">{errorMsg}</p>
				{/if}
				<div class="flex gap-2 pt-2">
					<Button type="submit" class="flex-1" disabled={isLoading}>{isLoading ? '...' : t('common.save')}</Button>
					<Button type="button" variant="outline" class="flex-1" onclick={() => goto('/rules')}>{t('common.cancel')}</Button>
				</div>
			</form>
		</CardContent>
	</Card>
</div>

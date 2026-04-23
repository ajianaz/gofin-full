<script lang="ts">
	import { onMount } from 'svelte';
	import { PageHeader, FormCard } from '$lib/components/shared/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Select } from '$lib/components/ui/select/index.js';
	import { exportService, walletService } from '$lib/services/index.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	let format = $state('csv');
	let startDate = $state('');
	let endDate = $state('');
	let walletId = $state('');
	let wallets = $state<{ id: string; name: string }[]>([]);
	let isLoading = $state(true);
	let isExporting = $state(false);
	let error = $state('');

	async function loadWallets() {
		isLoading = true;
		error = '';
		try {
			const data = await walletService.list();
			wallets = data.map((w) => ({ id: w.id, name: w.name }));
		} catch (e: any) {
			error = e.message || 'Failed to load wallets';
		} finally {
			isLoading = false;
		}
	}

	async function handleSubmit(e: Event) {
		e.preventDefault();
		isExporting = true;
		error = '';
		try {
			if (format === 'csv') {
				await exportService.downloadCSV(startDate || undefined, endDate || undefined, walletId || undefined);
			} else {
				await exportService.downloadOFX(startDate || undefined, endDate || undefined, walletId || undefined);
			}
		} catch (e: any) {
			error = e.message || 'Export failed';
		} finally {
			isExporting = false;
		}
	}

	onMount(() => {
		loadWallets();
	});
</script>

<PageHeader title={t('export.title')} description={t('export.description')} />

<FormCard title={t('export.exportData')}>
	{#if error}
		<p class="mb-4 text-sm text-destructive">{error}</p>
	{/if}

	<form class="grid gap-4" onsubmit={handleSubmit}>
		<div class="grid gap-2">
			<Label for="format">{t('export.format')}</Label>
			<Select bind:value={format} id="format">
				<option value="csv">CSV</option>
				<option value="ofx">OFX</option>
			</Select>
		</div>

		<div class="grid grid-cols-2 gap-4">
			<div class="grid gap-2">
				<Label for="start">{t('export.startDate')}</Label>
				<Input id="start" type="date" bind:value={startDate} />
			</div>
			<div class="grid gap-2">
				<Label for="end">{t('export.endDate')}</Label>
				<Input id="end" type="date" bind:value={endDate} />
			</div>
		</div>

		<div class="grid gap-2">
			<Label for="wallet">{t('export.wallet')}</Label>
			<Select bind:value={walletId} id="wallet" disabled={isLoading}>
				<option value="">{t('export.allWallets')}</option>
				{#each wallets as w}
					<option value={w.id}>{w.name}</option>
				{/each}
			</Select>
			{#if isLoading}
				<p class="text-sm text-muted-foreground">{t('common.loading')}</p>
			{/if}
		</div>

		<Button type="submit" disabled={isExporting || isLoading}>
			{isExporting ? t('export.exporting') : t('export.export')}
		</Button>
	</form>
</FormCard>

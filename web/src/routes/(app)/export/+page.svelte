<script lang="ts">
	import { PageHeader, FormCard } from '$lib/components/shared/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Select } from '$lib/components/ui/select/index.js';
	import { mockWallets } from '$lib/data/mock-wallets.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	let format = $state('csv');
	let startDate = $state('');
	let endDate = $state('');
	let walletId = $state('');
</script>

<PageHeader title={t('export.title')} description={t('export.description')} />

<FormCard title={t('export.exportData')}>
	<form class="grid gap-4" onsubmit={(e) => { e.preventDefault(); }}>
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
			<Select bind:value={walletId} id="wallet">
				<option value="">{t('export.allWallets')}</option>
				{#each mockWallets as w}
					<option value={w.id}>{w.name}</option>
				{/each}
			</Select>
		</div>

		<Button type="submit">{t('export.export')}</Button>
	</form>
</FormCard>

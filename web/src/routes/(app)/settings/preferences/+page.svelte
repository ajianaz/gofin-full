<script lang="ts">
	import { onMount } from 'svelte';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { ChevronDown } from '@lucide/svelte';
	import { preferenceService } from '$lib/services/index.js';
	import type { PreferenceItem } from '$lib/types/domain.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import { Select, SelectTrigger, SelectContent, SelectItem } from '$lib/components/ui/select/index.js';
	const t = localeStore.t;

	let preferences = $state<PreferenceItem[]>([]);
	let isLoading = $state(true);
	let errorMsg = $state('');

	const prefConfig: Record<string, { type: 'boolean' | 'select' | 'text' | 'number'; options?: string[] }> = {
		language: { type: 'select', options: ['id', 'en'] },
		currency: { type: 'select', options: ['IDR', 'USD', 'EUR'] },
		date_format: { type: 'select', options: ['DD MMM YYYY', 'YYYY-MM-DD', 'MM/DD/YYYY'] },
		group_style: { type: 'select', options: ['default', 'compact'] },
		budget_indicator: { type: 'boolean' },
		show_news: { type: 'boolean' },
		fiscal_year_start: { type: 'select', options: ['01', '04', '07', '10'] },
		transaction_count_per_page: { type: 'select', options: ['10', '15', '25', '50'] },
		two_factor_enabled: { type: 'boolean' },
		email_digest: { type: 'select', options: ['daily', 'weekly', 'monthly'] }
	};

	function getConfig(name: string) {
		return prefConfig[name] ?? { type: 'text' as const };
	}

	function parseValue(pref: PreferenceItem) {
		const config = getConfig(pref.name);
		if (config.type === 'boolean') {
			return pref.data === 'true';
		}
		if (config.type === 'number') {
			return Number(pref.data);
		}
		return pref.data;
	}

	async function handleCheckboxChange(pref: PreferenceItem) {
		const newVal = pref.data === 'true' ? 'false' : 'true';
		pref.data = newVal;
		try {
			await preferenceService.set(pref.name, newVal);
		} catch (e) {
			pref.data = newVal === 'true' ? 'false' : 'true';
			console.error(e);
		}
	}

	async function handleSelectChange(pref: PreferenceItem, value: string) {
		const oldVal = pref.data;
		pref.data = value;
		try {
			await preferenceService.set(pref.name, value);
		} catch (e) {
			pref.data = oldVal;
			console.error(e);
		}
	}

	async function handleInputChange(pref: PreferenceItem, value: string) {
		const oldVal = pref.data;
		pref.data = value;
		try {
			await preferenceService.set(pref.name, value);
		} catch (e) {
			pref.data = oldVal;
			console.error(e);
		}
	}

	onMount(async () => {
		try {
			preferences = await preferenceService.list();
		} catch (e) {
			errorMsg = t('common.error');
			console.error(e);
		} finally {
			isLoading = false;
		}
	});
</script>

<div class="flex flex-col gap-4">
	<h2 class="text-lg font-semibold text-foreground">{t('settings.preferences.title')}</h2>

	{#if isLoading}
<Card>
		<CardContent class="py-8">
			{#each Array(4) as _}
				<div class="flex flex-col gap-2 mb-4">
					<Skeleton class="h-4 w-24" />
					<Skeleton class="h-9 w-full rounded-md" />
				</div>
			{/each}
		</CardContent>
	</Card>
	{:else if errorMsg}
		<Card>
			<CardContent class="px-5 py-8 text-center text-sm text-destructive">
				{errorMsg}
			</CardContent>
		</Card>
	{:else}
		<Card>
			<CardHeader>
				<CardTitle class="text-base">{t('settings.preferences.appSettings')}</CardTitle>
			</CardHeader>
			<CardContent class="p-0">
				{#each preferences as pref (pref.id)}
					{@const config = getConfig(pref.name)}
					{@const value = parseValue(pref)}
					<div class="flex items-center justify-between px-5 py-3.5 border-b last:border-b-0">
						<p class="text-sm text-foreground">{t(`settings.preferences.${pref.name}`)}</p>
						{#if config.type === 'boolean'}
							<Checkbox checked={value as boolean} onchange={() => handleCheckboxChange(pref)} />
						{:else if config.type === 'select' && config.options}
							<div class="relative">
								<Select
									value={String(value)}
								onValueChange={(v) => handleSelectChange(pref, v)}
								>
								<SelectTrigger class="h-9 w-40">
								</SelectTrigger>
								<SelectContent>
									{#each config.options as opt (opt)}
										<SelectItem value={opt}>{opt}</SelectItem>
									{/each}
								</SelectContent>
								</Select>
								<ChevronDown class="pointer-events-none absolute right-2.5 top-1/2 size-3.5 -translate-y-1/2 text-muted-foreground" />
							</div>
						{:else if config.type === 'number'}
							<Input
								type="number"
								value={String(value)}
								class="h-9 w-32 text-sm"
								onchange={(e) => handleInputChange(pref, (e.target as HTMLInputElement).value)}
							/>
						{:else}
							<Input
								value={String(value)}
								class="h-9 w-40 text-sm"
								onchange={(e) => handleInputChange(pref, (e.target as HTMLInputElement).value)}
							/>
						{/if}
					</div>
				{/each}
			</CardContent>
		</Card>
	{/if}
</div>

<script lang="ts">
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { ChevronDown } from '@lucide/svelte';
	import { mockPreferences } from '$lib/data/mock-preferences.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;
</script>

<div class="flex flex-col gap-4">
	<h2 class="text-lg font-semibold text-foreground">{t('settings.preferences.title')}</h2>

	<Card>
		<CardHeader>
			<CardTitle class="text-base">{t('settings.preferences.appSettings')}</CardTitle>
		</CardHeader>
		<CardContent class="p-0">
			{#each mockPreferences as pref}
				<div class="flex items-center justify-between px-5 py-3.5 border-b last:border-b-0">
					<p class="text-sm text-foreground">{t(`settings.preferences.${pref.name}`)}</p>
					{#if pref.type === 'boolean'}
						<Checkbox checked={!!pref.value} />
					{:else if pref.type === 'select' && pref.options}
						<div class="relative">
							<select
								value={String(pref.value)}
								class="cn-input h-9 w-40 appearance-none bg-background pr-8 text-sm"
							>
								{#each pref.options as opt}
									<option value={opt}>{opt}</option>
								{/each}
							</select>
							<ChevronDown class="pointer-events-none absolute right-2.5 top-1/2 size-3.5 -translate-y-1/2 text-muted-foreground" />
						</div>
					{:else if pref.type === 'number'}
						<Input type="number" value={String(pref.value)} class="h-9 w-32 text-sm" />
					{:else}
						<Input value={String(pref.value)} class="h-9 w-40 text-sm" />
					{/if}
				</div>
			{/each}
		</CardContent>
	</Card>
</div>

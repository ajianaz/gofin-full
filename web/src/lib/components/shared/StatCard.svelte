<script lang="ts">
	import { Card, CardContent } from '$lib/components/ui/card/index.js';
	import type { Component } from 'svelte';

	let {
		title,
		value,
		icon,
		trend
	}: {
		title: string;
		value: string;
		icon?: Component;
		trend?: { value: string; positive: boolean };
	} = $props();
</script>

<Card>
	<CardContent class="p-4">
		<div class="flex items-center justify-between">
			<p class="text-sm text-muted-foreground">{title}</p>
			{#if icon}
				{@const Icon = icon}
				<div class="flex size-9 items-center justify-center rounded-lg bg-primary/10">
					<Icon class="size-4 text-primary" />
				</div>
			{/if}
		</div>
		<div class="mt-2">
			<p class="text-2xl font-bold text-foreground">{value}</p>
			{#if trend}
				<p class="text-xs mt-1 {trend.positive ? 'text-green-600' : 'text-red-600'}">
					{trend.value}
				</p>
			{/if}
		</div>
	</CardContent>
</Card>

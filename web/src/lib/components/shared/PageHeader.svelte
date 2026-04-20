<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { Component } from 'svelte';

	let {
		title,
		description,
		icon,
		actions
	}: {
		title: string;
		description?: string;
		icon?: Component;
		actions?: Snippet;
	} = $props();
</script>

<div class="flex items-center justify-between mb-6">
	<div class="flex items-center gap-3">
		{#if icon}
			{@const Icon = icon}
			<div class="flex size-10 items-center justify-center rounded-lg bg-primary/10">
				<Icon class="size-5 text-primary" />
			</div>
		{/if}
		<div>
			<h1 class="text-2xl font-bold text-foreground">{title}</h1>
			{#if description}
				<p class="text-sm text-muted-foreground mt-0.5">{description}</p>
			{/if}
		</div>
	</div>
	{#if actions}
		<div class="flex items-center gap-2">
			{@render actions()}
		</div>
	{/if}
</div>

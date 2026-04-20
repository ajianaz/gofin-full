<script lang="ts">
	import type { Snippet } from 'svelte';
	import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '$lib/components/ui/card/index.js';

	let {
		title,
		description,
		children
	}: {
		title: string;
		description?: string;
		children: Snippet;
	} = $props();
</script>

<Card class="max-w-2xl">
	<CardHeader>
		<CardTitle>{title}</CardTitle>
		{#if description}
			<CardDescription>{description}</CardDescription>
		{/if}
	</CardHeader>
	<CardContent>
		{@render children()}
	</CardContent>
</Card>

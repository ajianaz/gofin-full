<script lang="ts">
	import { cn } from "$lib/utils.js";
	import type { Snippet } from "svelte";
	import { Dialog as DialogPrimitive } from "bits-ui";

	let {
		class: className,
		open = $bindable(false),
		children,
		...restProps
	}: DialogPrimitive.RootProps & {
		class?: string;
		children: Snippet;
	} = $props();
</script>

<DialogPrimitive.Root bind:open {...restProps}>
	{@render children?.()}
</DialogPrimitive.Root>

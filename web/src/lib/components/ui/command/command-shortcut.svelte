<script lang="ts">
	import { cn, type WithElementRef } from "$lib/utils.js";
	import type { Snippet } from "svelte";
	import type { HTMLAttributes } from "svelte/elements";

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLSpanElement>> & {
		children?: Snippet;
	} = $props();
</script>

<span
	bind:this={ref}
	data-slot="command-shortcut"
	class={cn("cn-command-shortcut ml-auto text-xs tracking-widest text-muted-foreground", className)}
	{...restProps}
>
	{@render children?.()}
</span>

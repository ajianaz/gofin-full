<script lang="ts">
	import type { HTMLAttributes } from "svelte/elements";
	import { cn, type WithElementRef, type WithoutChildren } from "$lib/utils.js";
	import { MoreHorizontal } from "@lucide/svelte";

	let {
		ref = $bindable(null),
		class: className,
		...restProps
	}: WithoutChildren<WithElementRef<HTMLAttributes<HTMLSpanElement>>> = $props();
</script>

<span
	bind:this={ref}
	aria-hidden="true"
	data-slot="pagination-ellipsis"
	class={cn("cn-pagination-ellipsis flex items-center justify-center", className)}
	{...restProps}
>
	<MoreHorizontal class="size-4" />
	<span class="sr-only">More pages</span>
</span>

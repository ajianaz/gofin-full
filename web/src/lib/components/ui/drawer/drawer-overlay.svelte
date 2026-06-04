<script lang="ts">
	import { Drawer as DrawerPrimitive } from "vaul-svelte";
	import { cn } from "$lib/utils.js";
	import type { HTMLAttributes } from "svelte/elements";

	let {
		class: className,
		...restProps
	}: HTMLAttributes<HTMLDivElement> = $props();
</script>

<DrawerPrimitive.Overlay
	data-slot="drawer-overlay"
	class={cn("cn-drawer-overlay fixed inset-0 z-50", className)}
	{...restProps}
/>

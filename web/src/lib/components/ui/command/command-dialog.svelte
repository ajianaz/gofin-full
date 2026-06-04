<script lang="ts">
	import { cn } from "$lib/utils.js";
	import type { Snippet } from "svelte";
	import { Command as CommandPrimitive } from "bits-ui";

	let {
		ref = $bindable(null),
		class: className,
		open = $bindable(false),
		children,
		...restProps
	}: CommandPrimitive.DialogProps & {
		children: Snippet;
	} = $props();
</script>

<CommandPrimitive.Dialog bind:open>
	{@render children?.()}
</CommandPrimitive.Dialog>

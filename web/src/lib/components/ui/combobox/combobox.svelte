<script lang="ts">
	import { cn } from "$lib/utils.js";
	import type { ComponentProps } from "svelte";
	import { Command } from "$lib/components/ui/command/index.js";

	let {
		ref = $bindable(null),
		class: className,
		...restProps
	}: ComponentProps<typeof Command> = $props();
</script>

<Command
	bind:ref
	data-slot="combobox"
	class={cn("cn-combobox overflow-open", className)}
	{...restProps}
/>

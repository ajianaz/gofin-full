<script lang="ts">
	import { cn } from "$lib/utils.js";
	import { Command as CommandPrimitive } from "bits-ui";
	import { Search } from "@lucide/svelte";

	let {
		ref = $bindable(null),
		class: className,
		...restProps
	}: CommandPrimitive.InputProps = $props();
</script>

<div class="cn-command-input-wrapper flex items-center border-b px-3" data-slot="command-input-wrapper">
	<Search class="cn-command-input-search mr-2 size-4 shrink-0 opacity-50" />
	<CommandPrimitive.Input
		bind:ref
		data-slot="command-input"
		class={cn(
			"cn-command-input flex h-10 w-full bg-transparent py-3 text-sm outline-none placeholder:text-muted-foreground disabled:cursor-not-allowed disabled:opacity-50",
			className
		)}
		{...restProps}
	/>
</div>

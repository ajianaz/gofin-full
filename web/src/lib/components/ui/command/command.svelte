<script lang="ts">
	import { cn } from "$lib/utils.js";
	import { Command as CommandPrimitive } from "bits-ui";

	let {
		ref = $bindable(null),
		class: className,
		...restProps
	}: CommandPrimitive.RootProps = $props();
</script>

<CommandPrimitive.Root
	bind:ref
	data-slot="command"
	class={cn("cn-command flex h-full w-full flex-col overflow-hidden", className)}
	{...restProps}
/>

<script lang="ts">
	import { cn } from "$lib/utils.js";
	import { Command as CommandPrimitive } from "bits-ui";
	import { Loader2 } from "@lucide/svelte";

	let {
		ref = $bindable(null),
		class: className,
		...restProps
	}: CommandPrimitive.LoadingProps = $props();
</script>

<CommandPrimitive.Loading
	bind:ref
	data-slot="command-loading"
	class={cn("cn-command-loading flex items-center justify-center py-6", className)}
	{...restProps}
>
	<Loader2 class="size-4 animate-spin" />
</CommandPrimitive.Loading>

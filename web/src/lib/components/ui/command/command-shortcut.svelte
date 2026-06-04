<script lang="ts">
	import { cn } from "$lib/utils.js";

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: Record<string, unknown> = $props();
</script>

<span
	bind:this={ref}
	data-slot="command-shortcut"
	class={cn("cn-command-shortcut ml-auto text-xs tracking-widest text-muted-foreground", className as string)}
	{...restProps}
>
	{@render children?.()}
</span>

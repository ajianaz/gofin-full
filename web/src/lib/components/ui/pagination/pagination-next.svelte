<script lang="ts">
	import { Pagination as PaginationPrimitive } from "bits-ui";
	import { cn } from "$lib/utils.js";
	import { buttonVariants } from "$lib/components/ui/button/index.js";
	import { ChevronRight } from "@lucide/svelte";

	let {
		ref = $bindable(null),
		class: className,
		...restProps
	}: PaginationPrimitive.NextButtonProps = $props();
</script>

<PaginationPrimitive.NextButton
	bind:ref
	aria-label="Go to next page"
	class={cn(
		buttonVariants({ variant: "ghost", size: "default" }),
		"cn-pagination-next",
		className
	)}
	{...restProps}
>
	<span class="cn-pagination-next-text hidden sm:block">Next</span>
	<ChevronRight data-icon="inline-end" class="size-4" />
</PaginationPrimitive.NextButton>

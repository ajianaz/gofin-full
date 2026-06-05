<script lang="ts">
	import { Drawer as DrawerPrimitive } from "vaul-svelte";
	import type { HTMLAttributes } from "svelte/elements";

	let { ...restProps }: HTMLAttributes<HTMLDivElement> = $props();
</script>

<DrawerPrimitive.Portal {...restProps} />

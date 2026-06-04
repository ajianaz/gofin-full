<script lang="ts">
	import { Drawer as DrawerPrimitive } from "vaul-svelte";
	import type { HTMLButtonAttributes } from "svelte/elements";

	let { ...restProps }: HTMLButtonAttributes = $props();
</script>

<DrawerPrimitive.Trigger data-slot="drawer-trigger" {...restProps} />

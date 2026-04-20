<script lang="ts">
	import { Select as SelectPrimitive } from "bits-ui";

	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	type AnyRecord = Record<string, any>;

	type SelectRootSingleProps = Omit<SelectPrimitive.RootProps, "value" | "onValueChange" | "type"> & {
		value?: string;
		onValueChange?: (value: string) => void;
		type?: "single";
	} & AnyRecord;

	let {
		open = $bindable(false),
		value = $bindable(),
		type = "single",
		...restProps
	}: SelectRootSingleProps = $props();
</script>

<SelectPrimitive.Root bind:open bind:value={value as never} type={type} {...restProps} />

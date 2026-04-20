<script lang="ts">
	import { formatAmount } from '$lib/utils/format.js';

	let {
		amount,
		class: className = ''
	}: {
		amount: string;
		class?: string;
	} = $props();

	let formatted = $derived(formatAmount(amount));
</script>

<span class="font-semibold {formatted.color} {className}">
	{formatted.text}
</span>

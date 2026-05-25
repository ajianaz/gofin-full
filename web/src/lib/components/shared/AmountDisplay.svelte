<script lang="ts">
	import { formatAmount } from '$lib/utils/format.js';

	let {
		amount,
		symbol,
		class: className = ''
	}: {
		amount: string;
		symbol?: string;
		class?: string;
	} = $props();

	let formatted = $derived(formatAmount(amount, symbol));
</script>

<span class="font-semibold {formatted.color} {className}">
	{formatted.text}
</span>

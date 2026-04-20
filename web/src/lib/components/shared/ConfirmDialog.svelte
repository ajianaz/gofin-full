<script lang="ts">
	import {
		Dialog,
		DialogContent,
		DialogDescription,
		DialogFooter,
		DialogHeader,
		DialogTitle
	} from '$lib/components/ui/dialog/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	let {
		open = $bindable(false),
		title,
		description,
		confirmLabel,
		cancelLabel,
		onConfirm
	}: {
		open?: boolean;
		title: string;
		description: string;
		confirmLabel?: string;
		cancelLabel?: string;
		onConfirm: () => void;
	} = $props();
</script>

<Dialog bind:open>
	<DialogContent>
		<DialogHeader>
			<DialogTitle>{title}</DialogTitle>
			<DialogDescription>{description}</DialogDescription>
		</DialogHeader>
		<DialogFooter>
			<Button variant="outline" onclick={() => (open = false)}>{cancelLabel ?? t('common.cancel')}</Button>
			<Button variant="destructive" onclick={() => { onConfirm(); open = false; }}>{confirmLabel ?? t('common.delete')}</Button>
		</DialogFooter>
	</DialogContent>
</Dialog>

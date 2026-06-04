<script lang="ts">
	import { page } from '$app/state';
	import { Button } from '$lib/components/ui/button';

	let { status, message } = $props();

	// Use props (provided by SvelteKit error loader) as primary source,
	// fallback to page state. In +error.svelte, SvelteKit passes status
	// and message as props automatically.
	const errorStatus = $derived(status || page.status);
	const errorMessage = $derived(message || page.error?.message || 'An unexpected error occurred');

	const titles: Record<number, string> = {
		400: 'Bad Request',
		401: 'Unauthorized',
		403: 'Forbidden',
		404: 'Page Not Found',
		405: 'Method Not Allowed',
		429: 'Too Many Requests',
		500: 'Internal Server Error',
		502: 'Bad Gateway',
		503: 'Service Unavailable'
	};

	const title = $derived(titles[errorStatus] || 'Something went wrong');
</script>

<div class="flex min-h-[60vh] flex-col items-center justify-center px-4 text-center">
	<div class="max-w-md space-y-4">
		<div class="text-6xl font-bold text-muted-foreground">{errorStatus}</div>
		<h1 class="text-2xl font-semibold tracking-tight">{title}</h1>
		<p class="text-muted-foreground">{errorMessage}</p>
		<div class="flex items-center justify-center gap-3 pt-2">
			<Button variant="outline" onclick={() => history.back()}>Go Back</Button>
			<Button href="/">Back to Home</Button>
		</div>
	</div>
</div>

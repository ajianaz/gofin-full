<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { authStore } from '$lib/stores/auth.svelte.js';

	onMount(async () => {
		if (authStore.isAuthenticated) {
			goto('/dashboard');
		} else {
			goto('/login');
		}
	});
</script>

<div class="flex h-screen items-center justify-center bg-background">
	<p class="text-muted-foreground">Redirecting...</p>
</div>

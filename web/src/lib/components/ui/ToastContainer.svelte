<script lang="ts">
	import { browser } from '$app/environment';

	type ToastType = 'success' | 'error' | 'info';

	interface Toast {
		id: number;
		type: ToastType;
		title: string;
		description?: string;
	}

	let toasts: Toast[] = $state([]);
	let counter = 0;

	const DISMISS_AFTER = 4000;

	function add(type: ToastType, title: string, description?: string) {
		if (!browser) return;
		const id = ++counter;
		toasts.push({ id, type, title, description });
		setTimeout(() => dismiss(id), DISMISS_AFTER);
	}

	function dismiss(id: number) {
		toasts = toasts.filter((t) => t.id !== id);
	}

	// Exported API
	export function toast(type: ToastType, title: string, description?: string) {
		add(type, title, description);
	}

	export function showErrorToast(description?: string, title = 'Error') {
		add('error', title, description || 'An unexpected error occurred.');
	}

	export function showSuccessToast(message: string, title = 'Success') {
		add('success', title, message);
	}

	export function handleApiError(err: unknown): void {
		if (!browser) return;

		if (err && typeof err === 'object') {
			const apiErr = err as { status?: number; detail?: string; message?: string };
			const status = apiErr.status;
			const detail = apiErr.detail ?? apiErr.message;

			if (status === 401) {
				add('error', 'Unauthorized', detail || 'Session expired — please log in again.');
			} else if (status === 403) {
				add('error', 'Forbidden', detail || 'You do not have permission to perform this action.');
			} else if (status && status >= 500) {
				add('error', 'Server Error', detail || 'A server error occurred. Please try again later.');
			} else {
				add('error', 'Error', detail || 'Something went wrong.');
			}
		} else if (err instanceof Error) {
			add('error', 'Network Error', err.message || 'Network error — check your connection.');
		} else {
			add('error', 'Error');
		}
	}
</script>

{#if toasts.length > 0}
	<div class="fixed bottom-4 right-4 z-[100] flex flex-col gap-2">
		{#each toasts as t (t.id)}
			<button
				class="flex items-start gap-3 rounded-lg border p-4 shadow-lg backdrop-blur-sm transition-all duration-300 animate-in slide-in-from-right-full cursor-pointer text-left w-80
					{t.type === 'error' ? 'border-red-500/30 bg-red-950/80 text-red-200' : ''}
					{t.type === 'success' ? 'border-green-500/30 bg-green-950/80 text-green-200' : ''}
					{t.type === 'info' ? 'border-blue-500/30 bg-blue-950/80 text-blue-200' : ''}"
				onclick={() => dismiss(t.id)}
				aria-live="polite"
			>
				<div class="flex-1 min-w-0">
					<p class="text-sm font-semibold">{t.title}</p>
					{#if t.description}
						<p class="text-xs mt-1 opacity-80">{t.description}</p>
					{/if}
				</div>
				<span class="text-xs opacity-50 mt-0.5">✕</span>
			</button>
		{/each}
	</div>
{/if}

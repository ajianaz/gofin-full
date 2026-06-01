/** Safe browser check that works in both SvelteKit and unit-test contexts. */
function isBrowser(): boolean {
	return typeof window !== 'undefined';
}

type ToastType = 'success' | 'error' | 'info';

// Minimal event-based toast system — no external dependencies.
// The ToastContainer.svelte component listens and renders toasts.
type ToastListener = (toast: { type: ToastType; title: string; description?: string }) => void;

const listeners = new Set<ToastListener>();

function emit(type: ToastType, title: string, description?: string) {
	if (!isBrowser()) return;
	const payload = { type, title, description };
	listeners.forEach((fn) => fn(payload));
}

/** Subscribe to toast events (used by ToastContainer). */
export function onToast(listener: ToastListener): () => void {
	listeners.add(listener);
	return () => listeners.delete(listener);
}

/**
 * Show an error toast.
 */
export function showErrorToast(description?: string, title = 'Error') {
	emit('error', title, description || 'An unexpected error occurred.');
}

/**
 * Show a success toast.
 */
export function showSuccessToast(message: string, title = 'Success') {
	emit('success', title, message);
}

/** Alias for showErrorToast (shim for sonner API consumers). */
export const toast = {
	error: showErrorToast,
	success: showSuccessToast,
	info: (msg: string, title?: string) => emit('info', title || 'Info', msg)
};

/**
 * Inspect an HTTP status code (or ApiError object) and fire the
 * appropriate toast automatically.
 */
export function handleApiError(err: unknown): void {
	if (!isBrowser()) return;

	if (err && typeof err === 'object') {
		const apiErr = err as { status?: number; detail?: string; message?: string };
		const status = apiErr.status;
		const detail = apiErr.detail ?? apiErr.message;

		if (status === 401) {
			showErrorToast(detail || 'Session expired — please log in again.', 'Unauthorized');
		} else if (status === 403) {
			showErrorToast(detail || 'You do not have permission to perform this action.', 'Forbidden');
		} else if (status && status >= 500) {
			showErrorToast(detail || 'A server error occurred. Please try again later.', 'Server Error');
		} else {
			showErrorToast(detail || 'Something went wrong.');
		}
	} else if (err instanceof Error) {
		showErrorToast(err.message || 'Network error — check your connection.', 'Network Error');
	} else {
		showErrorToast();
	}
}

import type { Handle } from '@sveltejs/kit';

/**
 * Content-Security-Policy header.
 * In development the policy is relaxed so Vite HMR and inline styles work.
 * In production the policy is strict: no unsafe-eval, no unsafe-inline scripts.
 */
function cspValue(isDev: boolean): string {
	if (isDev) {
		return [
			"default-src 'self'",
			"script-src 'self' 'unsafe-inline' 'unsafe-eval'",
			"style-src 'self' 'unsafe-inline'",
			"img-src 'self' data: blob:",
			"font-src 'self' data:",
			"connect-src 'self' ws: wss:",
			"frame-ancestors 'none'",
			"base-uri 'self'",
			"form-action 'self'"
		].join('; ');
	}
	return [
		"default-src 'self'",
		"script-src 'self'",
		"style-src 'self' 'unsafe-inline'",
		"img-src 'self' data: blob:",
		"font-src 'self' data:",
		"connect-src 'self'",
		"frame-ancestors 'none'",
		"base-uri 'self'",
		"form-action 'self'",
		"upgrade-insecure-requests"
	].join('; ');
}

export const handle: Handle = async ({ event, resolve }) => {
	const response = await resolve(event);

	// Security headers (defense-in-depth; primary headers set by Go API)
	response.headers.set('Content-Security-Policy', cspValue(import.meta.env.DEV));
	response.headers.set('X-Content-Type-Options', 'nosniff');
	response.headers.set('X-Frame-Options', 'DENY');
	response.headers.set('X-XSS-Protection', '1; mode=block');
	response.headers.set('Referrer-Policy', 'strict-origin-when-cross-origin');

	return response;
};

import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ cookies }) => {
	// Server-side hint for SSR routing decisions.
	// We check if httpOnly session cookies are present to avoid the
	// flash-of-login-page on authenticated routes.
	//
	// IMPORTANT: This is NOT an auth check — the Go API validates the actual
	// httpOnly cookie on every request. This merely provides a SSR hint so
	// the (app) layout doesn't redirect to /login when the user has a valid
	// server-side session. Even if spoofed, the worst case is a brief flash
	// of the app layout before the API returns 401 and the client redirects.
	const hasAccessToken = !!cookies.get('access_token');
	const hasRefreshToken = !!cookies.get('refresh_token');

	return {
		session: {
			// Only treat as potentially authenticated if access_token cookie exists.
			// Refresh-only means session expired — let client handle refresh flow.
			isAuthenticated: hasAccessToken || hasRefreshToken
		}
	};
};

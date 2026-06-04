import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ cookies }) => {
	// Server-side auth check: detect if httpOnly session cookies exist.
	// This gives SSR a hint about authentication state, preventing the
	// flash-of-login-page that occurs when the client-side store defaults
	// to unauthenticated before hydration.
	//
	// The actual auth validation happens in the Go API via the httpOnly cookie.
	// Here we just check cookie presence for SSR routing decisions.
	const hasAccessToken = !!cookies.get('access_token');
	const hasRefreshToken = !!cookies.get('refresh_token');

	return {
		session: {
			// On the server, we can't validate the token without calling the API,
			// but we can check if cookies exist. If refresh token exists, assume
			// the user might be authenticated (API will validate on actual requests).
			isAuthenticated: hasAccessToken || hasRefreshToken
		}
	};
};

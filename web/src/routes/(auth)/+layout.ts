import { redirect } from '@sveltejs/kit';
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = async ({ parent }) => {
	const { session } = await parent();

	// If already authenticated, redirect away from auth pages
	if (session.isAuthenticated) {
		throw redirect(302, '/dashboard');
	}

	return {};
};

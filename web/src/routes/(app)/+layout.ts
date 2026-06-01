import { redirect } from '@sveltejs/kit';
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = async ({ parent }) => {
	const { session } = await parent();

	// If the user has no valid session, redirect to login
	if (!session.isAuthenticated) {
		throw redirect(302, '/login');
	}

	return {};
};

import { redirect } from '@sveltejs/kit';
import type { LayoutLoad } from './$types';
import { apiGet } from '$lib/services/apiFetch';
import { auth } from '$lib/stores/auth';

export const prerender = false;
export const ssr = false;

export const load: LayoutLoad = async ({ url }) => {
	// Don't guard the login page itself
	if (url.pathname === '/login') {
		return {};
	}

	try {
		// Fetch the current user from the backend
		const user = await apiGet<{ id: string; email: string; role: string }>('/auth/me');

		// Populate the auth store
		auth.setUser({
			id: user.id,
			email: user.email,
			role: user.role as 'admin' | 'investigator' | 'viewer'
		});

		return {};
	} catch (error) {
		// Clear auth state and redirect to login on any error (401 or otherwise)
		auth.clear();
		throw redirect(302, '/login');
	}
}

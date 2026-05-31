import { redirect } from '@sveltejs/kit';
import type { LayoutLoad } from './$types';
import { apiGet } from '$lib/services/apiFetch';
import { auth } from '$lib/stores/auth';

export const prerender = false;
export const ssr = false;

export const load: LayoutLoad = async () => {
	try {
		const user = await apiGet<{ id: string; email: string; role: string }>('/auth/me');

		auth.setUser({
			id: user.id,
			email: user.email,
			role: user.role as 'admin' | 'investigator' | 'viewer'
		});

		return {};
	} catch {
		auth.clear();
		throw redirect(302, '/login');
	}
};

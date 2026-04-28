import { redirect } from '@sveltejs/kit';
import { dev } from '$app/environment';
import { validateClientToken, createClientSession } from '$lib/server/client-access';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params, cookies }) => {
	const validation = await validateClientToken(params.token);

	if (validation.valid) {
		await createClientSession(validation.caseId);
		cookies.set('ca_session', validation.caseId, {
			httpOnly: true,
			sameSite: 'strict',
			secure: !dev,
			path: '/',
			maxAge: 60 * 60 * 24
		});
		redirect(303, `/client-access/${params.token}/files`);
	}

	return { error: validation.reason };
};

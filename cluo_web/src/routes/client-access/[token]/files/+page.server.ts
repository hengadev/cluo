import { redirect } from '@sveltejs/kit';
import { getClientCaseSummary } from '$lib/server/client-access';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params, cookies }) => {
	const caseId = cookies.get('ca_session');

	if (!caseId) {
		redirect(303, `/client-access/${params.token}`);
	}

	const summary = await getClientCaseSummary(caseId);

	return { summary, token: params.token };
};

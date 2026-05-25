import { redirect } from '@sveltejs/kit';
import { validateClientToken, getTokenMedia, getReportHtml } from '$lib/server/client-access';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params }) => {
	const validation = await validateClientToken(params.token);

	if (!validation.valid) {
		// Token invalid/expired — send back to landing which shows the right error card
		redirect(303, `/client-access/${params.token}`);
	}

	// Check media availability (fail open: show tab on error)
	const mediaResult = await getTokenMedia(params.token);
	const hasMedia = mediaResult === null ? true : mediaResult.length > 0;

	const reportResult = await getReportHtml(params.token);

	return {
		caseData: validation.caseData,
		token: params.token,
		hasMedia,
		rapportHtml: reportResult.status === 'ok' ? reportResult.html : null,
		rapportError: reportResult.status === 'error',
	};
};

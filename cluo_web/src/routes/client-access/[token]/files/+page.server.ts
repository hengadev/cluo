import { redirect } from '@sveltejs/kit';
import { validateClientToken, getTokenMedia, getReportHtml, getDocumentsByToken } from '$lib/server/client-access';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params }) => {
	const validation = await validateClientToken(params.token);

	if (!validation.valid) {
		// Token invalid/expired — send back to landing which shows the right error card
		redirect(303, `/client-access/${params.token}`);
	}

	const [mediaResult, reportResult, documentsResult] = await Promise.all([
		getTokenMedia(params.token),
		getReportHtml(params.token),
		getDocumentsByToken(params.token),
	]);

	// Check media availability (fail open: show tab on error)
	const hasMedia = mediaResult === null ? true : mediaResult.length > 0;
	const media = mediaResult ?? [];

	return {
		caseData: validation.caseData,
		token: params.token,
		hasMedia,
		media,
		rapportHtml: reportResult.status === 'ok' ? reportResult.html : null,
		rapportError: reportResult.status === 'error',
		documents: documentsResult.status === 'ok' ? documentsResult.documents : [],
		documentsError: documentsResult.status === 'error',
	};
};

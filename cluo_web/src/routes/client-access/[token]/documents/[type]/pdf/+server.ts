import { redirect } from '@sveltejs/kit';
import { validateClientToken, getApiBaseUrl } from '$lib/server/client-access';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({ params }) => {
	const validation = await validateClientToken(params.token);
	if (!validation.valid) {
		redirect(302, `/client-access/${encodeURIComponent(params.token)}/files?tab=documents&pdf_error=auth`);
	}

	const validTypes = ['estimate', 'mandate', 'contract', 'invoice'];
	if (!validTypes.includes(params.type)) {
		redirect(302, `/client-access/${encodeURIComponent(params.token)}/files?tab=documents&pdf_error=invalid_type`);
	}

	const res = await fetch(
		`${getApiBaseUrl()}/token/${encodeURIComponent(params.token)}/documents/${encodeURIComponent(params.type)}/pdf`
	);

	if (res.status === 404) {
		redirect(
			302,
			`/client-access/${encodeURIComponent(params.token)}/files?tab=documents&pdf_error=not_found`
		);
	}

	if (!res.ok) {
		redirect(
			302,
			`/client-access/${encodeURIComponent(params.token)}/files?tab=documents&pdf_error=server`
		);
	}

	const pdfBytes = await res.arrayBuffer();

	return new Response(pdfBytes, {
		headers: {
			'Content-Type': 'application/pdf',
			'Content-Disposition': `attachment; filename="${params.type}.pdf"`
		}
	});
};

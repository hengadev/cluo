import { getApiBaseUrl } from '$lib/server/client-access';
import { error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({ params }) => {
	const apiBase = getApiBaseUrl();

	try {
		const res = await fetch(`${apiBase}/token/${encodeURIComponent(params.token)}/media/archive`);
		if (!res.ok) {
			if (res.status === 401) {
				error(401, { message: 'Token expiré ou révoqué' });
			}
			error(500, { message: 'Failed to download media archive' });
		}

		const disposition = res.headers.get('Content-Disposition') || 'attachment; filename="medias.zip"';

		return new Response(res.body, {
			headers: {
				'Content-Type': 'application/zip',
				'Content-Disposition': disposition
			}
		});
	} catch (e) {
		if (e && typeof e === 'object' && 'status' in e) throw e;
		error(500, { message: 'Failed to download media archive' });
	}
};

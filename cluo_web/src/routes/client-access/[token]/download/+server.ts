import { error } from '@sveltejs/kit';
import { streamCaseArchive } from '$lib/server/client-access';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({ cookies }) => {
	const caseId = cookies.get('ca_session');

	if (!caseId) {
		error(401, 'Unauthorized');
	}

	const stream = await streamCaseArchive(caseId);

	// TODO: Log download event

	return new Response(stream, {
		headers: {
			'Content-Type': 'application/zip',
			'Content-Disposition': 'attachment; filename="case-files.zip"'
		}
	});
};

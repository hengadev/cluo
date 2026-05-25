import { error } from '@sveltejs/kit';
import { streamCaseArchive } from '$lib/server/client-access';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({ params }) => {
	try {
		const stream = await streamCaseArchive(params.token);

		return new Response(stream, {
			headers: {
				'Content-Type': 'application/zip',
				'Content-Disposition': 'attachment; filename="case-files.zip"'
			}
		});
	} catch {
		error(500, 'Failed to generate archive');
	}
};

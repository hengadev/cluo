import { error } from '@sveltejs/kit';
import { streamCaseArchive } from '$lib/server/client-access';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({ params }) => {
	try {
		const { stream, disposition } = await streamCaseArchive(params.token);

		return new Response(stream, {
			headers: {
				'Content-Type': 'application/zip',
				'Content-Disposition': disposition
			}
		});
	} catch (e) {
		if (e && typeof e === 'object' && 'status' in e) throw e;
		error(500, 'Failed to generate archive');
	}
};

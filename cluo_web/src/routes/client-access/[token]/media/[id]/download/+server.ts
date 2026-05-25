import { error } from '@sveltejs/kit';
import { validateClientToken, streamMediaFile } from '$lib/server/client-access';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({ params }) => {
	const validation = await validateClientToken(params.token);
	if (!validation.valid) {
		error(401, 'Token invalide ou expiré');
	}

	const result = await streamMediaFile(params.token, params.id);
	if (!result) {
		error(404, 'Média introuvable');
	}

	return new Response(result.stream, {
		headers: {
			'Content-Type': result.mimeType,
			'Content-Disposition': `attachment; filename*=UTF-8''${encodeURIComponent(result.fileName)}`,
			'Content-Length': String(result.fileSize),
		},
	});
};

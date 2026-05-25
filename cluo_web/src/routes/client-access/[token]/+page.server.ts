import { validateClientToken } from '$lib/server/client-access';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params }) => {
	const validation = await validateClientToken(params.token);

	if (validation.valid) {
		return {
			valid: true as const,
			caseData: validation.caseData,
			token: params.token
		};
	}

	return {
		valid: false as const,
		error: validation.reason
	};
};

import type { PageServerLoad } from './$types';
import { fetchProcessingStatus } from '$lib/api/server';

export const load: PageServerLoad = async ({ params, fetch }) => {
    return await fetchProcessingStatus(fetch, params.id);
};

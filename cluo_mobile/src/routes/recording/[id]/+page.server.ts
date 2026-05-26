import type { PageServerLoad } from './$types';
import { fetchMediaById } from '$lib/api/server';
import type { Transcript, AnalysisResult } from '$lib/types/recording';

export const load: PageServerLoad = async ({ params, fetch }) => {
    const recording = await fetchMediaById(fetch, params.id);
    return {
        recording,
        transcript: null as Transcript | null,
        analysis: null as AnalysisResult | null,
        error: null as string | null,
    };
};

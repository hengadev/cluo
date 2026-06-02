import type { PageLoad } from "./$types";
import { listRecordings, getCurrentCase } from "$lib/api";
import type { Recording } from "$lib/types/recording";
import type { Case } from "$lib/types/case";

interface HomePageData {
	recordings: Recording[];
	totalCount: number;
	currentCase: Case | null;
	error: string | null;
}

export const load: PageLoad = async ({ url }): Promise<HomePageData> => {
	const caseId = url.searchParams.get("caseId") ?? undefined;

	try {
		const [recordingsData, currentCase] = await Promise.all([
			listRecordings({ caseId }),
			getCurrentCase(caseId),
		]);
		return { recordings: recordingsData.recordings, totalCount: recordingsData.totalCount, currentCase, error: null };
	} catch (error) {
		return {
			recordings: [],
			totalCount: 0,
			currentCase: null,
			error: error instanceof Error ? error.message : "Échec du chargement des enregistrements",
		};
	}
};

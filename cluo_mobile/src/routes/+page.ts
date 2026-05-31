import type { PageLoad } from "./$types";
import { listRecordings, getCurrentCase } from "$lib/api";
import type { Recording } from "$lib/types/recording";
import type { Case } from "$lib/types/case";

interface HomePageData {
	recordings: Recording[];
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
		return { recordings: recordingsData.recordings, currentCase, error: null };
	} catch (error) {
		return {
			recordings: [],
			currentCase: null,
			error: error instanceof Error ? error.message : "Échec du chargement des enregistrements",
		};
	}
};

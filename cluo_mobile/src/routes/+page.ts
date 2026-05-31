import type { PageLoad } from "./$types";
import { listRecordings } from "$lib/api";
import type { Recording } from "$lib/types/recording";

interface HomePageData {
	recordings: Recording[];
	error: string | null;
}

export const load: PageLoad = async ({ url }): Promise<HomePageData> => {
	// caseId can be supplied as a query param: /?caseId=<uuid>
	// Falls back to the last-used caseId stored in localStorage (see audio.ts).
	const caseId = url.searchParams.get("caseId") ?? undefined;

	try {
		const data = await listRecordings({ caseId });
		return { recordings: data.recordings, error: null };
	} catch (error) {
		return {
			recordings: [],
			error: error instanceof Error ? error.message : "Échec du chargement des enregistrements",
		};
	}
};

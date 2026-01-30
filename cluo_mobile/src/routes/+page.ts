import type { PageLoad } from "./$types";
import { listRecordings } from "$lib/api";
import type { Recording } from "$lib/types/recording";

interface HomePageData {
	recordings: Recording[];
	error: string | null;
}

export const load: PageLoad = async (): Promise<HomePageData> => {
	try {
		const data = await listRecordings();
		return { recordings: data.recordings, error: null };
	} catch (error) {
		return {
			recordings: [],
			error: error instanceof Error ? error.message : "Failed to load recordings",
		};
	}
};

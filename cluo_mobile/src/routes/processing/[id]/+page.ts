import type { PageLoad } from "./$types";
import { getRecordingStatus } from "$lib/api";
import type { ProcessingStep } from "$lib/types/recording";

interface ProcessingData {
	recordingId: string;
	steps: ProcessingStep[];
	error: string | null;
}

export const load: PageLoad = async ({ params }): Promise<ProcessingData> => {
	const recordingId = params.id;

	try {
		const status = await getRecordingStatus(recordingId);
		return {
			recordingId,
			steps: status.processingSteps,
			error: null,
		};
	} catch (error) {
		return {
			recordingId,
			steps: [
				{ title: "Téléchargement audio", status: "pending" },
				{ title: "Traitement de la transcription", status: "pending" },
				{ title: "Génération du résumé", status: "pending" },
				{ title: "Terminé", status: "pending" },
			],
			error: error instanceof Error ? error.message : "Failed to load processing status",
		};
	}
};

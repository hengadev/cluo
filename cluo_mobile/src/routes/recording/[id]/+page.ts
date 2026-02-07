import type { PageLoad } from "./$types";
import { getRecording } from "$lib/api";
import type { Recording, Transcript, AnalysisResult } from "$lib/types/recording";

interface RecordingData {
	recording: Recording & { audioUrl?: string };
	transcript: Transcript | null;
	analysis: AnalysisResult | null;
	error: string | null;
}

export const load: PageLoad = async ({ params }): Promise<RecordingData> => {
	const recordingId = params.id;

	try {
		const data = await getRecording(recordingId);
		return {
			recording: data.recording,
			transcript: data.transcript,
			analysis: data.analysis,
			error: null,
		};
	} catch (error) {
		return {
			recording: {
				id: recordingId,
				title: "Recording",
				date: new Date().toLocaleDateString("en-GB", {
					day: "2-digit",
					month: "short",
					year: "numeric",
				}),
				startTime: "00:00",
				duration: 0,
				status: "failed",
			},
			transcript: null,
			analysis: null,
			error: error instanceof Error ? error.message : "Failed to load recording",
		};
	}
};

import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params, fetch }) => {
    const recordingId = params.id;

    // TODO: Replace with actual API call to your Golang backend
    // Example:
    // const response = await fetch(`${import.meta.env.VITE_API_URL}/recordings/${recordingId}`);
    // if (!response.ok) {
    //     throw error(response.status, 'Failed to load recording');
    // }
    // const recording = await response.json();

    // For now, returning placeholder data
    return {
        recording: {
            id: recordingId,
            title: "The title of the recording",
            date: "01 Jan, 2025",
            startTime: "00:00",
            duration: "05:23",
            audioUrl: "", // URL to the audio file
            transcript: "This is where the transcript would appear...",
            tags: ["Important", "Follow-up"],
        }
    };
};

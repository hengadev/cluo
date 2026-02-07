import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params, fetch }) => {
    const recordingId = params.id;

    // TODO: Replace with actual API call to your Golang backend
    // to fetch the current processing status
    // Example:
    // const response = await fetch(`${import.meta.env.VITE_API_URL}/recordings/${recordingId}/status`);
    // if (!response.ok) {
    //     throw error(response.status, 'Failed to load processing status');
    // }
    // const status = await response.json();

    // For now, returning mock data with first step completed
    return {
        recordingId,
        steps: [
            {
                title: "Téléchargement audio",
                status: "completed" as const,
            },
            {
                title: "Traitement de la transcription",
                status: "processing" as const,
            },
            {
                title: "Génération du résumé",
                status: "processing" as const,
            },
            {
                title: "Terminé",
                status: "processing" as const,
            },
        ]
    };
};

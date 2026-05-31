/**
 * API facade that switches between real and mock implementations
 * based on the VITE_MOCK_MODE environment variable.
 */

import * as realApi from "./audio";
import * as mockApi from "./mock";

/**
 * Check if mock mode is enabled.
 * Returns true if VITE_MOCK_MODE is "true" (string comparison).
 */
export const isMockMode = (): boolean => {
    return import.meta.env.VITE_MOCK_MODE === "true";
};

/**
 * Get the appropriate API implementation based on mock mode.
 */
function getApi() {
    return isMockMode() ? mockApi : realApi;
}

// Re-export all API functions with automatic mock mode switching

export async function uploadRecording(
    blob: Blob,
    metadata?: { caseId?: string; title?: string }
) {
    return getApi().uploadRecording(blob, metadata);
}

export async function getRecordingStatus(id: string) {
    return getApi().getRecordingStatus(id);
}

export async function getTranscript(id: string) {
    return getApi().getTranscript(id);
}

export async function confirmTranscript(id: string, text: string) {
    return getApi().confirmTranscript(id, text);
}

export async function analyzeTranscript(id: string) {
    return getApi().analyzeTranscript(id);
}

export async function getAnalysis(id: string) {
    return getApi().getAnalysis(id);
}

export async function listRecordings(options?: {
    limit?: number;
    offset?: number;
    caseId?: string;
    status?: string;
}) {
    return getApi().listRecordings(options);
}

export async function deleteRecording(id: string) {
    return getApi().deleteRecording(id);
}

/**
 * Get a single recording with its transcript and analysis.
 * Note: Real API implementation needs to be added to audio.ts
 */
export async function getRecording(id: string) {
    if (isMockMode()) {
        return mockApi.getRecording(id);
    }
    // For real API, we need to make multiple calls
    // This should be implemented as a single endpoint in the backend ideally
    return realApi.getRecording(id);
}

/**
 * Get the audio URL for a recording.
 */
export async function getAudioUrl(id: string) {
    if (isMockMode()) {
        return mockApi.getAudioUrl(id);
    }
    return realApi.getAudioUrl(id);
}

/**
 * Get the current active case.
 */
export async function getCurrentCase(caseId?: string) {
    return getApi().getCurrentCase(caseId);
}

// Re-export types
export type {
    UploadRecordingResponse,
    RecordingStatusResponse,
    TranscriptResponse,
    AnalysisResponse,
    RecordingsListResponse,
} from "../types/recording";
export type { Case } from "../types/case";

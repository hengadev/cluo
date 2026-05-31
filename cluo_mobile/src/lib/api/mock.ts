/**
 * Mock API implementations for local development.
 * These functions simulate realistic API responses with delays.
 */

import type {
    Recording,
    Transcript,
    AnalysisResult,
    ProcessingStep,
    RecordingStatus,
    UploadRecordingResponse,
    RecordingStatusResponse,
    TranscriptResponse,
    AnalysisResponse,
    RecordingsListResponse,
} from "../types/recording";

// Simulate network delay
const delay = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

// Mock data store (in-memory state for the session)
const mockRecordings: Map<string, Recording & { audioUrl?: string }> = new Map([
    [
        "1",
        {
            id: "1",
            title: "Interview with witness #1",
            date: "30 Jan, 2026",
            startTime: "14:30",
            duration: 323,
            status: "completed",
            audioUrl: "",
        },
    ],
    [
        "2",
        {
            id: "2",
            title: "Site visit notes",
            date: "30 Jan, 2026",
            startTime: "10:15",
            duration: 225,
            status: "completed",
            audioUrl: "",
        },
    ],
    [
        "3",
        {
            id: "3",
            title: "Phone call recording",
            date: "29 Jan, 2026",
            startTime: "16:00",
            duration: 730,
            status: "transcribing",
            audioUrl: "",
        },
    ],
    [
        "4",
        {
            id: "4",
            title: "Client meeting summary",
            date: "28 Jan, 2026",
            startTime: "09:00",
            duration: 1845,
            status: "completed",
            audioUrl: "",
        },
    ],
    [
        "5",
        {
            id: "5",
            title: "Field inspection audio",
            date: "27 Jan, 2026",
            startTime: "11:30",
            duration: 456,
            status: "analyzing",
            audioUrl: "",
        },
    ],
]);

const mockTranscripts: Map<string, Transcript> = new Map([
    [
        "1",
        {
            recordingId: "1",
            text: "This is the transcript for the witness interview. The witness stated that they observed the incident from approximately 50 meters away. They described the sequence of events in detail, noting the time was around 2:30 PM when the incident began.\n\nThe witness mentioned that visibility was good and there were no obstructions. They were able to clearly see the individuals involved and provided descriptions of their clothing and approximate heights.\n\nWhen asked about any sounds they heard, the witness recalled hearing loud voices followed by what they described as a 'crashing sound.' They immediately called emergency services after witnessing the event.",
            confidence: 0.94,
            isConfirmed: true,
            createdAt: "2026-01-30T14:35:00Z",
            updatedAt: "2026-01-30T14:40:00Z",
        },
    ],
    [
        "2",
        {
            recordingId: "2",
            text: "Site visit conducted on January 30th. The property is located at the corner of Main Street and Oak Avenue. Upon arrival, I noted the following conditions:\n\nThe exterior of the building shows signs of water damage on the north-facing wall. Paint is peeling in several areas and there appears to be mold growth near the foundation.\n\nThe roof has several missing shingles and the gutters are clogged with debris. Drainage appears to be inadequate, with water pooling near the entrance.\n\nInterior inspection revealed humidity issues in the basement. The electrical panel is outdated and may not meet current code requirements. Recommend full inspection by licensed electrician.",
            confidence: 0.91,
            isConfirmed: false,
            createdAt: "2026-01-30T10:20:00Z",
            updatedAt: "2026-01-30T10:20:00Z",
        },
    ],
    [
        "4",
        {
            recordingId: "4",
            text: "Client meeting held to discuss project timeline and deliverables. Attendees: John Smith (client), Sarah Johnson (project manager), Michael Chen (lead developer).\n\nKey discussion points:\n1. Project deadline extended by two weeks to accommodate additional feature requests\n2. Budget increase of 15% approved for enhanced security measures\n3. Weekly status meetings scheduled for every Tuesday at 10 AM\n\nAction items:\n- Sarah to send updated project plan by Friday\n- Michael to provide technical specifications for new features\n- John to confirm stakeholder availability for demo session\n\nNext meeting scheduled for February 4th to review progress.",
            confidence: 0.96,
            isConfirmed: true,
            createdAt: "2026-01-28T09:35:00Z",
            updatedAt: "2026-01-28T10:00:00Z",
        },
    ],
]);

const mockAnalysis: Map<string, AnalysisResult> = new Map([
    [
        "1",
        {
            id: "mock-analysis-1",
            transcriptionId: "mock-transcription-1",
            keyFindings: "Witness was positioned 50 meters from the incident location.\nIncident occurred at approximately 2:30 PM.\nVisibility conditions were good with no obstructions.\nWitness heard loud voices followed by a crashing sound.",
            summary: "Witness statement regarding incident on January 30th. Witness had clear line of sight and was able to provide detailed descriptions of individuals involved.",
            sentiment: "neutral",
            topics: JSON.stringify(["witness statement", "incident report", "physical descriptions"]),
            suggestedActions: "Follow up with witness for detailed physical descriptions.\nVerify incident timeline with other witnesses.\nRequest CCTV footage from the area.",
            createdAt: "2026-01-30T14:45:00Z",
        },
    ],
    [
        "4",
        {
            id: "mock-analysis-4",
            transcriptionId: "mock-transcription-4",
            keyFindings: "Project deadline extended by two weeks.\nBudget increased by 15% for security measures.\nWeekly status meetings scheduled every Tuesday at 10 AM.",
            summary: "Client meeting to discuss project timeline and deliverables. Several key decisions were made regarding budget, timeline, and communication cadence.",
            sentiment: "positive",
            topics: JSON.stringify(["project management", "budget", "timeline", "meetings"]),
            suggestedActions: "Sarah to send updated project plan by Friday.\nMichael to provide technical specifications for new features.\nJohn to confirm stakeholder availability for demo session.",
            createdAt: "2026-01-28T10:15:00Z",
        },
    ],
]);

// Track processing state for new uploads
const processingState: Map<string, { step: number; startTime: number }> = new Map();

/**
 * Format duration from seconds to MM:SS string
 */
function formatDuration(seconds: number): string {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins.toString().padStart(2, "0")}:${secs.toString().padStart(2, "0")}`;
}

/**
 * Mock: Upload a recording
 */
export async function uploadRecording(
    blob: Blob,
    metadata?: { caseId?: string; title?: string }
): Promise<UploadRecordingResponse> {
    await delay(800);

    const id = `mock-${Date.now()}`;
    const now = new Date();

    // Create new recording entry
    const recording: Recording & { audioUrl?: string } = {
        id,
        title: metadata?.title || `Recording ${now.toLocaleTimeString()}`,
        caseId: metadata?.caseId,
        date: now.toLocaleDateString("en-GB", {
            day: "2-digit",
            month: "short",
            year: "numeric",
        }),
        startTime: now.toLocaleTimeString("en-GB", {
            hour: "2-digit",
            minute: "2-digit",
        }),
        duration: 0, // Will be updated
        status: "uploading",
        audioUrl: URL.createObjectURL(blob),
    };

    mockRecordings.set(id, recording);
    processingState.set(id, { step: 0, startTime: Date.now() });

    return {
        id,
        status: "uploading",
    };
}

/**
 * Mock: Get recording status with processing steps
 */
export async function getRecordingStatus(id: string): Promise<RecordingStatusResponse> {
    await delay(300);

    const recording = mockRecordings.get(id);
    if (!recording) {
        throw new Error("Recording not found");
    }

    // Simulate processing progression
    const state = processingState.get(id);
    const steps: ProcessingStep[] = [
        { title: "Téléchargement audio", status: "pending" },
        { title: "Traitement de la transcription", status: "pending" },
        { title: "Génération du résumé", status: "pending" },
        { title: "Terminé", status: "pending" },
    ];

    if (state) {
        const elapsed = Date.now() - state.startTime;
        // Progress through steps every 3 seconds
        const currentStep = Math.min(Math.floor(elapsed / 3000), 4);

        for (let i = 0; i < steps.length; i++) {
            if (i < currentStep) {
                steps[i].status = "completed";
            } else if (i === currentStep) {
                steps[i].status = "processing";
            }
        }

        // Update recording status based on progress
        let status: RecordingStatus = "uploading";
        if (currentStep >= 1) status = "transcribing";
        if (currentStep >= 2) status = "analyzing";
        if (currentStep >= 4) {
            status = "completed";
            processingState.delete(id);

            // Create mock transcript when completed
            if (!mockTranscripts.has(id)) {
                mockTranscripts.set(id, {
                    recordingId: id,
                    text: "This is an automatically generated transcript for your recording. The audio has been processed and converted to text. Please review and edit as needed before confirming.\n\nNote: This is mock data for development purposes.",
                    confidence: 0.89,
                    isConfirmed: false,
                    createdAt: new Date().toISOString(),
                    updatedAt: new Date().toISOString(),
                });
            }
        }

        recording.status = status;
        mockRecordings.set(id, recording);

        return {
            id,
            status,
            processingSteps: steps,
        };
    }

    // For pre-existing recordings, return based on their status
    if (recording.status === "completed") {
        steps.forEach((s) => (s.status = "completed"));
    }

    return {
        id,
        status: recording.status,
        processingSteps: steps,
    };
}

/**
 * Mock: Get transcript
 */
export async function getTranscript(id: string): Promise<TranscriptResponse> {
    await delay(400);

    const transcript = mockTranscripts.get(id);
    if (!transcript) {
        throw new Error("Transcript not found. Processing may not be complete.");
    }

    return transcript;
}

/**
 * Mock: Confirm/update transcript
 */
export async function confirmTranscript(id: string, text: string): Promise<void> {
    await delay(500);

    const transcript = mockTranscripts.get(id);
    if (!transcript) {
        throw new Error("Transcript not found");
    }

    transcript.text = text;
    transcript.isConfirmed = true;
    transcript.updatedAt = new Date().toISOString();
    mockTranscripts.set(id, transcript);
}

/**
 * Mock: Trigger analysis
 */
export async function analyzeTranscript(id: string): Promise<void> {
    await delay(600);

    const transcript = mockTranscripts.get(id);
    if (!transcript) {
        throw new Error("Transcript not found");
    }

    // Create mock analysis if it doesn't exist
    if (!mockAnalysis.has(id)) {
        mockAnalysis.set(id, {
            id: `mock-analysis-${id}`,
            transcriptionId: `mock-transcription-${id}`,
            keyFindings: "Key observation extracted from the transcript.\nImportant detail identified in the recording.",
            summary: "Auto-generated summary of the recording content. Review and verify before use.",
            sentiment: "neutral",
            topics: JSON.stringify(["recording", "notes", "review"]),
            suggestedActions: "Review the transcript for accuracy.\nFollow up on any unclear sections.",
            createdAt: new Date().toISOString(),
        });
    }
}

/**
 * Mock: Get analysis results
 */
export async function getAnalysis(id: string): Promise<AnalysisResponse> {
    await delay(400);

    const analysis = mockAnalysis.get(id);
    if (!analysis) {
        throw new Error("Analysis not found. Please run analysis first.");
    }

    return analysis;
}

/**
 * Mock: List all recordings
 */
export async function listRecordings(options?: {
    limit?: number;
    offset?: number;
    caseId?: string;
    status?: string;
}): Promise<RecordingsListResponse> {
    await delay(500);

    let recordings = Array.from(mockRecordings.values()).map((r) => ({
        ...r,
        duration: formatDuration(typeof r.duration === "number" ? r.duration : 0),
    }));

    // Apply filters
    if (options?.caseId) {
        recordings = recordings.filter((r) => r.caseId === options.caseId);
    }
    if (options?.status) {
        recordings = recordings.filter((r) => r.status === options.status);
    }

    // Sort by date (newest first)
    recordings.sort((a, b) => {
        const dateA = new Date(a.date).getTime();
        const dateB = new Date(b.date).getTime();
        return dateB - dateA;
    });

    const totalCount = recordings.length;

    // Apply pagination
    const offset = options?.offset ?? 0;
    const limit = options?.limit ?? 20;
    recordings = recordings.slice(offset, offset + limit);

    return {
        recordings: recordings as Recording[],
        totalCount,
    };
}

/**
 * Mock: Delete recording
 */
export async function deleteRecording(id: string): Promise<void> {
    await delay(400);

    if (!mockRecordings.has(id)) {
        throw new Error("Recording not found");
    }

    mockRecordings.delete(id);
    mockTranscripts.delete(id);
    mockAnalysis.delete(id);
    processingState.delete(id);
}

/**
 * Mock: Get single recording with details
 */
export async function getRecording(id: string): Promise<{
    recording: Recording & { audioUrl?: string };
    transcript: Transcript | null;
    analysis: AnalysisResult | null;
}> {
    await delay(400);

    const recording = mockRecordings.get(id);
    if (!recording) {
        throw new Error("Recording not found");
    }

    const transcript = mockTranscripts.get(id) ?? null;
    const analysis = mockAnalysis.get(id) ?? null;

    return {
        recording: {
            ...recording,
            duration: formatDuration(typeof recording.duration === "number" ? recording.duration : 0),
        } as Recording & { audioUrl?: string },
        transcript,
        analysis,
    };
}

/**
 * Mock: Get audio URL for a recording
 */
export async function getAudioUrl(id: string): Promise<string> {
    await delay(200);

    const recording = mockRecordings.get(id);
    if (!recording) {
        throw new Error("Recording not found");
    }

    return recording.audioUrl ?? "";
}

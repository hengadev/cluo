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
import type { Case } from "../types/case";

// Simulate network delay
const delay = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

// Mock data store (in-memory state for the session)
const mockRecordings: Map<string, Recording & { audioUrl?: string }> = new Map([
    [
        "1",
        {
            id: "1",
            title: "Interview avec le témoin n°1",
            date: "30 janv. 2026",
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
            title: "Notes de visite sur site",
            date: "30 janv. 2026",
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
            title: "Enregistrement d'appel téléphonique",
            date: "29 janv. 2026",
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
            title: "Compte-rendu de réunion client",
            date: "28 janv. 2026",
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
            title: "Audio d'inspection sur le terrain",
            date: "27 janv. 2026",
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
            text: "Ceci est la transcription de l'interview du témoin. Le témoin a déclaré avoir observé l'incident depuis environ 50 mètres de distance. Il a décrit la séquence des événements en détail, notant que l'heure était autour de 14h30 lorsque l'incident a commencé.\n\nLe témoin a mentionné que la visibilité était bonne et qu'il n'y avait aucun obstacle. Il a pu voir clairement les personnes impliquées et a fourni des descriptions de leurs vêtements et de leurs tailles approximatives.\n\nLorsqu'on lui a demandé s'il avait entendu des bruits, le témoin s'est souvenu avoir entendu des voix fortes suivies de ce qu'il a décrit comme un « bruit de collision ». Il a immédiatement appelé les services d'urgence après avoir assisté à l'événement.",
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
            text: "Visite sur site effectuée le 30 janvier. La propriété est située à l'angle de la rue Principale et de l'avenue du Chêne. À l'arrivée, j'ai constaté les conditions suivantes :\n\nL'extérieur du bâtiment présente des signes de dégâts des eaux sur le mur nord. La peinture s'écaille à plusieurs endroits et il semble y avoir de la moisissure près des fondations.\n\nLe toit a plusieurs bardeaux manquants et les gouttières sont obstruées par des débris. Le drainage semble insuffisant, avec de l'eau stagnante près de l'entrée.\n\nL'inspection intérieure a révélé des problèmes d'humidité dans le sous-sol. Le tableau électrique est obsolète et pourrait ne pas répondre aux exigences actuelles. Recommandation de procéder à une inspection complète par un électricien agréé.",
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
            text: "Réunion client tenue pour discuter du calendrier du projet et des livrables. Participants : John Smith (client), Sarah Johnson (chef de projet), Michael Chen (développeur principal).\n\nPoints de discussion clés :\n1. Délai du projet prolongé de deux semaines pour accommodate les demandes de fonctionnalités supplémentaires\n2. Augmentation du budget de 15 % approuvée pour des mesures de sécurité renforcées\n3. Réunions d'avancement hebdomadaires planifiées chaque mardi à 10h\n\nPoints d'action :\n- Sarah doit envoyer le plan de projet mis à jour d'ici vendredi\n- Michael doit fournir les spécifications techniques pour les nouvelles fonctionnalités\n- John doit confirmer la disponibilité des parties prenantes pour la session de démonstration\n\nProchaine réunion prévue le 4 février pour faire le point sur l'avancement.",
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
            keyFindings: "Le témoin était positionné à 50 mètres du lieu de l'incident.\nL'incident s'est produit à environ 14h30.\nLes conditions de visibilité étaient bonnes sans aucun obstacle.\nLe témoin a entendu des voix fortes suivies d'un bruit de collision.",
            summary: "Déclaration du témoin concernant l'incident du 30 janvier. Le témoin avait une ligne de vue dégagée et a pu fournir des descriptions détaillées des personnes impliquées.",
            sentiment: "neutre",
            topics: JSON.stringify(["déclaration de témoin", "rapport d'incident", "descriptions physiques"]),
            suggestedActions: "Assurer un suivi avec le témoin pour des descriptions physiques détaillées.\nVérifier le calendrier de l'incident avec d'autres témoins.\nDemander les images de vidéosurveillance de la zone.",
            createdAt: "2026-01-30T14:45:00Z",
        },
    ],
    [
        "4",
        {
            id: "mock-analysis-4",
            transcriptionId: "mock-transcription-4",
            keyFindings: "Délai du projet prolongé de deux semaines.\nBudget augmenté de 15 % pour les mesures de sécurité.\nRéunions d'avancement hebdomadaires planifiées chaque mardi à 10h.",
            summary: "Réunion client pour discuter du calendrier du projet et des livrables. Plusieurs décisions clés ont été prises concernant le budget, le calendrier et la fréquence de communication.",
            sentiment: "positif",
            topics: JSON.stringify(["gestion de projet", "budget", "calendrier", "réunions"]),
            suggestedActions: "Sarah doit envoyer le plan de projet mis à jour d'ici vendredi.\nMichael doit fournir les spécifications techniques pour les nouvelles fonctionnalités.\nJohn doit confirmer la disponibilité des parties prenantes pour la session de démonstration.",
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
        title: metadata?.title || `Enregistrement ${now.toLocaleTimeString()}`,
        caseId: metadata?.caseId,
        date: now.toLocaleDateString("fr-FR", {
            day: "2-digit",
            month: "short",
            year: "numeric",
        }),
        startTime: now.toLocaleTimeString("fr-FR", {
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
        throw new Error("Enregistrement introuvable");
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
                    text: "Ceci est une transcription automatiquement générée pour votre enregistrement. L'audio a été traité et converti en texte. Veuillez relire et modifier si nécessaire avant de confirmer.\n\nNote : Il s'agit de données de test à des fins de développement.",
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
        throw new Error("Transcription introuvable. Le traitement n'est peut-être pas terminé.");
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
        throw new Error("Transcription introuvable");
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
        throw new Error("Transcription introuvable");
    }

    // Create mock analysis if it doesn't exist
    if (!mockAnalysis.has(id)) {
        mockAnalysis.set(id, {
            id: `mock-analysis-${id}`,
            transcriptionId: `mock-transcription-${id}`,
            keyFindings: "Observation clé extraite de la transcription.\nDétail important identifié dans l'enregistrement.",
            summary: "Résumé généré automatiquement du contenu de l'enregistrement. Veuillez relire et vérifier avant utilisation.",
            sentiment: "neutre",
            topics: JSON.stringify(["enregistrement", "notes", "révision"]),
            suggestedActions: "Relire la transcription pour vérifier l'exactitude.\nAssurer un suivi sur les sections peu claires.",
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
        throw new Error("Analyse introuvable. Veuillez d'abord lancer l'analyse.");
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
        throw new Error("Enregistrement introuvable");
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
        throw new Error("Enregistrement introuvable");
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
 * Mock: List all cases for the current user
 */
export async function getCases(): Promise<Case[]> {
    await delay(300);
    return [
        {
            id: "mock-case-1",
            title: "Réclamation pour préjudice corporel",
            status: "in_progress",
            externalReference: "CS-2024-892",
            clientId: "mock-client-1",
            clientName: "Sarah JENKINS",
            clientNumber: "CLI-9438",
        },
        {
            id: "mock-case-2",
            title: "Litige contrat de travail",
            status: "ready",
            externalReference: "CS-2024-741",
            clientId: "mock-client-2",
            clientName: "Marc DUPONT",
            clientNumber: "CLI-2201",
        },
        {
            id: "mock-case-3",
            title: "Divorce contentieux",
            status: "in_progress",
            externalReference: "CS-2025-013",
            clientId: "mock-client-3",
            clientName: "Julie MARTIN",
            clientNumber: "CLI-5567",
        },
    ];
}

/**
 * Mock: Get the current active case
 */
export async function getCurrentCase(caseId?: string): Promise<Case | null> {
    await delay(300);
    const cases = await getCases();
    return cases.find((c) => c.id === caseId) ?? cases[0] ?? null;
}

/**
 * Mock: Get audio URL for a recording
 */
export async function getAudioUrl(id: string): Promise<string> {
    await delay(200);

    const recording = mockRecordings.get(id);
    if (!recording) {
        throw new Error("Enregistrement introuvable");
    }

    return recording.audioUrl ?? "";
}

<script lang="ts">
    import { currentCase } from "$lib/stores/case";
    import { page } from "$app/stores";
    import { onMount } from "svelte";
    import {
        fetchCaseMedia,
        submitTranscriptionJob,
        getTranscriptionJobStatus,
        getTranscriptionByMediaFile,
    } from "$lib/services/api";
    import type { MediaFile, TranscriptionJob, Transcription } from "$lib/types/entities";
    import { getToastContext } from "$lib/custom/global/toast/state.svelte";
    import { TOAST_LEVELS } from "$lib/custom/global/toast/type";
    import Spinner from "$lib/components/Spinner.svelte";
    import { Mic, Play, Square, FileText, Loader2, AlertCircle, RefreshCw } from "@lucide/svelte";

    const toastState = getToastContext();

    interface RecordingState {
        media: MediaFile;
        transcriptionJob?: TranscriptionJob;
        transcription?: Transcription;
        loadingJob: boolean;
        loadingTranscription: boolean;
        isPlaying: boolean;
    }

    let recordings = $state<RecordingState[]>([]);
    let loading = $state(true);
    let currentAudio: HTMLAudioElement | null = null;
    let playingId: string | null = $state(null);

    $effect(() => {
        const caseId = $page.params.id;
        if (caseId && caseId !== $currentCase.id) {
            currentCase.setCase(caseId);
        }
    });

    onMount(loadRecordings);

    async function loadRecordings() {
        const caseId = $page.params.id;
        if (!caseId) return;
        loading = true;
        try {
            const response = await fetchCaseMedia(caseId, "audio");
            recordings = response.media.map((m) => ({
                media: m,
                loadingJob: false,
                loadingTranscription: false,
                isPlaying: false,
            }));
            // Load transcription status for each recording in the background
            for (const rec of recordings) {
                loadTranscriptionStatus(rec);
            }
        } catch (err) {
            toastState.add(TOAST_LEVELS.Error, "Erreur", "Impossible de charger les enregistrements.");
        } finally {
            loading = false;
        }
    }

    async function loadTranscriptionStatus(rec: RecordingState) {
        rec.loadingJob = true;
        try {
            const result = await getTranscriptionByMediaFile(rec.media.id);
            if (result.transcriptions.length > 0) {
                rec.transcription = result.transcriptions[0];
            }
        } catch {
            // No transcription yet — that's fine
        } finally {
            rec.loadingJob = false;
        }
    }

    async function handleTranscribe(rec: RecordingState) {
        rec.loadingJob = true;
        try {
            const job = await submitTranscriptionJob(rec.media.id);
            rec.transcriptionJob = job;
            toastState.add(TOAST_LEVELS.Info, "Transcription lancée", "La transcription est en cours de traitement.");

            // Poll for completion
            pollJob(rec, job.jobId);
        } catch (err) {
            toastState.add(TOAST_LEVELS.Error, "Erreur", "Impossible de lancer la transcription.");
            rec.loadingJob = false;
        }
    }

    async function pollJob(rec: RecordingState, jobId: string) {
        const maxAttempts = 60; // 5 minutes at 5s interval
        let attempts = 0;

        const interval = setInterval(async () => {
            attempts++;
            try {
                const job = await getTranscriptionJobStatus(jobId);
                rec.transcriptionJob = job;

                if (job.status === "completed") {
                    clearInterval(interval);
                    rec.loadingJob = false;
                    // Load the transcription
                    await loadTranscriptionStatus(rec);
                    toastState.add(TOAST_LEVELS.Info, "Transcription terminée", `"${rec.media.fileName}" a été transcrit.`);
                } else if (job.status === "failed" || job.status === "cancelled") {
                    clearInterval(interval);
                    rec.loadingJob = false;
                    toastState.add(TOAST_LEVELS.Error, "Erreur", `La transcription a échoué : ${job.errorMessage || "raison inconnue"}.`);
                }
            } catch {
                // Ignore poll errors
            }

            if (attempts >= maxAttempts) {
                clearInterval(interval);
                rec.loadingJob = false;
            }
        }, 5000);
    }

    function togglePlayback(rec: RecordingState) {
        if (playingId === rec.media.id && currentAudio) {
            // Stop current
            currentAudio.pause();
            currentAudio.currentTime = 0;
            currentAudio = null;
            playingId = null;
            rec.isPlaying = false;
            return;
        }

        // Stop any previous audio
        if (currentAudio) {
            currentAudio.pause();
            currentAudio.currentTime = 0;
            const prev = recordings.find((r) => r.media.id === playingId);
            if (prev) prev.isPlaying = false;
        }

        const audio = new Audio(rec.media.url);
        audio.onended = () => {
            playingId = null;
            rec.isPlaying = false;
            currentAudio = null;
        };
        audio.play();
        currentAudio = audio;
        playingId = rec.media.id;
        rec.isPlaying = true;
    }

    function formatDuration(ms: number): string {
        const seconds = Math.floor(ms / 1000);
        const m = Math.floor(seconds / 60);
        const s = seconds % 60;
        return `${m}:${s.toString().padStart(2, "0")}`;
    }

    function formatDate(dateStr: string): string {
        return new Date(dateStr).toLocaleDateString("fr-FR", {
            day: "2-digit",
            month: "short",
            year: "numeric",
            hour: "2-digit",
            minute: "2-digit",
        });
    }

    function formatFileSize(bytes: number): string {
        if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(0)} Ko`;
        return `${(bytes / (1024 * 1024)).toFixed(1)} Mo`;
    }

    function getJobStatusLabel(status: string): string {
        switch (status) {
            case "pending": return "En attente";
            case "processing": return "En cours";
            case "completed": return "Terminée";
            case "failed": return "Échouée";
            case "cancelled": return "Annulée";
            default: return status;
        }
    }

    function getJobStatusClass(status: string): string {
        switch (status) {
            case "completed": return "bg-success/15 text-success";
            case "failed":
            case "cancelled": return "bg-destructive/15 text-destructive";
            default: return "bg-muted text-muted-foreground";
        }
    }
</script>

<div class="p-8 flex flex-col flex-1 min-h-0 gap-6">
    <div>
        <h1 class="text-2xl font-bold text-foreground">Enregistrements</h1>
        <p class="text-sm text-muted-foreground">
            {recordings.length} enregistrement{recordings.length !== 1 ? "s" : ""} audio
        </p>
    </div>

    {#if loading}
        <div class="flex items-center justify-center py-12">
            <Spinner size="lg" />
        </div>
    {:else if recordings.length === 0}
        <div class="border border-dashed border-border rounded-lg bg-muted/20 flex flex-col items-center justify-center flex-1 gap-4 min-h-[50vh]">
            <Mic class="w-12 h-12 text-muted-foreground" />
            <p class="text-muted-foreground text-center">Aucun enregistrement pour ce dossier.</p>
            <p class="text-xs text-muted-foreground text-center max-w-sm">
                Les enregistrements sont capturés depuis l'application mobile lors des missions sur le terrain.
            </p>
        </div>
    {:else}
        <div class="flex-1 min-h-0 overflow-y-auto">
            <div class="flex flex-col gap-3">
                {#each recordings as rec, index (rec.media.id)}
                    <div
                        class="border border-border-card rounded-card p-5 bg-background hover:shadow-popover transition-shadow animate-fade-in"
                        style="animation-delay: {index * 50}ms;"
                    >
                        <div class="flex items-start gap-4">
                            <!-- Play button -->
                            <button
                                type="button"
                                onclick={() => togglePlayback(rec)}
                                class="w-12 h-12 rounded-full bg-foreground text-background flex items-center justify-center flex-shrink-0 hover:opacity-90 active:scale-95 transition-all shadow-mini cursor-pointer"
                                title={rec.isPlaying ? "Arrêter" : "Écouter"}
                            >
                                {#if rec.isPlaying}
                                    <Square size={18} fill="currentColor" />
                                {:else}
                                    <Play size={18} fill="currentColor" class="ml-0.5" />
                                {/if}
                            </button>

                            <!-- Info -->
                            <div class="flex-1 min-w-0">
                                <p class="font-medium text-foreground truncate">{rec.media.fileName}</p>
                                <div class="flex items-center gap-3 mt-1">
                                    <span class="text-xs text-muted-foreground">{formatFileSize(rec.media.fileSize)}</span>
                                    <span class="text-xs text-muted-foreground">{formatDate(rec.media.createdAt)}</span>
                                    {#if rec.media.caption}
                                        <span class="text-xs text-muted-foreground truncate">— {rec.media.caption}</span>
                                    {/if}
                                </div>

                                <!-- Transcription status & actions -->
                                <div class="mt-3 flex items-center gap-3">
                                    {#if rec.loadingJob}
                                        <span class="inline-flex items-center gap-1.5 text-xs text-muted-foreground">
                                            <Loader2 size={12} class="animate-spin" />
                                            {rec.transcriptionJob?.status === "processing"
                                                ? `Transcription en cours… (${rec.transcriptionJob.progress}%)`
                                                : "Chargement…"}
                                        </span>
                                    {:else if rec.transcription}
                                        <span class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium {getJobStatusClass('completed')}">
                                            <FileText size={10} />
                                            Transcription disponible
                                        </span>
                                        <span class="text-xs text-muted-foreground">
                                            {formatDuration(rec.transcription.duration)} · {rec.transcription.language.toUpperCase()}
                                        </span>
                                    {:else if rec.transcriptionJob}
                                        <span class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium {getJobStatusClass(rec.transcriptionJob.status)}">
                                            {rec.transcriptionJob.status === "failed"
                                                ? AlertCircle
                                                : Loader2
                                            }
                                            {getJobStatusLabel(rec.transcriptionJob.status)}
                                        </span>
                                        {#if rec.transcriptionJob.status === "failed"}
                                            <button
                                                type="button"
                                                onclick={() => handleTranscribe(rec)}
                                                class="inline-flex items-center gap-1 text-xs text-foreground hover:text-muted-foreground cursor-pointer"
                                            >
                                                <RefreshCw size={10} />
                                                Réessayer
                                            </button>
                                        {/if}
                                    {:else}
                                        <button
                                            type="button"
                                            onclick={() => handleTranscribe(rec)}
                                            class="inline-flex items-center gap-1.5 px-3 py-1 rounded-input border border-border-input text-xs font-medium text-foreground hover:bg-muted active:scale-[0.98] transition-all cursor-pointer"
                                        >
                                            <FileText size={12} />
                                            Lancer la transcription
                                        </button>
                                    {/if}
                                </div>

                                <!-- Transcription text -->
                                {#if rec.transcription?.transcript}
                                    <div class="mt-3 p-3 rounded-input bg-muted/50 border border-border-card">
                                        <p class="text-sm text-foreground whitespace-pre-wrap leading-relaxed">{rec.transcription.transcript}</p>
                                    </div>
                                {/if}
                            </div>
                        </div>
                    </div>
                {/each}
            </div>
        </div>
    {/if}
</div>

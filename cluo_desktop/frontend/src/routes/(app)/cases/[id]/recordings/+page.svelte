<script lang="ts">
    import { currentCase } from "$lib/stores/case";
    import { page } from "$app/stores";
    import { onMount } from "svelte";
    import {
        fetchCaseMedia,
        submitTranscriptionJob,
        getTranscriptionByMediaFile,
        analyzeTranscript,
        getAnalysisByTranscriptionId,
    } from "$lib/services/api";
    import { jobTracker } from "$lib/services/jobTracker";
    import { notificationStore } from "$lib/stores/notifications.svelte";
    import type { MediaFile, TranscriptionJob, Transcription, TranscriptAnalysis } from "$lib/types/entities";
    import { getToastContext } from "$lib/custom/global/toast/state.svelte";
    import { TOAST_LEVELS } from "$lib/custom/global/toast/type";
    import Spinner from "$lib/components/Spinner.svelte";
    import { Mic, Play, Square, FileText, Loader2, AlertCircle, RefreshCw, Sparkles, ChevronDown, ChevronRight, Tag } from "@lucide/svelte";

    const toastState = getToastContext();

    interface RecordingState {
        media: MediaFile;
        transcriptionJob?: TranscriptionJob;
        transcription?: Transcription;
        analysis?: TranscriptAnalysis | null;
        loadingJob: boolean;
        loadingTranscription: boolean;
        loadingAnalysis: boolean;
        analysisExpanded: boolean;
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
                analysis: null,
                loadingJob: false,
                loadingTranscription: false,
                loadingAnalysis: false,
                analysisExpanded: false,
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
                // Load any existing analysis for this transcription
                loadAnalysis(rec);
            }
        } catch {
            // No transcription yet — that's fine
        } finally {
            rec.loadingJob = false;
        }
    }

    async function loadAnalysis(rec: RecordingState) {
        if (!rec.transcription) return;
        rec.loadingAnalysis = true;
        try {
            rec.analysis = await getAnalysisByTranscriptionId(rec.transcription.id);
        } catch {
            rec.analysis = null;
        } finally {
            rec.loadingAnalysis = false;
        }
    }

    async function handleAnalyze(rec: RecordingState) {
        if (!rec.transcription) return;
        const caseId = $page.params.id;
        rec.loadingAnalysis = true;
        try {
            rec.analysis = await analyzeTranscript(rec.transcription.id);
            rec.analysisExpanded = true;
            toastState.add(TOAST_LEVELS.Info, "Analyse terminée", `« ${rec.media.fileName} » a été analysé.`);
            if (caseId) {
                notificationStore.push({
                    kind: "analysis_completed",
                    title: "Analyse prête",
                    content: `« ${rec.media.fileName} » a été analysé.`,
                    caseId,
                });
            }
        } catch (err) {
            toastState.add(TOAST_LEVELS.Error, "Erreur", "Impossible d'analyser la transcription.");
        } finally {
            rec.loadingAnalysis = false;
        }
    }

    function parseTopics(analysis: TranscriptAnalysis): string[] {
        try {
            const parsed = JSON.parse(analysis.topics);
            return Array.isArray(parsed) ? parsed.filter((t) => typeof t === "string") : [];
        } catch {
            return [];
        }
    }

    function sentimentLabel(sentiment: string): string {
        switch (sentiment) {
            case "positive": return "Positif";
            case "neutral": return "Neutre";
            case "negative": return "Négatif";
            case "mixed": return "Mixte";
            default: return sentiment;
        }
    }

    function sentimentClass(sentiment: string): string {
        switch (sentiment) {
            case "positive": return "bg-success/15 text-success";
            case "negative": return "bg-destructive/15 text-destructive";
            case "mixed": return "bg-primary/15 text-primary";
            default: return "bg-muted text-muted-foreground";
        }
    }

    async function handleTranscribe(rec: RecordingState) {
        const caseId = $page.params.id;
        rec.loadingJob = true;
        try {
            const job = await submitTranscriptionJob(rec.media.id);
            rec.transcriptionJob = job;
            toastState.add(TOAST_LEVELS.Info, "Transcription lancée", "La transcription est en cours de traitement.");

            // Track the job globally so its result is surfaced through the
            // notification bell even if the investigator navigates away from
            // this page. The page no longer polls on its own.
            jobTracker.trackJob(job.jobId, rec.media.fileName, caseId ?? "");
        } catch (err) {
            toastState.add(TOAST_LEVELS.Error, "Erreur", "Impossible de lancer la transcription.");
        } finally {
            rec.loadingJob = false;
        }
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

<div class="page-content flex-1 min-h-0">
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
                        class="border border-border-card rounded-card p-4 bg-background hover:shadow-card transition-shadow"
                    >
                        <div class="flex items-start gap-4">
                            <!-- Play button -->
                            <button
                                type="button"
                                onclick={() => togglePlayback(rec)}
                                class="w-12 h-12 rounded-full bg-foreground text-background flex items-center justify-center flex-shrink-0 hover:opacity-90 active:scale-95 transition-interactive shadow-mini cursor-pointer"
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
                                            class="inline-flex items-center gap-1.5 px-3 py-1 rounded-input border border-border-input text-xs font-medium text-foreground hover:bg-muted active:scale-[0.98] transition-interactive cursor-pointer"
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

                                <!-- Transcript analysis -->
                                {#if rec.transcription}
                                    <div class="mt-3">
                                        {#if rec.loadingAnalysis}
                                            <span class="inline-flex items-center gap-1.5 text-xs text-muted-foreground">
                                                <Loader2 size={12} class="animate-spin" />
                                                Analyse en cours…
                                            </span>
                                        {:else if rec.analysis}
                                            <div class="rounded-input border border-border-card bg-muted/20">
                                                <button
                                                    type="button"
                                                    onclick={() => (rec.analysisExpanded = !rec.analysisExpanded)}
                                                    class="w-full flex items-center justify-between gap-2 px-3 py-2 cursor-pointer hover:bg-muted/40 transition-interactive"
                                                >
                                                    <span class="inline-flex items-center gap-2">
                                                        <Sparkles size={14} class="text-primary" />
                                                        <span class="text-xs font-medium text-foreground">Analyse IA</span>
                                                        <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium {sentimentClass(rec.analysis.sentiment)}">
                                                            {sentimentLabel(rec.analysis.sentiment)}
                                                        </span>
                                                    </span>
                                                    {#if rec.analysisExpanded}
                                                        <ChevronDown size={14} class="text-muted-foreground" />
                                                    {:else}
                                                        <ChevronRight size={14} class="text-muted-foreground" />
                                                    {/if}
                                                </button>
                                                {#if rec.analysisExpanded}
                                                    <div class="px-3 pb-3 pt-3 flex flex-col gap-3 border-t border-border-card">
                                                        {#if rec.analysis.summary}
                                                            <div>
                                                                <p class="text-xs font-medium text-muted-foreground mb-1">Résumé</p>
                                                                <p class="text-sm text-foreground leading-relaxed">{rec.analysis.summary}</p>
                                                            </div>
                                                        {/if}
                                                        {#if rec.analysis.keyFindings}
                                                            <div>
                                                                <p class="text-xs font-medium text-muted-foreground mb-1">Points clés</p>
                                                                <p class="text-sm text-foreground whitespace-pre-wrap leading-relaxed">{rec.analysis.keyFindings}</p>
                                                            </div>
                                                        {/if}
                                                        {#if parseTopics(rec.analysis).length > 0}
                                                            <div>
                                                                <p class="text-xs font-medium text-muted-foreground mb-1">Sujets</p>
                                                                <div class="flex flex-wrap gap-1.5">
                                                                    {#each parseTopics(rec.analysis) as topic}
                                                                        <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs bg-muted text-muted-foreground border border-border-card">
                                                                            <Tag size={10} />
                                                                            {topic}
                                                                        </span>
                                                                    {/each}
                                                                </div>
                                                            </div>
                                                        {/if}
                                                        {#if rec.analysis.suggestedActions}
                                                            <div>
                                                                <p class="text-xs font-medium text-muted-foreground mb-1">Actions suggérées</p>
                                                                <p class="text-sm text-foreground whitespace-pre-wrap leading-relaxed">{rec.analysis.suggestedActions}</p>
                                                            </div>
                                                        {/if}
                                                        <button
                                                            type="button"
                                                            onclick={() => handleAnalyze(rec)}
                                                            class="inline-flex items-center gap-1 text-xs text-foreground hover:text-muted-foreground cursor-pointer self-start"
                                                        >
                                                            <RefreshCw size={10} />
                                                            Ré-analyser
                                                        </button>
                                                    </div>
                                                {/if}
                                            </div>
                                        {:else}
                                            <button
                                                type="button"
                                                onclick={() => handleAnalyze(rec)}
                                                class="inline-flex items-center gap-1.5 px-3 py-1 rounded-input border border-border-input text-xs font-medium text-foreground hover:bg-muted active:scale-[0.98] transition-interactive cursor-pointer"
                                            >
                                                <Sparkles size={12} />
                                                Analyser
                                            </button>
                                        {/if}
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

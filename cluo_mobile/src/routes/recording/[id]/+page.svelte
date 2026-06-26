<script lang="ts">
    import { onMount } from "svelte";
    import { ArrowLeft, FileText, Sparkles, Play, Trash2, Download, Loader2 } from "@lucide/svelte";
    import { goto } from "$app/navigation";
    import AudioPlayer from "$lib/components/AudioPlayer.svelte";
    import Spinner from "$lib/components/ui/Spinner.svelte";
    import { deleteRecording, getAudioUrl, getRecording } from "$lib/api";
    import { snackbar } from "$lib/stores/snackbar";
    import type { Recording, Transcript, AnalysisResult } from "$lib/types/recording";

    let { data } = $props();

    let recording = $state<(Recording & { audioUrl?: string }) | null>(null);
    let transcript = $state<Transcript | null>(null);
    let analysis = $state<AnalysisResult | null>(null);
    let isLoading = $state(true);
    let pageError = $state<string | null>(null);
    let isDeleting = $state(false);
    let isDownloading = $state(false);

    onMount(async () => {
        try {
            const result = await getRecording(data.id);
            recording = result.recording;
            transcript = result.transcript;
            analysis = result.analysis;
        } catch (e) {
            pageError = e instanceof Error ? e.message : "Échec du chargement de l'enregistrement";
        } finally {
            isLoading = false;
        }
    });

    function goBack() {
        if (history.length > 0) history.back();
        else window.location.href = "/";
    }

    const statusLabels: Record<string, string> = {
        uploading: "Téléchargement",
        transcribing: "Transcription",
        analyzing: "Analyse",
        completed: "Terminé",
        failed: "Échoué",
    };

    async function handleDownload() {
        if (isDownloading || !recording) return;

        try {
            isDownloading = true;
            const audioUrl = await getAudioUrl(recording.id);

            if (!audioUrl) {
                throw new Error("Audio not available");
            }

            const a = document.createElement("a");
            a.href = audioUrl;
            a.download = `${recording.title || "recording"}.webm`;
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
        } catch (error) {
            console.error("Download failed:", error);
            snackbar.error("Échec du téléchargement de l'audio", () => handleDownload());
        } finally {
            isDownloading = false;
        }
    }

    async function handleDelete() {
        if (isDeleting || !recording) return;

        const confirmed = confirm(
            `Êtes-vous sûr de vouloir supprimer « ${recording.title} » ? Cette action est irréversible.`
        );

        if (!confirmed) return;

        try {
            isDeleting = true;
            await deleteRecording(recording.id);
            goto("/");
        } catch (error) {
            console.error("Delete failed:", error);
            snackbar.error("Échec de la suppression de l'enregistrement", () => handleDelete());
        } finally {
            isDeleting = false;
        }
    }
</script>

<div class="min-h-screen flex flex-col gap-6 pb-24">
    <!-- Header with back button -->
    <div class="flex items-center gap-3">
        <button
            onclick={goBack}
            class="flex items-center justify-center w-10 h-10 rounded-full hover:bg-dark-50 transition-colors"
        >
            <ArrowLeft class="text-dark-700" />
        </button>
        <h1 class="text-dark-900 font-extrabold text-xl">Détails de l'enregistrement</h1>
    </div>

    {#if isLoading}
        <div class="flex items-center justify-center p-12">
            <Spinner size="md" />
        </div>
    {:else if pageError}
        <div class="flex items-center justify-center p-8 bg-red-50 rounded-2xl">
            <p class="text-red-600">{pageError}</p>
        </div>
    {:else if recording}
        <!-- Recording Info Card -->
        <div class="flex flex-col gap-4 p-4 border border-dark-100 rounded-2xl">
            <div class="flex justify-between items-start">
                <div class="min-w-0 flex-1">
                    <h2 class="text-dark-800 font-bold text-lg truncate">
                        {recording.title}
                    </h2>
                    <div class="flex gap-2 items-center mt-1">
                        <p class="text-dark-600 text-sm">{recording.date}</p>
                        <span class="text-dark-300">•</span>
                        <p class="text-dark-400 text-sm">{recording.startTime}</p>
                    </div>
                </div>
                <div class="flex items-center gap-2">
                    <p
                        class="flex justify-center items-center border border-dark-100 rounded-3xl bg-dark-50 text-dark-600 py-1 px-3 text-sm font-medium"
                    >
                        {recording.duration}
                    </p>
                </div>
            </div>

            {#if recording.audioUrl}
                <AudioPlayer src={recording.audioUrl} />
            {:else}
                <div
                    class="flex items-center gap-4 p-4 bg-dark-50 rounded-xl border border-dark-100"
                >
                    <button
                        class="flex items-center justify-center w-12 h-12 bg-dark-700 rounded-full hover:bg-dark-600 transition-colors"
                        disabled
                    >
                        <Play class="text-dark-400" fill="currentColor" size={20} />
                    </button>
                    <p class="text-dark-500 text-sm">Audio non disponible</p>
                </div>
            {/if}
        </div>

        <!-- Status Badge -->
        <div class="flex items-center justify-center py-2 px-4 rounded-full bg-dark-100 w-fit">
            <p class="text-dark-700 text-sm font-medium">
                {statusLabels[recording.status] ?? recording.status}
            </p>
        </div>

        <!-- Transcript Section Link -->
        {#if transcript}
            <a
                href="/recording/{recording.id}/transcript"
                class="flex items-center gap-3 p-4 bg-dark-50 hover:bg-dark-100 rounded-2xl transition-colors no-underline"
            >
                <div class="flex items-center justify-center w-10 h-10 bg-dark-200 rounded-full">
                    <FileText class="text-dark-700" size={20} />
                </div>
                <div class="flex-1">
                    <p class="text-dark-900 font-semibold">Transcription</p>
                    <p class="text-dark-600 text-sm">
                        {transcript.isConfirmed ? "Confirmée" : "En attente de révision"}
                    </p>
                </div>
            </a>
        {/if}

        <!-- Analysis Section Link -->
        {#if analysis}
            <a
                href="/recording/{recording.id}/analysis"
                class="flex items-center gap-3 p-4 bg-dark-50 hover:bg-dark-100 rounded-2xl transition-colors no-underline"
            >
                <div class="flex items-center justify-center w-10 h-10 bg-dark-200 rounded-full">
                    <Sparkles class="text-dark-700" size={20} />
                </div>
                <div class="flex-1">
                    <p class="text-dark-900 font-semibold">Analyse IA</p>
                    <p class="text-dark-600 text-sm">
                        {analysis.sentiment === "neutre" ? "Neutre" : analysis.sentiment === "positif" ? "Positif" : analysis.sentiment === "négatif" ? "Négatif" : analysis.sentiment || "Terminé"}
                    </p>
                </div>
            </a>
        {/if}

        <!-- Quick Actions -->
        <div class="flex flex-col gap-3 mt-4">
            {#if !transcript}
                <a
                    href="/processing/{recording.id}"
                    class="flex items-center justify-center w-full py-4 bg-dark-700 hover:bg-dark-600 text-white rounded-xl transition-colors font-semibold no-underline"
                >
                    Voir l'état du traitement
                </a>
            {:else if !analysis}
                <a
                    href="/recording/{recording.id}/transcript"
                    class="flex items-center justify-center w-full py-4 bg-dark-700 hover:bg-dark-600 text-white rounded-xl transition-colors font-semibold no-underline"
                >
                    Réviser la transcription
                </a>
            {/if}

            <div class="flex gap-3">
                <button
                    onclick={handleDownload}
                    disabled={isDownloading}
                    class="flex-1 flex items-center justify-center gap-2 px-4 py-3 bg-dark-700 hover:bg-dark-600 disabled:bg-dark-400 text-white rounded-xl transition-colors"
                >
                    {#if isDownloading}
                        <Loader2 size={18} class="animate-spin" />
                    {:else}
                        <Download size={18} />
                    {/if}
                    <span class="font-medium">Télécharger</span>
                </button>
                <button
                    onclick={handleDelete}
                    disabled={isDeleting}
                    class="flex-1 flex items-center justify-center gap-2 px-4 py-3 bg-destructive hover:bg-destructive/90 disabled:bg-destructive/50 text-white rounded-xl transition-colors"
                >
                    {#if isDeleting}
                        <Loader2 size={18} class="animate-spin" />
                    {:else}
                        <Trash2 size={18} />
                    {/if}
                    <span class="font-medium">Supprimer</span>
                </button>
            </div>
        </div>
    {/if}
</div>

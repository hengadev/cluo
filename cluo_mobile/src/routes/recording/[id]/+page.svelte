<script lang="ts">
    import { onMount } from "svelte";
    import { ArrowLeft, FileText, Sparkles, Play, Trash2, Download, Loader2, Pencil } from "@lucide/svelte";
    import { goto } from "$app/navigation";
    import AudioPlayer from "$lib/components/AudioPlayer.svelte";
    import Spinner from "$lib/components/ui/Spinner.svelte";
    import { deleteRecording, getAudioUrl, getRecording, updateRecording } from "$lib/api";
    import { snackbar } from "$lib/stores/snackbar";
    import type { Recording, RecordingPurpose, Transcript, AnalysisResult } from "$lib/types/recording";

    let { data } = $props();

    let recording = $state<(Recording & { audioUrl?: string }) | null>(null);
    let transcript = $state<Transcript | null>(null);
    let analysis = $state<AnalysisResult | null>(null);
    let isLoading = $state(true);
    let pageError = $state<string | null>(null);
    let isDeleting = $state(false);
    let isDownloading = $state(false);
    let purposeDialogOpen = $state(false);
    let pendingPurpose = $state<RecordingPurpose>("general");
    let isUpdatingPurpose = $state(false);
    let titleDialogOpen = $state(false);
    let pendingTitle = $state("");
    let isUpdatingTitle = $state(false);

    const purposeLabels: Record<RecordingPurpose, string> = {
        general: "Général",
        witness_interview: "Audition témoin",
    };

    function openPurposeDialog() {
        if (!recording) return;
        pendingPurpose = recording.purpose;
        purposeDialogOpen = true;
    }

    function openTitleDialog() {
        if (!recording) return;
        pendingTitle = recording.title;
        titleDialogOpen = true;
    }

    async function confirmTitleChange() {
        if (!recording || isUpdatingTitle) return;
        const trimmed = pendingTitle.trim();
        if (!trimmed || trimmed === recording.title) {
            titleDialogOpen = false;
            return;
        }
        try {
            isUpdatingTitle = true;
            await updateRecording(recording.id, { title: trimmed });
            recording.title = trimmed;
            titleDialogOpen = false;
        } catch {
            snackbar.error("Échec de la mise à jour du titre");
        } finally {
            isUpdatingTitle = false;
        }
    }

    async function confirmPurposeChange() {
        if (!recording || isUpdatingPurpose) return;
        try {
            isUpdatingPurpose = true;
            await updateRecording(recording.id, { purpose: pendingPurpose });
            recording.purpose = pendingPurpose;
            purposeDialogOpen = false;
        } catch (error) {
            snackbar.error("Échec de la mise à jour de la finalité");
        } finally {
            isUpdatingPurpose = false;
        }
    }

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

    function formatDuration(d: number | string): string {
        const secs = typeof d === "number" ? d : 0;
        if (secs <= 0) return "--:--";
        const m = Math.floor(secs / 60);
        const s = Math.floor(secs % 60);
        return `${m.toString().padStart(2, "0")}:${s.toString().padStart(2, "0")}`;
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

<div class="min-h-screen flex flex-col gap-6 pb-6">
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
                    <div class="flex items-center gap-2">
                        <h2 class="text-dark-800 font-bold text-lg truncate">
                            {recording.title}
                        </h2>
                        <button
                            onclick={openTitleDialog}
                            class="flex-shrink-0 text-dark-400 hover:text-dark-700 transition-colors"
                            aria-label="Renommer l'enregistrement"
                        >
                            <Pencil size={14} />
                        </button>
                    </div>
                    <div class="flex gap-2 items-center mt-1">
                        <p class="text-dark-600 text-sm">{recording.date}</p>
                        <span class="text-dark-300">•</span>
                        <p class="text-dark-400 text-sm">{recording.startTime}</p>
                    </div>
                </div>
                <div class="flex items-center gap-2">
                    <button
                        onclick={openPurposeDialog}
                        class="flex justify-center items-center border border-dark-100 rounded-3xl bg-dark-50 hover:bg-dark-100 text-dark-600 py-1 px-3 text-sm font-medium transition-colors"
                    >
                        {purposeLabels[recording.purpose]}
                    </button>
                    <p
                        class="flex justify-center items-center border border-dark-100 rounded-3xl bg-dark-50 text-dark-600 py-1 px-3 text-sm font-medium"
                    >
                        {formatDuration(recording.duration)}
                    </p>
                </div>
            </div>

            {#if recording.audioUrl}
                <AudioPlayer src={recording.audioUrl} duration={typeof recording.duration === 'number' ? recording.duration : 0} />
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

{#if titleDialogOpen}
    <div
        class="fixed inset-0 z-50 flex items-end justify-center bg-black/60"
        role="dialog"
        aria-modal="true"
        aria-label="Renommer l'enregistrement"
    >
        <div class="bg-background rounded-t-2xl p-6 w-full max-w-lg shadow-popover">
            <p class="text-dark-900 font-semibold text-base mb-4">Renommer l'enregistrement</p>
            <input
                type="text"
                bind:value={pendingTitle}
                placeholder={recording?.title ?? ""}
                class="w-full bg-dark-50 border border-dark-200 text-dark-800 placeholder-dark-400 px-3 py-3 rounded-xl text-base focus:outline-none focus:ring-1 focus:ring-dark-400 mb-6"
                onkeydown={(e) => e.key === "Enter" && confirmTitleChange()}
            />
            <div class="flex gap-3">
                <button
                    class="flex-1 px-4 py-3 rounded-xl border border-dark-200 text-dark-700 hover:bg-dark-50 transition-colors"
                    onclick={() => (titleDialogOpen = false)}
                >
                    Annuler
                </button>
                <button
                    class="flex-1 px-4 py-3 rounded-xl bg-dark-700 hover:bg-dark-600 text-white font-medium transition-colors disabled:opacity-50"
                    onclick={confirmTitleChange}
                    disabled={isUpdatingTitle || !pendingTitle.trim() || pendingTitle.trim() === recording?.title}
                >
                    {isUpdatingTitle ? "Enregistrement…" : "Confirmer"}
                </button>
            </div>
        </div>
    </div>
{/if}

{#if purposeDialogOpen}
    <div
        class="fixed inset-0 z-50 flex items-end justify-center bg-black/60"
        role="dialog"
        aria-modal="true"
        aria-label="Finalité de l'enregistrement"
    >
        <div class="bg-background rounded-t-2xl p-6 w-full max-w-lg shadow-popover">
            <p class="text-dark-900 font-semibold text-base mb-4">Finalité de l'enregistrement</p>
            <div class="flex flex-col gap-2 mb-6">
                {#each ([{ value: "general" as RecordingPurpose, label: "Général" }, { value: "witness_interview" as RecordingPurpose, label: "Audition témoin" }]) as option}
                    <button
                        onclick={() => (pendingPurpose = option.value)}
                        class="flex items-center gap-3 p-4 rounded-xl border transition-colors {pendingPurpose === option.value ? 'border-dark-700 bg-dark-50' : 'border-dark-100 hover:bg-dark-50'}"
                    >
                        <span class="w-4 h-4 rounded-full border-2 flex-shrink-0 {pendingPurpose === option.value ? 'border-dark-700 bg-dark-700' : 'border-dark-300'}"></span>
                        <span class="text-dark-900 font-medium">{option.label}</span>
                    </button>
                {/each}
            </div>
            <div class="flex gap-3">
                <button
                    class="flex-1 px-4 py-3 rounded-xl border border-dark-200 text-dark-700 hover:bg-dark-50 transition-colors"
                    onclick={() => (purposeDialogOpen = false)}
                >
                    Annuler
                </button>
                <button
                    class="flex-1 px-4 py-3 rounded-xl bg-dark-700 hover:bg-dark-600 text-white font-medium transition-colors disabled:opacity-50"
                    onclick={confirmPurposeChange}
                    disabled={isUpdatingPurpose || pendingPurpose === recording?.purpose}
                >
                    {isUpdatingPurpose ? "Enregistrement…" : "Confirmer"}
                </button>
            </div>
        </div>
    </div>
{/if}

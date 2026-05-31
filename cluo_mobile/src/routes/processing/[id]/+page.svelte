<script lang="ts">
    import { Check, ChevronLeft, Ellipsis, XCircle, RotateCw } from "@lucide/svelte";
    import Spinner from "$lib/components/ui/Spinner.svelte";
    import { goto } from "$app/navigation";
    import { onMount, onDestroy } from "svelte";
    import type { ProcessingStep } from "$lib/types/recording";
    import { getRecordingStatus } from "$lib/api";

    // Data is passed from +page.ts load function
    let { data } = $props();
    let recordingId = $derived(data.recordingId);
    let steps = $state<ProcessingStep[]>(data.steps);
    let error = $state(data.error);
    let isRetrying = $state(false);

    let pollInterval: number | null = null;
    const POLL_INTERVAL = 2000; // Poll every 2 seconds

    function goBack() {
        if (history.length > 0) history.back();
        else goto("/");
    }

    async function fetchStatus() {
        try {
            const status = await getRecordingStatus(recordingId);
            steps = status.processingSteps;

            // If processing is complete, stop polling and navigate to transcript
            if (status.status === "completed") {
                stopPolling();
                // Navigate to transcript page after a short delay
                setTimeout(() => {
                    goto(`/recording/${recordingId}/transcript`);
                }, 500);
            } else if (status.status === "failed") {
                stopPolling();
                error = status.error ?? "Processing failed";
            }
        } catch (err) {
            console.error("Failed to fetch status:", err);
            // Don't stop polling on transient errors
        }
    }

    function startPolling() {
        if (pollInterval !== null) return;
        pollInterval = window.setInterval(fetchStatus, POLL_INTERVAL);
    }

    function stopPolling() {
        if (pollInterval !== null) {
            clearInterval(pollInterval);
            pollInterval = null;
        }
    }

    async function retry() {
        isRetrying = true;
        error = null;
        await fetchStatus();
        isRetrying = false;
    }

    // Start polling on mount
    onMount(() => {
        // Only start polling if there's no initial error
        if (!error) {
            startPolling();
        }
    });

    // Clean up on destroy
    onDestroy(() => {
        stopPolling();
    });

    // Derived state for completion
    const isCompleted = $derived(steps.every((s) => s.status === "completed"));
    const hasFailed = $derived(error !== null);
</script>

<div class="min-h-screen flex flex-col gap-8 pb-24 mt-8 px-4">
    <div class="flex flex-col gap-2">
        <div class="flex items-center justify-between mb-8">
            <button onclick={goBack}>
                <ChevronLeft />
            </button>
            <button>
                <Ellipsis />
            </button>
        </div>
        <h1 class="text-dark-900 font-extrabold text-2xl">
            Traitement de l'enregistrement
        </h1>
        <p class="text-dark-600 text-base">
            Veuillez patienter pendant le traitement de votre enregistrement...
        </p>
    </div>

    {#if hasFailed}
        <!-- Error State -->
        <div class="flex flex-col items-center gap-4 p-8 bg-red-50 rounded-2xl">
            <XCircle class="text-red-500" size={48} />
            <div class="text-center">
                <p class="text-red-700 font-semibold text-lg">Échec du traitement</p>
                <p class="text-red-600 text-sm mt-1">{error}</p>
            </div>
            <button
                onclick={retry}
                class="flex items-center gap-2 px-6 py-3 bg-red-600 hover:bg-red-500 text-white rounded-xl transition-colors font-semibold"
                disabled={isRetrying}
            >
                <RotateCw size={18} class={isRetrying ? 'animate-spin' : ''} />
                <span>Réessayer</span>
            </button>
        </div>
    {:else}
        <!-- Processing Steps -->
        <div class="flex flex-col gap-4">
            {#each steps as step}
                <div
                    class="flex items-center gap-4 p-4 border-1 border-dark-100 rounded-2xl bg-background-alt"
                >
                    <div
                        class="flex items-center justify-center w-12 h-12 rounded-full {step.status ===
                        'completed'
                            ? 'bg-green-500'
                            : step.status === 'failed'
                                ? 'bg-red-500'
                                : 'bg-dark-50'}"
                    >
                        {#if step.status === "completed"}
                            <Check class="text-white" size={24} strokeWidth={3} />
                        {:else if step.status === "failed"}
                            <XCircle class="text-white" size={24} strokeWidth={3} />
                        {:else}
                            <Spinner size="md" />
                        {/if}
                    </div>

                    <div class="flex-1">
                        <p
                            class="text-dark-800 font-semibold text-base {step.status ===
                            'completed'
                                ? 'line-through text-dark-500'
                                : ''}"
                        >
                            {step.title}
                        </p>
                        {#if step.status === "completed"}
                            <p class="text-green-600 text-sm">Terminé</p>
                        {:else if step.status === "failed"}
                            <p class="text-red-600 text-sm">Échoué</p>
                        {:else if step.status === "processing"}
                            <p class="text-dark-500 text-sm">En cours...</p>
                        {:else}
                            <p class="text-dark-400 text-sm">En attente</p>
                        {/if}
                    </div>
                </div>
            {/each}
        </div>

        {#if isCompleted}
            <div class="mt-4">
                <a
                    href="/recording/{recordingId}/transcript"
                    class="flex items-center justify-center w-full px-6 py-4 bg-dark-700 hover:bg-dark-600 text-foreground rounded-xl transition-colors font-semibold no-underline"
                >
                    Voir la transcription
                </a>
            </div>
        {/if}
    {/if}

    <!-- Privacy Notice -->
    <div class="flex items-center justify-center p-4 bg-dark-50 rounded-2xl mt-4">
        <p class="text-dark-600 text-sm text-center">
            La transcription et l'analyse sont traitées sur une infrastructure privée
        </p>
    </div>
</div>

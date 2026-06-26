<script lang="ts">
    import { onMount } from "svelte";
    import { ChevronDown } from "@lucide/svelte";
    import { auth } from "$lib/stores/auth";

    import Input from "$lib/components/ui/Input.svelte";
    import Spinner from "$lib/components/ui/Spinner.svelte";
    import PastRecording from "./PastRecording.svelte";
    import CurrentCase from "./CurrentCase.svelte";
    import CasePicker from "./CasePicker.svelte";
    import { currentCase as currentCaseStore } from "$lib/stores/current-case";
    import { listRecordings, uploadRecording } from "$lib/api";
    import { flush } from "$lib/upload-queue";
    import { queueCount } from "$lib/stores/upload-queue-count";

    import type { Case } from "$lib/types/case";
    import type { Recording } from "$lib/types/recording";

    let { data } = $props();

    const PAGE_SIZE = 20;

    let recordings = $state<Recording[]>(data.recordings);
    let totalCount = $state<number>(data.totalCount);
    let error = $state<string | null>(data.error);
    let currentCase = $state<Case | null>(data.currentCase);
    let pickerOpen = $state(false);
    let recordingsLoading = $state(false);
    let loadingMore = $state(false);
    let loadMoreError = $state<string | null>(null);
    let searchQuery = $state("");
    let fetchSeq = 0;
    const remainingCount = $derived(Math.max(0, totalCount - recordings.length));
    const hasMore = $derived(remainingCount > 0);

    const filteredRecordings = $derived(
        searchQuery.trim() === ""
            ? recordings
            : recordings.filter((r) =>
                  r.title.toLowerCase().includes(searchQuery.trim().toLowerCase()),
              ),
    );

    // Sync local state to shared store so the layout/Footer can read it
    $effect(() => {
        currentCaseStore.set(currentCase);
    });

    const greeting = $derived($auth.user?.name ?? $auth.user?.email?.split('@')[0] ?? '');

    async function handleOnline() {
        const result = await flush(uploadRecording);
        if (result.succeeded.length > 0 && currentCase) {
            fetchRecordings(currentCase.id);
        }
        queueCount.refresh();
    }

    onMount(() => {
        queueCount.refresh();
        window.addEventListener("online", handleOnline);

        // Flush any queued uploads immediately if already online (not just on reconnect).
        if (navigator.onLine) {
            handleOnline();
        }

        // SSR runs without localStorage, so currentCase may be null on first paint.
        // Restore the last selected case from the cached object — no API call needed.
        if (!currentCase) {
            try {
                const raw = localStorage.getItem("cluo_current_case");
                if (raw) {
                    const restored = JSON.parse(raw) as Case;
                    currentCase = restored;
                    fetchRecordings(restored.id);
                }
            } catch {
                // localStorage unavailable or JSON malformed — remain empty
            }
        }

        return () => window.removeEventListener("online", handleOnline);
    });

    async function fetchRecordings(caseId: string) {
        const seq = ++fetchSeq;
        recordingsLoading = true;
        error = null;
        try {
            const res = await listRecordings({ caseId });
            if (seq !== fetchSeq) return;
            recordings = res.recordings;
            totalCount = res.totalCount;
        } catch (e) {
            if (seq !== fetchSeq) return;
            error = e instanceof Error ? e.message : "Échec du chargement des enregistrements";
        } finally {
            if (seq === fetchSeq) recordingsLoading = false;
        }
    }

    async function loadMore() {
        if (loadingMore || !hasMore || !currentCase) return;
        loadingMore = true;
        loadMoreError = null;
        try {
            const res = await listRecordings({
                caseId: currentCase.id,
                offset: recordings.length,
                limit: PAGE_SIZE,
            });
            recordings = [...recordings, ...res.recordings];
            totalCount = res.totalCount;
        } catch (e) {
            loadMoreError = e instanceof Error ? e.message : "Échec du chargement des enregistrements";
        } finally {
            loadingMore = false;
        }
    }

    function handleCaseSelect(c: Case) {
        currentCase = c;
        searchQuery = "";
        try {
            localStorage.setItem("cluo_current_case_id", c.id);
            localStorage.setItem("cluo_current_case", JSON.stringify(c));
        } catch {
            // localStorage unavailable
        }
        fetchRecordings(c.id);
    }
</script>

<div class="min-h-screen flex flex-col gap-8 pb-28">
    <p class="text-dark-900 font-extrabold text-xl">Bonjour {greeting},</p>
    <div class="grid gap-4">
        <div class="flex justify-between items-center">
            <p class="font-extrabold text-lg text-dark-800">Affaire active</p>
            <button
                onclick={() => (pickerOpen = true)}
                class="flex items-center gap-1 text-dark-600 text-sm cursor-pointer hover:text-dark-900 transition-colors"
            >
                <span>Changer d'affaire</span>
                <ChevronDown size={16} />
            </button>
        </div>
        <CurrentCase {currentCase} />
    </div>
    <Input placeholder="Recherche parmi les enregistrements" bind:value={searchQuery} type="search" />
    <div class="flex flex-col gap-4">
        <div class="flex items-center gap-2">
            <p class="text-dark-700 font-bold text-base">Enregistrements</p>
            {#if $queueCount > 0}
                <span class="bg-tertiary text-background text-xs font-semibold px-2 py-0.5 rounded-full">
                    {$queueCount} en attente
                </span>
            {/if}
        </div>
        {#if recordingsLoading}
            <div class="flex items-center justify-center p-8">
                <Spinner size="md" />
            </div>
        {:else if error}
            <div class="flex items-center justify-center p-4 bg-red-50 rounded-2xl">
                <p class="text-red-600 text-sm">{error}</p>
            </div>
        {:else if recordings.length === 0}
            <div class="flex items-center justify-center p-8 bg-dark-50 rounded-2xl">
                <p class="text-dark-600">Aucun enregistrement pour le moment. Commencez par enregistrer des notes !</p>
            </div>
        {:else if filteredRecordings.length === 0}
            <div class="flex items-center justify-center p-8 bg-dark-50 rounded-2xl">
                <p class="text-dark-600">Aucun enregistrement ne correspond à «\u00a0{searchQuery}\u00a0».</p>
            </div>
        {:else}
            <div class="flex flex-col gap-2">
                {#each filteredRecordings as recording}
                    <PastRecording
                        id={recording.id}
                        title={recording.title}
                        date={recording.date}
                        startTime={recording.startTime}
                        duration={recording.duration}
                        status={recording.status}
                    />
                {/each}

                {#if loadMoreError}
                    <p class="text-red-600 text-sm text-center py-2">{loadMoreError}</p>
                {/if}
                {#if hasMore}
                    <button
                        onclick={loadMore}
                        disabled={loadingMore}
                        class="text-sm text-dark-500 hover:text-dark-900 transition-colors py-2 cursor-pointer disabled:cursor-not-allowed disabled:opacity-60 flex items-center justify-center gap-2"
                    >
                        {#if loadingMore}
                            <Spinner size="sm" />
                            <span>Chargement...</span>
                        {:else}
                            <span>Voir plus ({remainingCount})</span>
                        {/if}
                    </button>
                {/if}
            </div>
        {/if}
    </div>
</div>

<CasePicker
    bind:open={pickerOpen}
    activeId={currentCase?.id ?? null}
    onselect={handleCaseSelect}
/>

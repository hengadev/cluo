<script lang="ts">
    import { ChevronDown } from "@lucide/svelte";
    import { auth } from "$lib/stores/auth";

    import Input from "$lib/components/ui/Input.svelte";
    import Recording from "./PastRecording.svelte";
    import CurrentCase from "./CurrentCase.svelte";
    import CasePicker from "./CasePicker.svelte";

    import type { Case } from "$lib/types/case";

    let { data } = $props();
    const recordings = data.recordings;
    const error = data.error;

    let currentCase = $state<Case | null>(data.currentCase);
    let pickerOpen = $state(false);

    const greeting = $derived($auth.user?.name ?? $auth.user?.email?.split('@')[0] ?? '');

    function handleCaseSelect(c: Case) {
        currentCase = c;
        try {
            localStorage.setItem("cluo_current_case_id", c.id);
        } catch {
            // localStorage unavailable
        }
    }
</script>

<div class="min-h-screen flex flex-col gap-8">
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
    <div class="flex gap-4">
        <Input placeholder="Recherche parmi les enregistrements" />
        <button class="text-dark-500">Modifier</button>
    </div>
    <div class="flex flex-col gap-4">
        <p class="text-dark-700 font-bold text-base">Enregistrements</p>
        {#if error}
            <div class="flex items-center justify-center p-4 bg-red-50 rounded-2xl">
                <p class="text-red-600 text-sm">{error}</p>
            </div>
        {:else if recordings.length === 0}
            <div class="flex items-center justify-center p-8 bg-dark-50 rounded-2xl">
                <p class="text-dark-600">Aucun enregistrement pour le moment. Commencez par enregistrer des notes !</p>
            </div>
        {:else}
            <div class="flex flex-col gap-2">
                {#each recordings as recording}
                    <Recording
                        id={recording.id}
                        title={recording.title}
                        date={recording.date}
                        startTime={recording.startTime}
                        duration={recording.duration}
                        status={recording.status}
                    />
                {/each}
            </div>
        {/if}
    </div>
</div>

<CasePicker
    bind:open={pickerOpen}
    activeId={currentCase?.id ?? null}
    onselect={handleCaseSelect}
/>

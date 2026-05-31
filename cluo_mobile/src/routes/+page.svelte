<script lang="ts">
    import { ChevronDown } from "@lucide/svelte";

    import Input from "$lib/components/ui/Input.svelte";
    import Recording from "./PastRecording.svelte";
    import CurrentCase from "./CurrentCase.svelte";

    // Data is passed from +page.ts load function
    let { data } = $props();
    const recordings = data.recordings;
    const error = data.error;
</script>

<div class="min-h-screen flex flex-col gap-8">
    <p class="text-dark-900 font-extrabold text-xl">Bonjour John,</p>
    <div class="grid gap-4">
        <div class="flex justify-between items-center">
            <p class="font-extrabold text-lg text-dark-800">Affaire active</p>
            <div class="flex items-center gap-1">
                <p class="text-dark-600 text-sm">Changer d'affaire</p>
                <ChevronDown />
            </div>
        </div>
        <CurrentCase />
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

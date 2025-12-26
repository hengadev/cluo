<script lang="ts">
    import { ChevronDown } from "@lucide/svelte";

    import Input from "$lib/components/ui/Input.svelte";
    import Recording from "./PastRecording.svelte";
    import CurrentCase from "./CurrentCase.svelte";
</script>

<div class="min-h-screen flex flex-col gap-8">
    <p class="text-dark-900 font-extrabold text-xl">Bonjour John,</p>
    <div class="grid gap-4">
        <div class="flex justify-between items-center">
            <p class="font-extrabold text-lg text-dark-800">Active case</p>
            <div class="flex items-center gap-1">
                <p class="text-dark-600 text-sm">Switch Case</p>
                <ChevronDown />
            </div>
        </div>
        <CurrentCase />
    </div>
    <div class="flex gap-4">
        <Input placeholder="Recherche parmi les enregistrements" />
        <button class="text-dark-500">Edit</button>
    </div>
    <div class="flex flex-col gap-4">
        <p class="text-dark-700 font-bold text-base">Recordings</p>
        <div class="flex flex-col gap-2">
            {#each Array(3) as _, index}
                <Recording id={index} />
            {/each}
        </div>
    </div>
</div>

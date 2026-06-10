<script lang="ts">
    import type { Case, CaseStatus } from "$lib/types/case";

    interface Props {
        currentCase?: Case | null;
    }

    let { currentCase = null }: Props = $props();

    const statusLabels: Record<CaseStatus, string> = {
        in_progress: "En cours",
        ready: "Prêt",
        released: "Clôturé",
    };

    const statusBadgeClasses: Record<CaseStatus, string> = {
        in_progress: "bg-accent-subtle text-accent-subtle-foreground",
        ready: "bg-success text-success-foreground",
        released: "bg-muted text-muted-foreground",
    };

    const statusLabel = $derived(currentCase ? (statusLabels[currentCase.status] ?? currentCase.status) : "");
    const badgeClass = $derived(currentCase ? (statusBadgeClasses[currentCase.status] ?? "bg-dark-900 text-white") : "bg-dark-900 text-white");
    const displayId = $derived(currentCase?.externalReference ?? currentCase?.id ?? "");
</script>

{#if currentCase}
<div class="border border-dark-100 px-4 py-2 rounded-xl grid gap-4">
    <div class="flex gap-4 items-center">
        <span class="{badgeClass} px-4 py-2 rounded-2xl">{statusLabel}</span>
        <p class="text-dark-600 text-sm">ID: {displayId}</p>
    </div>
    <p class="text-dark-900 font-extrabold text-lg">{currentCase.title}</p>
    <div class="flex items-center gap-8">
        {#if currentCase.clientName}
        <div class="flex flex-col gap-1 text-xs">
            <p class="uppercase text-dark-500">client</p>
            <p class="text-dark-800 font-bold uppercase">{currentCase.clientName}</p>
        </div>
        {/if}
        {#if currentCase.clientNumber}
        <div class="flex flex-col gap-1 text-xs">
            <p class="uppercase text-dark-500">n° client</p>
            <p class="uppercase">{currentCase.clientNumber}</p>
        </div>
        {/if}
    </div>
</div>
{:else}
<div class="border border-dark-100 px-4 py-2 rounded-xl flex items-center justify-center h-24">
    <p class="text-dark-500 text-sm">Aucune affaire active</p>
</div>
{/if}

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

    const statusLabel = $derived(currentCase ? (statusLabels[currentCase.status] ?? currentCase.status) : "");
    const displayId = $derived(currentCase?.externalReference ?? currentCase?.id ?? "");
    const clientMeta = $derived(
        [currentCase?.clientName, currentCase?.clientNumber].filter(Boolean).join(" · ")
    );
</script>

{#if currentCase}
<div class="bg-background-alt rounded-2xl shadow-card p-4 flex flex-col gap-1.5">
    <div class="flex items-start justify-between gap-3">
        <p class="text-dark-900 font-semibold text-[17px] leading-snug break-words flex-1">{currentCase.title}</p>
        <span class="text-dark-400 text-sm flex-shrink-0 mt-0.5">{statusLabel}</span>
    </div>
    {#if clientMeta}
    <p class="text-dark-400 text-sm">{clientMeta}</p>
    {/if}
    {#if displayId}
    <p class="text-dark-300 text-xs font-mono">{displayId}</p>
    {/if}
</div>
{:else}
<p class="text-dark-400 text-sm py-1">Aucune affaire active</p>
{/if}

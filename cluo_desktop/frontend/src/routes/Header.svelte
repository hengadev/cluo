<script lang="ts">
    import { Folder } from "@lucide/svelte";
    import Search from "$lib/custom/header/Search.svelte";
    import { items, type HeaderItem } from "$lib/constructor/header";
    import ThemeToggle from "$lib/components/ThemeToggle.svelte";
    // TODO: pour la partie client ou type d'enquete, il faut un select avec tous les anciens clients + un bouton plus pour ajouter un nouvel element
    import { Button, Tooltip } from "bits-ui";
    import { currentCase } from "$lib/stores/case";
    import { fetchCase, fetchClient } from "$lib/services/api";
    import type { Case, Client, CaseStatus } from "$lib/types/entities";

    const STATUS_LABELS: Record<CaseStatus, string> = {
        in_progress: "En cours",
        ready: "Prêt",
        released: "Clôturé",
    };

    const STATUS_CLASSES: Record<CaseStatus, string> = {
        in_progress: "bg-blue-100 text-blue-800",
        ready: "bg-green-100 text-green-800",
        released: "bg-purple-100 text-purple-800",
    };

    let caseData: Case | null = $state(null);
    let clientData: Client | null = $state(null);

    $effect(() => {
        const caseId = $currentCase.id;
        if (!caseId) {
            caseData = null;
            clientData = null;
            return;
        }
        fetchCase(caseId).then(async (c) => {
            caseData = c;
            const client = await fetchClient(c.clientId);
            clientData = client;
        });
    });
</script>

<div class="header border-1 border-dark-50 animate-fade-in" style="animation-delay: 100ms;">
    <div class="grid">
        <div class="left">
            {#if caseData}
                <div class="current-case">
                    <div class="p-2 rounded-input bg-foreground">
                        <Folder size={16} class="text-background" />
                    </div>
                    <p>{clientData?.name ?? '…'}</p>
                </div>
                <p>&bull;</p>
                <p>{caseData.title}</p>
                <span class="status-badge {STATUS_CLASSES[caseData.status]}">
                    {STATUS_LABELS[caseData.status]}
                </span>
            {:else}
                <div class="no-case">
                    <Folder size={16} />
                    <p>Aucune affaire ouverte</p>
                </div>
            {/if}
        </div>
    </div>
    <Search />
    <div class="flex align-center gap-2">
        <div class="buttons">
            <ThemeToggle />
            {#each items as item}
                {@const DialogOrPopover = item.uiComponent}
                <DialogOrPopover>
                    {@render headerItem(item)}
                </DialogOrPopover>
            {/each}
        </div>
    </div>
</div>

{#snippet headerItem(item: HeaderItem)}
    {@const Icon = item.icon}
    <Tooltip.Provider>
        <Tooltip.Root delayDuration={100}>
            <Tooltip.Trigger
                class="rounded-10px p-3 bg-dark-50 ring-offset-background active:scale-[0.98] active:transition:all
		focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center focus-visible:ring-2 focus-visible:ring-offset-2
                    hover:bg-dark-100/50 hover:scale-105 transition-all duration-200 {item.bg} text-{item.fg}"
            >
                <Button.Root class="cursor-pointer">
                    <Icon size={24} strokeWidth={1.75} />
                </Button.Root>
            </Tooltip.Trigger>
            <Tooltip.Content sideOffset={8} side="bottom">
                <div
                    class="rounded-input text-[1rem] align-center bg-foreground text-background font-medium border-dark-10 shadow-popover outline-hidden z-0 flex items-center justify-center border p-2"
                >
                    {item.title}
                </div>
            </Tooltip.Content>
        </Tooltip.Root>
    </Tooltip.Provider>
{/snippet}

<style>
    .header {
        grid-area: header;
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 0.5rem 2rem;
        gap: 2rem;
    }
    .left,
    .buttons {
        flex: 1;
    }
    .left {
        display: flex;
        gap: 0.5rem;
        align-items: center;
        font-weight: 500;
    }
    .current-case {
        display: flex;
        gap: 0.5rem;
        align-items: center;
    }
    .status-badge {
        display: inline-block;
        padding: 0.125rem 0.5rem;
        border-radius: 9999px;
        font-size: 0.75rem;
        font-weight: 500;
    }
    .no-case {
        display: flex;
        gap: 0.5rem;
        align-items: center;
        color: var(--foreground-alt);
        font-weight: 400;
    }
    .buttons {
        display: flex;
        justify-content: right;
        gap: 0.5rem;
        margin-left: auto;
        text-align: right;
    }
</style>

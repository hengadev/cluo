<script lang="ts">
    import { Dialog } from "bits-ui";
    import { X, Check } from "@lucide/svelte";
    import type { Case } from "$lib/types/case";

    interface Props {
        open: boolean;
        cases: Case[];
        activeId: string | null;
        onselect: (c: Case) => void;
    }

    let { open = $bindable(false), cases, activeId, onselect }: Props = $props();

    const statusLabels: Record<string, string> = {
        in_progress: "En cours",
        ready: "Prêt",
        released: "Clôturé",
    };

    function pick(c: Case) {
        onselect(c);
        open = false;
    }
</script>

<Dialog.Root bind:open>
    <Dialog.Portal>
        <Dialog.Overlay class="fixed inset-0 bg-black/40 z-40" />
        <Dialog.Content
            class="fixed bottom-0 left-0 right-0 z-50 bg-white rounded-t-2xl px-4 pt-4 pb-8 max-h-[80vh] flex flex-col gap-4 shadow-xl"
        >
            <div class="flex items-center justify-between">
                <Dialog.Title class="font-extrabold text-lg text-dark-900">
                    Choisir une affaire
                </Dialog.Title>
                <Dialog.Close class="text-dark-400 hover:text-dark-700 transition-colors cursor-pointer">
                    <X size={20} />
                </Dialog.Close>
            </div>

            <div class="flex flex-col gap-2 overflow-y-auto">
                {#each cases as c (c.id)}
                    <button
                        onclick={() => pick(c)}
                        class="flex items-center justify-between px-4 py-3 rounded-xl border text-left transition-colors cursor-pointer
                            {activeId === c.id
                                ? 'border-dark-900 bg-dark-50'
                                : 'border-dark-100 hover:bg-dark-50'}"
                    >
                        <div class="flex flex-col gap-1">
                            <p class="font-bold text-dark-900 text-sm">{c.title}</p>
                            <div class="flex items-center gap-2 text-xs text-dark-500">
                                <span>{statusLabels[c.status] ?? c.status}</span>
                                {#if c.externalReference}
                                    <span>·</span>
                                    <span>{c.externalReference}</span>
                                {/if}
                                {#if c.clientName}
                                    <span>·</span>
                                    <span>{c.clientName}</span>
                                {/if}
                            </div>
                        </div>
                        {#if activeId === c.id}
                            <Check size={18} class="text-dark-900 shrink-0" />
                        {/if}
                    </button>
                {/each}
            </div>
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>

<script lang="ts">
    import { Dialog } from "bits-ui";
    import { X, Check, Search } from "@lucide/svelte";
    import type { Case } from "$lib/types/case";
    import { getCases } from "$lib/api";

    interface Props {
        open: boolean;
        activeId: string | null;
        onselect: (c: Case) => void;
    }

    let { open = $bindable(false), activeId, onselect }: Props = $props();

    const LIMIT = 8;

    let allCases = $state<Case[]>([]);
    let loading = $state(false);
    let query = $state("");
    let statusFilter = $state<string | null>(null);
    let showAll = $state(false);

    const statusLabels: Record<string, string> = {
        in_progress: "En cours",
        ready: "Prêt",
        released: "Clôturé",
    };

    $effect(() => {
        if (open) {
            if (allCases.length === 0) {
                loading = true;
                getCases().then((cases) => {
                    allCases = cases;
                    loading = false;
                });
            }
        } else {
            query = "";
            statusFilter = null;
            showAll = false;
        }
    });

    const filtered = $derived(
        allCases.filter((c) => {
            const matchesStatus = !statusFilter || c.status === statusFilter;
            const q = query.toLowerCase();
            const matchesQuery =
                !q ||
                c.title.toLowerCase().includes(q) ||
                (c.clientName?.toLowerCase().includes(q) ?? false) ||
                (c.externalReference?.toLowerCase().includes(q) ?? false);
            return matchesStatus && matchesQuery;
        })
    );

    const visible = $derived(showAll ? filtered : filtered.slice(0, LIMIT));
    const remaining = $derived(filtered.length - LIMIT);

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

            <!-- Search -->
            <div class="flex items-center gap-2 px-3 py-2 border border-dark-100 rounded-xl bg-dark-50">
                <Search size={16} class="text-dark-400 shrink-0" />
                <input
                    bind:value={query}
                    placeholder="Rechercher une affaire..."
                    class="flex-1 bg-transparent text-sm text-dark-900 placeholder:text-dark-400 outline-none"
                />
            </div>

            <!-- Status filter chips -->
            <div class="flex gap-2 flex-wrap">
                {#each [null, "in_progress", "ready", "released"] as s}
                    <button
                        onclick={() => { statusFilter = s; showAll = false; }}
                        class="px-3 py-1 rounded-full text-xs font-medium transition-colors cursor-pointer
                            {statusFilter === s
                                ? 'bg-dark-900 text-white'
                                : 'bg-dark-50 text-dark-500 hover:bg-dark-100'}"
                    >
                        {s === null ? "Toutes" : statusLabels[s]}
                    </button>
                {/each}
            </div>

            <!-- Case list -->
            <div class="flex flex-col gap-2 overflow-y-auto">
                {#if loading}
                    <p class="text-sm text-dark-400 text-center py-6">Chargement...</p>
                {:else if filtered.length === 0}
                    <p class="text-sm text-dark-400 text-center py-6">Aucune affaire trouvée.</p>
                {:else}
                    {#each visible as c (c.id)}
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

                    {#if !showAll && remaining > 0}
                        <button
                            onclick={() => (showAll = true)}
                            class="text-sm text-dark-500 hover:text-dark-900 transition-colors py-2 cursor-pointer"
                        >
                            Voir plus ({remaining})
                        </button>
                    {/if}
                {/if}
            </div>
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>

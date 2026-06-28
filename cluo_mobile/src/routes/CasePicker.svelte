<script lang="ts">
    import { Dialog } from "bits-ui";
    import Spinner from "$lib/components/ui/Spinner.svelte";
    import { X, Check, Search } from "@lucide/svelte";
    import type { Case, CaseStatus } from "$lib/types/case";
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
    let fetchError = $state<string | null>(null);
    let query = $state("");
    let statusFilter = $state<CaseStatus | null>(null);
    let showAll = $state(false);

    const statusLabels: Record<CaseStatus, string> = {
        in_progress: "En cours",
        ready: "Prêt",
        released: "Clôturé",
    };

    const statusOptions: (CaseStatus | null)[] = [null, "in_progress", "ready", "released"];

    $effect(() => {
        if (open) {
            loading = true;
            fetchError = null;
            getCases()
                .then((cases) => { allCases = cases; })
                .catch(() => { fetchError = "Impossible de charger les affaires."; })
                .finally(() => { loading = false; });
        } else {
            query = "";
            statusFilter = null;
            showAll = false;
        }
    });

    function escHtml(s: string): string {
        return s.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;").replace(/"/g, "&quot;");
    }

    function fuzzyMatch(text: string, pattern: string): { match: boolean; score: number; indices: number[] } {
        const t = text.toLowerCase();
        const p = pattern.toLowerCase();
        let score = 0;
        let tIdx = 0;
        const indices: number[] = [];
        let consecutive = 0;

        for (let pIdx = 0; pIdx < p.length; pIdx++) {
            let found = false;
            for (; tIdx < t.length; tIdx++) {
                if (t[tIdx] === p[pIdx]) {
                    indices.push(tIdx);
                    score += 1 + consecutive * 2;
                    consecutive++;
                    tIdx++;
                    found = true;
                    break;
                } else {
                    consecutive = 0;
                }
            }
            if (!found) return { match: false, score: 0, indices: [] };
        }
        // Bonus for earlier match position
        if (indices.length > 0) score += Math.max(0, 10 - indices[0]);
        return { match: true, score, indices };
    }

    function highlightFuzzy(text: string, indices: number[]): string {
        if (!indices.length) return escHtml(text);
        const set = new Set(indices);
        let out = "";
        for (let i = 0; i < text.length; i++) {
            const ch = escHtml(text[i]);
            out += set.has(i)
                ? `<mark class="bg-dark-200 text-dark-900 rounded-sm not-italic font-bold">${ch}</mark>`
                : ch;
        }
        return out;
    }

    interface MatchedCase {
        c: Case;
        score: number;
        titleHtml: string;
    }

    const filteredMatches = $derived.by((): MatchedCase[] => {
        const matchesStatus = (c: Case) => !statusFilter || c.status === statusFilter;
        const q = query.trim();

        if (!q) {
            return allCases
                .filter(matchesStatus)
                .map((c) => ({ c, score: 0, titleHtml: escHtml(c.title) }));
        }

        const results: MatchedCase[] = [];
        for (const c of allCases) {
            if (!matchesStatus(c)) continue;
            const titleMatch = fuzzyMatch(c.title, q);
            const clientMatch = c.clientName ? fuzzyMatch(c.clientName, q) : { match: false, score: 0, indices: [] };
            const refMatch = c.externalReference ? fuzzyMatch(c.externalReference, q) : { match: false, score: 0, indices: [] };

            if (!titleMatch.match && !clientMatch.match && !refMatch.match) continue;

            results.push({
                c,
                score: Math.max(titleMatch.score, clientMatch.score, refMatch.score),
                titleHtml: highlightFuzzy(c.title, titleMatch.indices),
            });
        }

        return results.sort((a, b) => b.score - a.score);
    });

    const visible = $derived(showAll ? filteredMatches : filteredMatches.slice(0, LIMIT));
    const remaining = $derived(filteredMatches.length - LIMIT);

    function pick(c: Case) {
        onselect(c);
        open = false;
    }
</script>

<Dialog.Root bind:open>
    <Dialog.Portal>
        <Dialog.Overlay class="fixed inset-0 bg-black/40 z-40" />
        <Dialog.Content
            class="fixed bottom-0 left-0 right-0 z-50 bg-background rounded-t-2xl px-4 pt-4 pb-8 max-h-[80vh] flex flex-col gap-4 shadow-popover"
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
                    class="flex-1 bg-transparent text-base text-dark-900 placeholder:text-dark-400 outline-none"
                />
            </div>

            <!-- Status filter chips -->
            <div class="flex gap-2 flex-wrap">
                {#each statusOptions as s}
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
                    <div class="flex items-center justify-center p-8">
                        <Spinner size="md" />
                    </div>
                {:else if fetchError}
                    <p class="text-sm text-red-500 text-center py-6">{fetchError}</p>
                {:else if filteredMatches.length === 0}
                    <p class="text-sm text-dark-400 text-center py-6">Aucune affaire trouvée.</p>
                {:else}
                    {#each visible as { c, titleHtml } (c.id)}
                        <button
                            onclick={() => pick(c)}
                            class="flex items-center justify-between px-4 py-3 rounded-xl border text-left transition-colors cursor-pointer
                                {activeId === c.id
                                    ? 'border-dark-900 bg-dark-50'
                                    : 'border-dark-100 hover:bg-dark-50'}"
                        >
                            <div class="flex flex-col gap-1 min-w-0 flex-1">
                                <p class="font-bold text-dark-900 text-sm truncate">{@html titleHtml}</p>
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

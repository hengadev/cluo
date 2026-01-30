<script lang="ts">
    import { currentCase } from "$lib/stores/case";

    // For now, a simple dashboard that shows a mock case selection
    // This will be replaced with actual case list from backend

    const mockCases = [
        { id: "CASE-2024-0847", title: "Affaire Dupont", status: "En cours" },
        { id: "CASE-2024-0852", title: "Affaire Martin", status: "Nouveau" },
        { id: "CASE-2024-0855", title: "Affaire Bernard", status: "En attente" },
    ];

    function selectCase(caseId: string) {
        currentCase.setCase(caseId);
    }
</script>

<div class="p-8">
    <h1 class="text-3xl font-bold mb-8">Tableau de bord</h1>

    <div class="grid gap-6">
        <section>
            <h2 class="text-xl font-semibold mb-4">Dossiers récents</h2>
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {#each mockCases as caseItem}
                    <button
                        class="border border-border-card rounded-card p-4 bg-background hover:border-border-input-hover transition-colors text-left"
                        onclick={() => selectCase(caseItem.id)}
                    >
                        <h3 class="font-semibold text-foreground">{caseItem.title}</h3>
                        <p class="text-sm text-muted-foreground">{caseItem.id}</p>
                        <span class="inline-block mt-2 px-2 py-1 text-xs rounded-full bg-muted text-foreground">
                            {caseItem.status}
                        </span>
                    </button>
                {/each}
            </div>
        </section>
    </div>
</div>

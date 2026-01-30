<script lang="ts">
    import { currentCase } from "$lib/stores/case";
    import { page } from "$app/stores";
    import { goto } from "$app/navigation";

    // Update the current case store when navigating to a case's documents
    $effect(() => {
        const caseId = $page.params.id;
        if (caseId && caseId !== $currentCase.id) {
            currentCase.setCase(caseId);
        }
    });

    const documentTypes = [
        { type: "facture", title: "Facture", icon: "📄" },
        { type: "mandat", title: "Mandat", icon: "🤝" },
        { type: "devis", title: "Devis", icon: "📝" },
    ];
</script>

<div class="p-8">
    <h1 class="text-3xl font-bold mb-8">Documents</h1>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {#each documentTypes as doc}
            <button
                class="border border-border-card rounded-card p-6 bg-background hover:border-border-input-hover transition-colors text-left"
                onclick={() => goto(`/cases/{$page.params.id}/documents/${doc.type}`)}
            >
                <span class="text-4xl mb-4 block">{doc.icon}</span>
                <h3 class="font-semibold text-foreground text-lg">{doc.title}</h3>
                <p class="text-sm text-muted-foreground mt-2">Voir les {doc.title.toLowerCase()}</p>
            </button>
        {/each}
    </div>
</div>

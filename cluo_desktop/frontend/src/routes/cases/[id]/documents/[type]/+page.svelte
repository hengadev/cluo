<script lang="ts">
    import { currentCase } from "$lib/stores/case";
    import { page } from "$app/stores";

    // Update the current case store when navigating to a case's document type
    $effect(() => {
        const caseId = $page.params.id;
        if (caseId && caseId !== $currentCase.id) {
            currentCase.setCase(caseId);
        }
    });

    // Map document types to their display names
    const docType = $derived($page.params.type);

    const docTypeNames: Record<string, string> = {
        facture: "Factures",
        mandat: "Mandats",
        devis: "Devis"
    };

    const displayName = $derived(docTypeNames[docType] || docType);
</script>

<div class="p-8">
    <h1 class="text-2xl font-bold mb-4">{displayName}</h1>
    <p class="text-muted-foreground">Gestion des {displayName.toLowerCase()} du dossier</p>
</div>

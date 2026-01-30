<script lang="ts">
    import { currentCase } from "$lib/stores/case";
    import { page } from "$app/stores";
    import Facture from "$lib/custom/content/Facture.svelte";
    import Mandat from "$lib/custom/content/Mandat.svelte";
    import Devis from "$lib/custom/content/Devis.svelte";

    // Update the current case store when navigating to a case's document type
    $effect(() => {
        const caseId = $page.params.id;
        if (caseId && caseId !== $currentCase.id) {
            currentCase.setCase(caseId);
        }
    });

    // Map document types to their components
    const docType = $derived($page.params.type);

    // Get the component to render
    const Component = $derived(() => {
        switch (docType) {
            case "facture":
                return Facture;
            case "mandat":
                return Mandat;
            case "devis":
                return Devis;
            default:
                return null;
        }
    });
</script>

{#if Component()}
    {@const Comp = Component()}
    <Comp />
{:else}
    <div class="p-8">
        <p class="text-muted-foreground">Type de document inconnu: {docType}</p>
    </div>
{/if}

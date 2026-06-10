<script lang="ts">
    import { currentCase } from "$lib/stores/case";
    import { page } from "$app/stores";

    // Update the current case store when navigating to a case's reseaux
    $effect(() => {
        const caseId = $page.params.id;
        if (caseId && caseId !== $currentCase.id) {
            currentCase.setCase(caseId);
        }
    });
</script>

<div class="page-content">
    <h1 class="text-2xl font-bold">Réseaux</h1>
    <p class="text-muted-foreground">Gestion des réseaux du dossier</p>
</div>

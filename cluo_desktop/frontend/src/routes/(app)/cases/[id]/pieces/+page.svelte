<script lang="ts">
    import { currentCase } from "$lib/stores/case";
    import { page } from "$app/stores";

    // Update the current case store when navigating to a case's pieces
    $effect(() => {
        const caseId = $page.params.id;
        if (caseId && caseId !== $currentCase.id) {
            currentCase.setCase(caseId);
        }
    });
</script>

<div class="p-8">
    <h1 class="text-2xl font-bold mb-4">Pièces</h1>
    <p class="text-muted-foreground">Gestion des pièces du dossier</p>
</div>

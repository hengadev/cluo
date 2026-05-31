<script lang="ts">
    import Informations from "./informations/Informations.svelte";
    import { currentCase } from "$lib/stores/case";
    import { page } from "$app/stores";

    // Update the current case store when navigating to a case
    $effect(() => {
        const caseId = $page.params.id;
        if (caseId && caseId !== $currentCase.id) {
            currentCase.setCase(caseId);
        }
    });
</script>

<Informations />

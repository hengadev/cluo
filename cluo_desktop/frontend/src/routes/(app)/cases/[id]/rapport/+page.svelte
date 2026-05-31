<script lang="ts">
    import Rapport from "./Rapport.svelte";
    import { currentCase } from "$lib/stores/case";
    import { page } from "$app/stores";

    // Update the current case store when navigating to a case's rapport
    $effect(() => {
        const caseId = $page.params.id;
        if (caseId && caseId !== $currentCase.id) {
            currentCase.setCase(caseId);
        }
    });
</script>

<Rapport />

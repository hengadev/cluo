<script lang="ts">
    import { Mic } from "@lucide/svelte";
    import { currentCase } from "$lib/stores/case";
    import { page } from "$app/stores";

    $effect(() => {
        const caseId = $page.params.id;
        if (caseId && caseId !== $currentCase.id) {
            currentCase.setCase(caseId);
        }
    });
</script>

<div class="p-8 flex flex-col flex-1 min-h-0 gap-6">
    <div>
        <h1 class="text-2xl font-bold text-foreground">Enregistrements</h1>
        <p class="text-sm text-muted-foreground">Enregistrements audio du dossier</p>
    </div>

    <div class="border border-dashed border-border rounded-lg p-12 bg-muted/20 flex flex-col items-center justify-center gap-4 flex-1 min-h-[50vh]">
        <Mic class="w-12 h-12 text-muted-foreground" />
        <p class="text-muted-foreground text-center">Aucun enregistrement pour ce dossier.</p>
    </div>
</div>

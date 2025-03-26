<script lang="ts">
    import Informations from "$lib/custom/content/informations/Informations.svelte";
    import Photos from "$lib/custom/content/photos/Photos.svelte";
    import Facture from "$lib/custom/content/Facture.svelte";
    import Rapport from "$lib/custom/content/Rapport.svelte";
    import Mandat from "$lib/custom/content/Mandat.svelte";
    import Devis from "$lib/custom/content/Devis.svelte";
    import Reseaux from "$lib/custom/content/Reseaux.svelte";

    import { type SidebarState, SIDEBAR_STATES } from "$lib/types/sidebar";

    type Props = {
        sidebarState: SidebarState;
    };
    let { sidebarState }: Props = $props();
</script>

<div class="content">
    {#if sidebarState === SIDEBAR_STATES.Informations}
        <Informations />
    {:else if sidebarState === SIDEBAR_STATES.Photos}
        <Photos />
    {:else if sidebarState === SIDEBAR_STATES.Facture}
        <Facture />
    {:else if sidebarState === SIDEBAR_STATES.Rapport}
        <Rapport />
    {:else if sidebarState === SIDEBAR_STATES.Mandat}
        <Mandat />
    {:else if sidebarState === SIDEBAR_STATES.Devis}
        <Devis />
    {:else if sidebarState === SIDEBAR_STATES.Reseaux}
        <Reseaux />
    {/if}
</div>

<style>
    .content {
        grid-area: content;
        height: 100%;
        overflow: auto;
    }
</style>

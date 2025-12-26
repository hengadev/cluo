<script lang="ts">
    import Informations from "$lib/custom/content/informations/Informations.svelte";
    import Users from "$lib/custom/content/users/Users.svelte";
    import Photos from "$lib/custom/content/photos/Photos.svelte";
    import Facture from "$lib/custom/content/Facture.svelte";
    import Rapport from "$lib/custom/content/rapport/Rapport.svelte";
    import Mandat from "$lib/custom/content/Mandat.svelte";
    import Devis from "$lib/custom/content/Devis.svelte";
    import Reseaux from "$lib/custom/content/Reseaux.svelte";
    import Pieces from "$lib/custom/content/pieces/Pieces.svelte";
    import Recordings from "$lib/custom/content/recordings/Recordings.svelte";

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
    {:else if sidebarState === SIDEBAR_STATES.Utilisateurs}
        <Users />
    {:else if sidebarState === SIDEBAR_STATES.Pièces}
        <Pieces />
    {:else if sidebarState === SIDEBAR_STATES.Enregistrement}
        <Recordings />
    {:else if sidebarState === SIDEBAR_STATES.Facture}
        <Facture id="94t8a9fawg48ae3" />
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

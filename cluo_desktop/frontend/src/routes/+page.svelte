<script lang="ts">
    import "../app.css";
    import "../reset.css";

    import Header from "./Header.svelte";
    import Sidebar from "./Sidebar.svelte";
    import Toaster from "$lib/custom/global/toast/Toaster.svelte";
    import Informations from "./Informations.svelte";
    import Photos from "./Photos.svelte";
    import Facture from "./Facture.svelte";
    import Rapport from "./Rapport.svelte";
    import Mandat from "./Mandat.svelte";
    import Devis from "./Devis.svelte";
    import Reseaux from "./Reseaux.svelte";

    import { setToastContext } from "$lib/custom/global/toast/state.svelte";
    setToastContext();

    import { type SidebarState, SIDEBAR_STATES } from "$lib/types/sidebar";
    let sidebarState = $state<SidebarState>(SIDEBAR_STATES.Informations);
</script>

<Toaster />

<div class="page">
    <Header />
    <Sidebar bind:sidebarState />
    {@render renderContent()}
</div>

{#snippet renderContent()}
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
{/snippet}

<style lang="postcss">
    .page {
        display: grid;
        grid-template-areas:
            "sidebar header"
            "sidebar content";
        /* grid-template-rows: 100px 1fr; */
        grid-template-rows: auto 1fr;
        grid-template-columns: auto 1fr;
        height: 100vh;
    }
</style>

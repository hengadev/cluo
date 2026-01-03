<script lang="ts">
    import "../app.css";
    import "../reset.css";

    import Header from "./Header.svelte";
    import Sidebar from "./Sidebar.svelte";
    import Toaster from "$lib/custom/global/toast/Toaster.svelte";

    import { setToastContext } from "$lib/custom/global/toast/state.svelte";
    setToastContext();

    import { type SidebarState, SIDEBAR_STATES } from "$lib/types/sidebar";
    import Content from "./Content.svelte";
    // let sidebarState = $state<SidebarState>(SIDEBAR_STATES.Informations);
    let sidebarState = $state<SidebarState>(SIDEBAR_STATES.Photos);
</script>

<Toaster />

<div class="page">
    <Header />
    <Sidebar bind:sidebarState />
    <Content {sidebarState} />
</div>

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

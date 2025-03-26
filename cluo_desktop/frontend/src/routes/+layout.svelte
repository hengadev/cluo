<script lang="ts">
    import "../app.css";
    import "../reset.css";
    let { children } = $props();

    import Header from "./Header.svelte";
    import Sidebar from "./Sidebar.svelte";
    import Toaster from "$lib/custom/global/toast/Toaster.svelte";

    import { setToastContext } from "$lib/custom/global/toast/state.svelte";
    setToastContext();
</script>

<Toaster />

<div class="page">
    <Header />
    <Sidebar />
    {@render children()}
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

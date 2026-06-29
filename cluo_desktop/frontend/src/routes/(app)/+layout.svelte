<script lang="ts">
    import { onMount } from "svelte";
    import Header from "../Header.svelte";
    import Sidebar from "../Sidebar.svelte";
    import { initBackgroundCheckers } from "$lib/services/backgroundInit";

    // One-shot background checkers (overdue-invoice notifications in v1). Run
    // once per session on layout mount; each checker is idempotent.
    onMount(() => {
        initBackgroundCheckers();
    });
</script>

<div class="page">
    <Header />
    <Sidebar />
    <div class="content">
        <slot />
    </div>
</div>

<style lang="postcss">
    .page {
        display: grid;
        grid-template-areas:
            "sidebar header"
            "sidebar content";
        grid-template-rows: auto 1fr;
        grid-template-columns: auto 1fr;
        height: 100dvh;
    }

    .content {
        grid-area: content;
        height: 100%;
        overflow: auto;
        display: flex;
        flex-direction: column;
    }
</style>

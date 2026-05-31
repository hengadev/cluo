<script lang="ts">
    import Header from "../Header.svelte";
    import Sidebar from "../Sidebar.svelte";
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
        height: 100vh;
    }

    .content {
        grid-area: content;
        height: 100%;
        overflow: auto;
        display: flex;
        flex-direction: column;
    }
</style>

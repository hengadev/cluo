<script lang="ts">
    import "../app.css";
    import { page } from "$app/state";

    import Footer from "./Footer.svelte";

    let { children } = $props();

    // Hide footer on auth pages
    const showFooter = $derived(!page.url.pathname.startsWith("/auth"));
</script>

<div class="relative px-4">
    {@render children()}
    {#if showFooter}
        <div class="fixed bottom-0 inset-x-0 z-50">
            <Footer />
        </div>
    {/if}
</div>

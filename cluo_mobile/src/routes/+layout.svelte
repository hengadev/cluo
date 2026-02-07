<script lang="ts">
    import "../app.css";
    import { page } from "$app/state";
    import { onMount } from 'svelte';
    import { dev } from '$app/environment';

    import Footer from "./Footer.svelte";
    import Snackbar from "$lib/components/Snackbar.svelte";

    let { children } = $props();

    // Hide footer on auth pages
    const showFooter = $derived(!page.url.pathname.startsWith("/auth"));

    // Register service worker for PWA functionality
    onMount(() => {
        if (!dev && 'serviceWorker' in navigator) {
            navigator.serviceWorker.register('/service-worker.js');
        }
    });
</script>

<div class="relative px-4">
    {@render children()}
    {#if showFooter}
        <div class="fixed bottom-0 inset-x-0 z-50">
            <Footer />
        </div>
    {/if}
</div>

<Snackbar />

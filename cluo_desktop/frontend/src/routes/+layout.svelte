<script lang="ts">
    import { onMount } from 'svelte';
    import { theme } from '$lib/stores/theme';
    import { updateDialogOpen } from '$lib/stores/update';
    import '../app.css';
    import '../reset.css';

    import Header from "./Header.svelte";
    import Sidebar from "./Sidebar.svelte";
    import Toaster from "$lib/custom/global/toast/Toaster.svelte";
    import UpdateDialog from "$lib/custom/global/UpdateDialog.svelte";

    import { setToastContext } from "$lib/custom/global/toast/state.svelte";
    setToastContext();

    onMount(async () => {
        theme.set($theme);
        try {
            const { CheckForUpdate } = await import('$lib/wailsjs/go/updater/Updater');
            const info = await CheckForUpdate();
            if (info.available) {
                updateDialogOpen.set(true);
            }
        } catch {
            // ManifestURL not configured (dev build) or network error — silently skip
        }
    });

    // Disable prerendering for dynamic routes
    export const prerender = false;
    export const ssr = false;
</script>

<Toaster />
<UpdateDialog bind:open={$updateDialogOpen} />

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
    }
</style>

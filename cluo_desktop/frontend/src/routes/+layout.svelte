<script lang="ts">
    import { onMount } from 'svelte';
    import { theme } from '$lib/stores/theme';
    import { updateDialogOpen } from '$lib/stores/update';
    import '../app.css';
    import '../reset.css';

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
</script>

<Toaster />
<UpdateDialog bind:open={$updateDialogOpen} />

<slot />

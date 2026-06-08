<script lang="ts">
    import { onMount } from 'svelte';
    import { theme } from '$lib/stores/theme';
    import { updateDialogOpen } from '$lib/stores/update';
    import '../app.css';
    import '../reset.css';

    import Toaster from "$lib/custom/global/toast/Toaster.svelte";
    import UpdatePrompt from "$lib/custom/global/UpdatePrompt.svelte";
    import UpdateDialog from "$lib/custom/global/UpdateDialog.svelte";

    import { setToastContext } from "$lib/custom/global/toast/state.svelte";
    setToastContext();

    onMount(() => {
        theme.set($theme);
    });
</script>

<Toaster />
<UpdatePrompt />
<UpdateDialog bind:open={$updateDialogOpen} />

<slot />

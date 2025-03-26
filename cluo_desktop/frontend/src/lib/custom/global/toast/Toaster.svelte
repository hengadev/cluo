<script lang="ts">
    import Toast from "./Toast.svelte";
    import { getToastContext } from "$lib/custom/global/toast/state.svelte";

    let toastState = getToastContext();
</script>

<div class="absolute bottom-4 right-8 flex flex-col-reverse gap-2 shadow-sm">
    {#each toastState.toasts as toast}
        <Toast {toast} />
    {/each}
</div>

<style>
</style>

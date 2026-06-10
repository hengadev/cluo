<script lang="ts">
    import { snackbar, type Snackbar } from "$lib/stores/snackbar";
    import { X } from "@lucide/svelte";
    import { fly } from "svelte/transition";

    let current: Snackbar | null = $state(null);

    snackbar.subscribe((value) => {
        current = value;
    });

    function handleAction() {
        if (current?.action) {
            current.action.onClick();
            snackbar.dismiss();
        }
    }
</script>

{#if current}
    <div
        class="fixed bottom-24 left-4 right-4 z-50"
        transition:fly={{ y: 50, duration: 200 }}
    >
        <div class="flex items-center gap-3 px-4 py-3 bg-red-600 text-white rounded-xl shadow-popover">
            <p class="flex-1 text-sm font-medium">{current.message}</p>
            {#if current.action}
                <button
                    onclick={handleAction}
                    class="px-3 py-1 bg-white/20 hover:bg-white/30 rounded-lg text-sm font-semibold transition-colors"
                >
                    {current.action.label}
                </button>
            {/if}
            <button
                onclick={() => snackbar.dismiss()}
                class="opacity-70 hover:opacity-100 transition-opacity"
            >
                <X size={18} />
            </button>
        </div>
    </div>
{/if}

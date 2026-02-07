<script lang="ts">
    import { Button } from "bits-ui";
    import { X, CircleCheck, CircleX, CircleAlert } from "@lucide/svelte";
    import { type ToastLevel, type Toast, TOAST_LEVELS } from "./type";
    import { fly } from "svelte/transition";
    import { quintOut } from "svelte/easing";

    import { getToastContext } from "$lib/custom/global/toast/state.svelte";
    const toastState = getToastContext();

    type Props = { toast: Toast };

    let { toast }: Props = $props();
</script>

<div
    class="flex items-start gap-4 border-2 border-muted rounded-input p-4 max-w-[420px]"
    transition:fly={{ y: 10, duration: 500, easing: quintOut }}
>
    {@render toastIcon(toast.level)}
    <div>
        <p class="font-semibold">{toast.title}</p>
        <p class="">{toast.message}</p>
    </div>
    <Button.Root
        class="cursor-pointer"
        onclick={() => toastState.remove(toast.id)}
    >
        <X size={24} />
    </Button.Root>
</div>

{#snippet toastIcon(level: ToastLevel)}
    {#if level === TOAST_LEVELS.Info}
        <CircleCheck size={32} color="green" />
    {:else if level === TOAST_LEVELS.Error}
        <CircleX size={32} color="red" />
    {:else if level === TOAST_LEVELS.Alert}
        <CircleAlert size={32} color="orange" />
    {/if}
{/snippet}

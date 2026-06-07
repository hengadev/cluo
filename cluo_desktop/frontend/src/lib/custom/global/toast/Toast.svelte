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
    class="flex items-start gap-3 bg-background border border-border-input rounded-card-sm shadow-popover p-4 w-[360px]"
    transition:fly={{ y: 8, duration: 300, easing: quintOut }}
>
    {@render toastIcon(toast.level)}
    <div class="flex-1 min-w-0">
        <p class="text-sm font-semibold">{toast.title}</p>
        <p class="text-sm text-foreground-alt mt-0.5 leading-snug">{toast.message}</p>
    </div>
    <Button.Root
        class="cursor-pointer text-foreground-alt hover:text-foreground transition-colors mt-0.5 shrink-0"
        onclick={() => toastState.remove(toast.id)}
    >
        <X size={14} />
    </Button.Root>
</div>

{#snippet toastIcon(level: ToastLevel)}
    {#if level === TOAST_LEVELS.Info}
        <CircleCheck size={16} class="text-success shrink-0 mt-0.5" />
    {:else if level === TOAST_LEVELS.Error}
        <CircleX size={16} class="text-destructive shrink-0 mt-0.5" />
    {:else if level === TOAST_LEVELS.Alert}
        <CircleAlert size={16} class="text-tertiary shrink-0 mt-0.5" />
    {/if}
{/snippet}

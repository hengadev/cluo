<script lang="ts">
    import {
        MessageSquareMore,
        Lightbulb,
        PencilLine,
        AudioLines,
    } from "@lucide/svelte";

    interface Props {
        width: number;
        minWidth: number;
        maxWidth: number;
    }

    let { width, minWidth, maxWidth }: Props = $props();

    let selected = $state(0);

    type AIButton = {
        icon: typeof import("@lucide/svelte").Icon;
        title: string;
    };

    const aiButtons: AIButton[] = [
        { icon: MessageSquareMore, title: "Chat" },
        { icon: Lightbulb, title: "Ideas" },
        { icon: PencilLine, title: "Review" },
        { icon: AudioLines, title: "Audio" },
    ];
</script>

<div
    class="grid grid-rows-[auto_1fr] border-l-1 border-dark-50 overflow-x-hidden"
    style="width: {width}px; min-width: {minWidth}px; max-width: {maxWidth}px;"
>
    <div
        class="flex justify-center border-b-1 border-dark-50 text-center"
        style="gap: {width <= minWidth ? '0.5rem' : '1.5rem'};"
    >
        {#each aiButtons as item, index}
            {@render button(index, item)}
        {/each}
    </div>
</div>

{#snippet button(index: number, item: AIButton)}
    {@const Icon = item.icon}
    {@const isSelected = index === selected}
    {@const isCompact = width <= minWidth}
    <button
        class="flex items-center gap-2 transition-all"
        style="padding: {isCompact ? '0.75rem 0.5rem' : '1rem'};"
        onclick={() => (selected = index)}
    >
        {#if !isCompact}
            <Icon
                size={24}
                strokeWidth={1.5}
                class={isSelected ? `text-dark-900` : `text-dark-300`}
            />
        {/if}
        <p class={isSelected ? `text-dark-900` : `text-dark-300`}>
            {item.title}
        </p>
    </button>
{/snippet}

<script lang="ts">
    import {
        MessageSquareMore,
        Lightbulb,
        PencilLine,
        AudioLines,
    } from "@lucide/svelte";

    let selected = $state(0);

    // Panel resizing state
    const MIN_AI_PANEL_WIDTH = 400;
    const MAX_AI_PANEL_WIDTH = 800;
    const DEFAULT_AI_PANEL_WIDTH = 400;

    let aiPanelWidth = $state(DEFAULT_AI_PANEL_WIDTH);
    let isDragging = $state(false);
    let containerRef: HTMLDivElement;

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

<div class="grid grid-cols-2 h-full">
    <div class="">rapport the part with the text editor</div>
    <div class="grid grid-rows-[auto_1fr] border-l-1 border-dark-50">
        <div
            class="flex justify-center gap-6 border-b-1 border-dark-50 text-center"
        >
            {#each aiButtons as item, index}
                {@render button(index, item)}
            {/each}
        </div>
    </div>
</div>

{#snippet button(index: number, item: AIButton)}
    {@const Icon = item.icon}
    {@const isSelected = index === selected}
    <button class="flex gap-2 p-4" onclick={() => (selected = index)}>
        <Icon
            size={32}
            strokeWidth={1.5}
            class={isSelected ? `text-dark-900` : `text-dark-300`}
        />
        <p class={isSelected ? `text-dark-900` : `text-dark-300`}>
            {item.title}
        </p>
    </button>
{/snippet}

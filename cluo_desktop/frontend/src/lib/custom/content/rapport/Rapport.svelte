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

    function handleMouseDown(e: MouseEvent) {
        isDragging = true;
        e.preventDefault();
    }

    function handleMouseMove(e: MouseEvent) {
        if (!isDragging || !containerRef) return;

        const containerRect = containerRef.getBoundingClientRect();
        const newWidth = containerRect.right - e.clientX;

        aiPanelWidth = Math.max(
            MIN_AI_PANEL_WIDTH,
            Math.min(MAX_AI_PANEL_WIDTH, newWidth),
        );
    }

    function handleMouseUp() {
        isDragging = false;
    }
</script>

<svelte:window onmousemove={handleMouseMove} onmouseup={handleMouseUp} />

<div class="flex h-full" bind:this={containerRef}>
    <!-- Left panel: Text editor (flexible width) -->
    <div class="flex-1 min-w-0">rapport the part with the text editor</div>

    <!-- Draggable divider -->
    <div
        class="w-1 cursor-col-resize hover:bg-dark-300 transition-colors relative group"
        onmousedown={handleMouseDown}
        role="separator"
        aria-orientation="vertical"
    >
        <div
            class="absolute inset-y-0 -left-1 -right-1 group-hover:bg-dark-10"
        ></div>
    </div>

    <!-- Right panel: AI panel (constrained width) -->
    <div
        class="grid grid-rows-[auto_1fr] border-l-1 border-dark-50"
        style="width: {aiPanelWidth}px; min-width: {MIN_AI_PANEL_WIDTH}px; max-width: {MAX_AI_PANEL_WIDTH}px;"
    >
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

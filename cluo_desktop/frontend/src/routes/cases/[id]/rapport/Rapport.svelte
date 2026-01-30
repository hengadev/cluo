<script lang="ts">
    import AISidebar from "./AISidebar.svelte";
    import TextEditor from "./TextEditor.svelte";

    // Panel resizing state
    const MIN_AI_PANEL_WIDTH = 400;
    const MAX_AI_PANEL_WIDTH = 800;
    const DEFAULT_AI_PANEL_WIDTH = 400;

    let aiPanelWidth = $state(DEFAULT_AI_PANEL_WIDTH);
    let isDragging = $state(false);
    let containerRef: HTMLDivElement;

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
    <div class="flex-1 min-w-0">
        <TextEditor />
    </div>

    <!-- Draggable divider -->
    <div
        class="w-1 cursor-col-resize hover:bg-dark-300 transition-colors relative group"
        onmousedown={handleMouseDown}
        role="separator"
        aria-orientation="vertical"
        tabindex="0"
        aria-label="Resize panels"
    >
        <div
            class="absolute inset-y-0 -left-1 -right-1 group-hover:bg-dark-10"
        ></div>
    </div>

    <!-- Right panel: AI Sidebar -->
    <AISidebar
        width={aiPanelWidth}
        minWidth={MIN_AI_PANEL_WIDTH}
        maxWidth={MAX_AI_PANEL_WIDTH}
    />
</div>

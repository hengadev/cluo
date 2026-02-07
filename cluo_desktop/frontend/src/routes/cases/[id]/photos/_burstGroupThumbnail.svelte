<script lang="ts">
    import type { Image } from "./types";
    import { Layers } from "@lucide/svelte";

    interface Props {
        images: Image[];
        isInReport: boolean;
        selectMode: boolean;
        isSelected: boolean;
        onSelectionChange: () => void;
        onOpenBurstModal: () => void;
    }

    let {
        images,
        isInReport,
        selectMode,
        isSelected,
        onSelectionChange,
        onOpenBurstModal,
    }: Props = $props();

    let firstImage = $derived(images[0]);
    let count = $derived(images.length);

    function handleClick(): void {
        if (selectMode) {
            onSelectionChange();
        } else {
            // Open burst modal to select which images to add
            onOpenBurstModal();
        }
    }
</script>

<div
    class="burst-thumbnail-wrapper group border border-border-card rounded-card bg-background hover:border-border-input-hover transition-colors cursor-pointer {selectMode
        ? 'select-none'
        : ''} {isSelected ? 'ring-2 ring-primary ring-offset-2' : ''}"
    onclick={handleClick}
>
    <!-- Stacked Cards Container -->
    <div class="burst-stack">
        <!-- Stack Layer 3 (bottom) -->
        {#if count >= 4}
            <div class="burst-card-layer burst-card-layer-3">
                <img
                    src={images[3].url}
                    alt=""
                    class="w-full h-full object-cover"
                />
            </div>
        {/if}

        <!-- Stack Layer 2 -->
        {#if count >= 3}
            <div class="burst-card-layer burst-card-layer-2">
                <img
                    src={images[2].url}
                    alt=""
                    class="w-full h-full object-cover"
                />
            </div>
        {/if}

        <!-- Stack Layer 1 -->
        {#if count >= 2}
            <div class="burst-card-layer burst-card-layer-1">
                <img
                    src={images[1].url}
                    alt=""
                    class="w-full h-full object-cover"
                />
            </div>
        {/if}

        <!-- Main Image (front) -->
        <div class="burst-card-main">
            <img
                src={firstImage.url}
                alt={firstImage.filename}
                class="w-full h-full object-cover"
            />
        </div>
    </div>

    <!-- Count Badge -->
    <div
        class="absolute bottom-2 left-2 bg-background/90 backdrop-blur-sm rounded-full px-2.5 py-1 flex items-center gap-1.5 shadow-sm z-10"
    >
        <Layers size={12} class="text-foreground" />
        <span class="text-xs font-semibold text-foreground">{count}</span>
    </div>

    {#if selectMode}
        <!-- Selection Checkbox Overlay -->
        <div
            class="absolute top-2 left-2 w-6 h-6 rounded-md border-2 flex items-center justify-center transition-all z-10 {isSelected
                ? 'bg-primary border-primary'
                : 'bg-white/80 border-white'}"
        >
            {#if isSelected}
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="14"
                    height="14"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="3"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    class="text-primary-foreground"
                >
                    <polyline points="20 6 9 17 4 12"></polyline>
                </svg>
            {/if}
        </div>
    {:else if isInReport}
        <!-- Check Indicator (some photos in report) -->
        <div
            class="absolute top-2 right-2 w-8 h-8 bg-green-500 rounded-full flex items-center justify-center shadow-mini z-10"
        >
            <svg
                xmlns="http://www.w3.org/2000/svg"
                width="16"
                height="16"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="3"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="text-white"
            >
                <polyline points="20 6 9 17 4 12"></polyline>
            </svg>
        </div>
    {/if}
</div>

<style>
    .burst-thumbnail-wrapper {
        position: relative;
        aspect-ratio: 1 / 1;
        overflow: visible;
    }

    .burst-stack {
        position: relative;
        width: 90%;
        height: 90%;
        margin: 5%;
    }

    .burst-card-main,
    .burst-card-layer {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        border-radius: 0.25rem;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
        transition: transform 0.2s ease;
        overflow: hidden;
    }

    /* Main front card */
    .burst-card-main {
        z-index: 4;
        border: 2px solid white;
    }

    /* Layer 1 - slightly behind */
    .burst-card-layer-1 {
        z-index: 3;
        transform: translate(4px, 3px);
        border: 2px solid white;
    }

    /* Layer 2 - more behind */
    .burst-card-layer-2 {
        z-index: 2;
        transform: translate(8px, 6px);
        border: 2px solid white;
    }

    /* Layer 3 - furthest behind */
    .burst-card-layer-3 {
        z-index: 1;
        transform: translate(12px, 9px);
        border: 2px solid white;
    }

    /* On hover, spread the stack */
    .burst-thumbnail-wrapper:hover .burst-card-layer-1 {
        transform: translate(6px, 4px);
    }

    .burst-thumbnail-wrapper:hover .burst-card-layer-2 {
        transform: translate(12px, 8px);
    }

    .burst-thumbnail-wrapper:hover .burst-card-layer-3 {
        transform: translate(18px, 12px);
    }
</style>

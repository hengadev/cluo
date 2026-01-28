<script lang="ts">
    import type { Image } from "./types";
    import { Plus, Check } from "@lucide/svelte";

    interface Props {
        image: Image;
        isInReport: boolean;
        selectMode: boolean;
        isSelected: boolean;
        onSelectionChange: () => void;
        onAdd: () => void;
    }

    let {
        image,
        isInReport,
        selectMode,
        isSelected,
        onSelectionChange,
        onAdd,
    }: Props = $props();

    function handleClick(): void {
        if (selectMode) {
            onSelectionChange();
        } else if (!isInReport) {
            onAdd();
        }
    }
</script>

<div
    class="relative group border border-border-card rounded-card overflow-hidden bg-background hover:border-border-input-hover transition-colors cursor-pointer {selectMode
        ? 'select-none'
        : ''} {isSelected
        ? 'ring-2 ring-primary ring-offset-2'
        : ''}"
    onclick={handleClick}
>
    <!-- Image Preview -->
    <img
        src={image.url}
        alt={image.filename}
        class="w-full aspect-square object-cover"
    />

    {#if selectMode}
        <!-- Selection Checkbox Overlay -->
        <div class="absolute top-2 left-2 w-6 h-6 rounded-md border-2 flex items-center justify-center transition-all {isSelected
            ? 'bg-primary border-primary'
            : 'bg-white/80 border-white'}">
            {#if isSelected}
                <Check size={14} class="text-primary-foreground" />
            {/if}
        </div>
    {:else}
        <!-- Normal Mode Actions -->
        {#if !isInReport}
            <!-- Add Button (visible on hover) -->
            <button
                class="absolute top-2 right-2 w-8 h-8 bg-white rounded-full shadow-mini flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity hover:bg-white/90 active:scale-95"
                onclick={onAdd}
                title="Ajouter au rapport"
            >
                <Plus size={16} class="text-foreground" />
            </button>
        {:else}
            <!-- Check Indicator (always visible when in report) -->
            <div
                class="absolute top-2 right-2 w-8 h-8 bg-green-500 rounded-full flex items-center justify-center shadow-mini"
            >
                <Check size={16} class="text-white" />
            </div>
        {/if}
    {/if}
</div>

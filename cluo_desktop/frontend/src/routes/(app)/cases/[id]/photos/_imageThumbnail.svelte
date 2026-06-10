<script lang="ts">
    import type { Image } from "./types";
    import { Plus, Check, Eye, EyeOff, Trash2 } from "@lucide/svelte";

    interface Props {
        image: Image;
        isInReport: boolean;
        selectMode: boolean;
        isSelected: boolean;
        onSelectionChange: () => void;
        onAdd: () => void;
        onTogglePublish?: () => void;
        onDelete?: () => void;
    }

    let {
        image,
        isInReport,
        selectMode,
        isSelected,
        onSelectionChange,
        onAdd,
        onTogglePublish,
        onDelete,
    }: Props = $props();

    let showActions = $state(false);

    function handleClick(): void {
        if (selectMode) {
            onSelectionChange();
        } else if (!isInReport) {
            onAdd();
        }
    }

    function handleContextMenu(e: MouseEvent): void {
        e.preventDefault();
        showActions = !showActions;
    }

    function handleMouseLeave(): void {
        showActions = false;
    }
</script>

<div
    class="relative group border border-border-card rounded-card overflow-hidden bg-background hover:border-border-input-hover hover:shadow-card transition-all duration-300 cursor-pointer aspect-square {selectMode
        ? 'select-none'
        : ''} {isSelected ? 'ring-2 ring-foreground ring-offset-2 ring-offset-background' : ''} {image.isPublished ? '' : 'opacity-80'}"
    onclick={handleClick}
    oncontextmenu={handleContextMenu}
    onmouseleave={handleMouseLeave}
>
    <!-- Image Preview -->
    <img
        src={image.url}
        alt={image.filename}
        class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
        loading="lazy"
    />

    <!-- Published indicator -->
    <div class="absolute bottom-2 left-2 px-1.5 py-0.5 rounded text-[10px] font-medium {image.isPublished
        ? 'bg-success/20 text-success'
        : 'bg-muted text-muted-foreground'}">
        {image.isPublished ? "Publié" : "Brouillon"}
    </div>

    {#if selectMode}
        <!-- Selection Checkbox Overlay -->
        <div
            class="absolute top-2 left-2 w-7 h-7 rounded-md border-2 flex items-center justify-center transition-all shadow-mini {isSelected
                ? 'bg-foreground border-foreground'
                : 'bg-background border-dark'}"
        >
            {#if isSelected}
                <Check size={16} class="text-background" />
            {/if}
        </div>
    {:else}
        <!-- Normal Mode Actions -->
        {#if !isInReport}
            <!-- Add Button (visible on hover) -->
            <button
                class="absolute top-2 right-2 w-8 h-8 bg-background rounded-full shadow-mini flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity hover:bg-background/90 active:scale-95"
                onclick={(e) => { e.stopPropagation(); onAdd(); }}
                title="Ajouter au rapport"
            >
                <Plus size={16} class="text-foreground" />
            </button>
        {:else}
            <!-- Check Indicator (always visible when in report) -->
            <div
                class="absolute top-2 right-2 w-8 h-8 bg-success rounded-full flex items-center justify-center shadow-mini"
            >
                <Check size={16} class="text-success-foreground" />
            </div>
        {/if}

        <!-- Context action menu -->
        {#if showActions}
            <div class="absolute top-2 right-2 flex flex-col gap-1 z-10">
                {#if onTogglePublish}
                    <button
                        class="w-8 h-8 bg-background rounded-full shadow-mini flex items-center justify-center hover:bg-muted active:scale-95"
                        onclick={(e) => { e.stopPropagation(); onTogglePublish(); showActions = false; }}
                        title={image.isPublished ? "Dépublier" : "Publier"}
                    >
                        {#if image.isPublished}
                            <EyeOff size={14} class="text-muted-foreground" />
                        {:else}
                            <Eye size={14} class="text-foreground" />
                        {/if}
                    </button>
                {/if}
                {#if onDelete}
                    <button
                        class="w-8 h-8 bg-background rounded-full shadow-mini flex items-center justify-center hover:bg-destructive/10 active:scale-95"
                        onclick={(e) => { e.stopPropagation(); onDelete(); showActions = false; }}
                        title="Supprimer"
                    >
                        <Trash2 size={14} class="text-destructive" />
                    </button>
                {/if}
            </div>
        {/if}
    {/if}
</div>

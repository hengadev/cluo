<script lang="ts">
    import ImageThumbnail from "./_imageThumbnail.svelte";
    import type { Image } from "./types";
    import type { ViewMode } from "./_floatingToolbar.svelte";

    interface Props {
        images: Image[];
        reportedIds: Set<string>;
        viewMode: ViewMode;
        selectMode: boolean;
        selectedIds: Set<string>;
        onSelectionChange: (id: string) => void;
        onAdd: (image: Image) => void;
    }

    let {
        images,
        reportedIds,
        viewMode,
        selectMode,
        selectedIds,
        onSelectionChange,
        onAdd,
    }: Props = $props();

    let gridCols = $derived(viewMode === "grid-compact" ? "grid-cols-3" : "grid-cols-2");
</script>

<div class="border border-border-card rounded-card p-4 bg-background">
    <h2 class="text-lg font-semibold mb-8 text-foreground">
        Bibliothèque d'images du projet
    </h2>

    {#if images.length === 0}
        <div
            class="flex flex-col items-center justify-center py-12 text-center"
        >
            <p class="text-muted-foreground">Aucune image disponible</p>
        </div>
    {:else}
        <div
            class="grid {gridCols} gap-3 max-h-[70vh] overflow-y-auto py-4 px-2"
        >
            {#each images as image (image.id)}
                <ImageThumbnail
                    {image}
                    isInReport={reportedIds.has(image.id)}
                    {selectMode}
                    isSelected={selectedIds.has(image.id)}
                    onSelectionChange={() => onSelectionChange(image.id)}
                    onAdd={() => onAdd(image)}
                />
            {/each}
        </div>
    {/if}
</div>

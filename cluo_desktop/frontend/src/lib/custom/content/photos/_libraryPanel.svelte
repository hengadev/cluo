<script lang="ts">
    import ImageThumbnail from "./_imageThumbnail.svelte";
    import BurstGroupThumbnail from "./_burstGroupThumbnail.svelte";
    import type { Image, BurstGroup } from "./types";
    import type { ViewMode } from "./_floatingToolbar.svelte";

    interface Props {
        images: Image[];
        burstGroups: BurstGroup[];
        reportedIds: Set<string>;
        viewMode: ViewMode;
        selectMode: boolean;
        selectedIds: Set<string>;
        onSelectionChange: (id: string) => void;
        onAdd: (image: Image) => void;
    }

    let {
        images,
        burstGroups,
        reportedIds,
        viewMode,
        selectMode,
        selectedIds,
        onSelectionChange,
        onAdd,
    }: Props = $props();

    let gridStyle = $derived(
        viewMode === "grid-compact"
            ? "grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));"
            : "grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));",
    );

    // Track which image IDs are in burst groups
    let burstGroupImageIds = $derived(
        new Set(burstGroups.flatMap((g) => g.images.map((img) => img.id))),
    );

    // Filter out images that are in burst groups
    let soloImages = $derived(
        images.filter((img) => !burstGroupImageIds.has(img.id)),
    );
</script>

<div
    class="border border-border-card rounded-card p-4 bg-background h-full flex flex-col w-full"
>
    <h2 class="text-lg font-semibold mb-4 text-foreground py-2">
        Bibliothèque d'images du projet
    </h2>

    {#if images.length === 0}
        <div
            class="flex flex-col items-center justify-center flex-1 text-center"
        >
            <p class="text-muted-foreground">Aucune image disponible</p>
        </div>
    {:else}
        <div class="flex-1 min-h-0 overflow-y-auto p-4">
            <div class="grid gap-3" style={gridStyle}>
                {#each burstGroups as group (group.id)}
                    <BurstGroupThumbnail
                        images={group.images}
                        isInReport={group.images.some((img) =>
                            reportedIds.has(img.id),
                        )}
                        {selectMode}
                        isSelected={group.images.some((img) =>
                            selectedIds.has(img.id),
                        )}
                        onSelectionChange={() => {
                            for (const img of group.images) {
                                onSelectionChange(img.id);
                            }
                        }}
                        onAdd={onAdd}
                    />
                {/each}
                {#each soloImages as image (image.id)}
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
        </div>
    {/if}
</div>

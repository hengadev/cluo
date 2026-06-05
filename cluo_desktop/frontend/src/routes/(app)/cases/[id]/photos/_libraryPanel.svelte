<script lang="ts">
    import ImageThumbnail from "./_imageThumbnail.svelte";
    import BurstGroupThumbnail from "./_burstGroupThumbnail.svelte";
    import BurstGroupModal from "./_burstGroupModal.svelte";
    import type { Image, BurstGroup } from "./types";

    interface Props {
        images: Image[];
        burstGroups: BurstGroup[];
        reportedIds: Set<string>;
        selectMode: boolean;
        selectedIds: Set<string>;
        onSelectionChange: (id: string) => void;
        onAdd: (image: Image) => void;
        onImport?: () => void;
    }

    let {
        images,
        burstGroups,
        reportedIds,
        selectMode,
        selectedIds,
        onSelectionChange,
        onAdd,
        onImport,
    }: Props = $props();

    let openBurstGroupId = $state<string | null>(null);

    function handleOpenBurstModal(groupId: string): void {
        openBurstGroupId = groupId;
    }

    function handleCloseBurstModal(): void {
        openBurstGroupId = null;
    }

    let gridStyle = "grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));";

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
            class="flex flex-col items-center justify-center flex-1 text-center gap-3"
        >
            <p class="text-muted-foreground">Aucune image disponible</p>
            {#if onImport}
                <button
                    type="button"
                    onclick={onImport}
                    class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center px-4 text-sm font-semibold active:scale-[0.98] cursor-pointer"
                >
                    Importer des photos
                </button>
            {/if}
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
                        onOpenBurstModal={() => handleOpenBurstModal(group.id)}
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

    <!-- Burst Group Modal -->
    {#if openBurstGroupId}
        {@const burstGroup = burstGroups.find((g) => g.id === openBurstGroupId)}
        {#if burstGroup}
            <BurstGroupModal
                images={burstGroup.images}
                {reportedIds}
                onClose={handleCloseBurstModal}
                {onAdd}
            />
        {/if}
    {/if}
</div>

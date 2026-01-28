<script lang="ts">
    import ReportImageCard from "./_reportImageCard.svelte";
    import { Image } from "@lucide/svelte";
    import type { ReportImage } from "./types";

    interface Props {
        images: ReportImage[];
        onRemove: (id: string) => void;
        onReorder: (images: ReportImage[]) => void;
        onCaptionChange: (id: string, caption: string) => void;
    }

    let { images, onRemove, onReorder, onCaptionChange }: Props = $props();

    let draggedId = $state<string | null>(null);
    let dragOverId = $state<string | null>(null);

    function handleDragStart(e: DragEvent, id: string): void {
        draggedId = id;
        if (e.dataTransfer) {
            e.dataTransfer.effectAllowed = "move";
            e.dataTransfer.setData("text/plain", id);
        }
    }

    function handleDragOver(e: DragEvent): void {
        e.preventDefault();
        if (e.dataTransfer) {
            e.dataTransfer.dropEffect = "move";
        }
    }

    function handleDrop(e: DragEvent, targetId: string): void {
        e.preventDefault();
        if (!draggedId || draggedId === targetId) {
            dragOverId = null;
            return;
        }

        const draggedIndex = images.findIndex((img) => img.id === draggedId);
        const targetIndex = images.findIndex((img) => img.id === targetId);

        if (draggedIndex === -1 || targetIndex === -1) {
            dragOverId = null;
            return;
        }

        // Reorder and recalculate order numbers
        const reordered = [...images];
        const [removed] = reordered.splice(draggedIndex, 1);
        reordered.splice(targetIndex, 0, removed);

        onReorder(
            reordered.map((img, i) => ({
                ...img,
                order: i + 1,
            })),
        );

        draggedId = null;
        dragOverId = null;
    }

    function handleDragEnd(): void {
        draggedId = null;
        dragOverId = null;
    }

    function handleDragEnter(targetId: string): void {
        if (draggedId && draggedId !== targetId) {
            dragOverId = targetId;
        }
    }

    function handleDragLeave(): void {
        // Don't clear dragOverId here - it would flicker too much
        // It gets cleared on drop and drag end
    }
</script>

<div class="border border-border-card rounded-card p-4 bg-background">
    <div class="flex justify-between items-center mb-4">
        <h2 class="text-lg font-semibold text-foreground">Inclus dans le rapport</h2>
        <span class="text-sm text-muted-foreground">{images.length} image{images.length !== 1 ? 's' : ''}</span>
    </div>

    {#if images.length === 0}
        <!-- Empty State -->
        <div class="flex flex-col items-center justify-center py-12 text-center">
            <Image size={48} class="text-muted-foreground/50 mb-4" />
            <p class="text-muted-foreground">Aucune image dans le rapport</p>
            <p class="text-sm text-muted-foreground/70 mt-1">Cliquez sur le bouton + des images pour les ajouter</p>
        </div>
    {:else}
        <!-- Image List -->
        <div class="space-y-3 max-h-[70vh] overflow-y-auto pr-1">
            {#each images as image (image.id)}
                <ReportImageCard
                    {image}
                    onRemove={onRemove}
                    onCaptionChange={onCaptionChange}
                    isDragging={draggedId === image.id}
                    isDragOver={dragOverId === image.id}
                    onDragStart={handleDragStart}
                    onDragOver={handleDragOver}
                    onDragEnter={handleDragEnter}
                    onDrop={handleDrop}
                    onDragEnd={handleDragEnd}
                />
            {/each}
        </div>
    {/if}
</div>

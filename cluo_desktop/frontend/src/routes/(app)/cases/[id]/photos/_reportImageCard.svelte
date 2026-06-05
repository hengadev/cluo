<script lang="ts">
    import type { ReportImage } from "./types";
    import { GripVertical, X } from "@lucide/svelte";

    interface Props {
        image: ReportImage;
        onRemove: (id: string) => void;
        onCaptionChange: (id: string, caption: string) => void;
        isDragging: boolean;
        isDragOver: boolean;
        onDragStart: (e: DragEvent, id: string) => void;
        onDragOver: (e: DragEvent) => void;
        onDragEnter: (id: string) => void;
        onDrop: (e: DragEvent, id: string) => void;
        onDragEnd: () => void;
    }

    let {
        image,
        onRemove,
        onCaptionChange,
        isDragging,
        isDragOver,
        onDragStart,
        onDragOver,
        onDragEnter,
        onDrop,
        onDragEnd,
    }: Props = $props();

    let caption = $state(image.reportCaption);

    $effect(() => {
        caption = image.reportCaption;
    });

    function handleCaptionChange(): void {
        onCaptionChange(image.id, caption);
    }
</script>

<div
    class="flex gap-4 border rounded-card p-3 bg-muted/30 hover:shadow-md hover:-translate-y-0.5 transition-all duration-200 {isDragging
        ? 'opacity-50'
        : ''} {isDragOver
        ? 'border-accent bg-accent/50'
        : 'border-border-card'}"
    draggable="true"
    ondragstart={(e) => onDragStart(e, image.id)}
    ondragover={onDragOver}
    ondragenter={() => onDragEnter(image.id)}
    ondrop={(e) => onDrop(e, image.id)}
    ondragend={onDragEnd}
>
    <!-- Drag Handle + Order -->
    <div class="flex flex-col items-center gap-1 pt-2 select-none">
        <GripVertical class="text-muted-foreground cursor-grab active:cursor-grabbing" size={16} />
        <span class="text-sm font-semibold text-muted-foreground">{image.order}</span>
    </div>

    <!-- Thumbnail -->
    <img
        src={image.url}
        alt={image.filename}
        class="w-24 h-24 rounded-lg object-cover flex-shrink-0 bg-muted"
    />

    <!-- Content -->
    <div class="flex-1 min-w-0 flex flex-col">
        <h3 class="font-medium truncate mb-2 text-foreground" title={image.filename}>
            {image.filename}
        </h3>

        <!-- Caption Area -->
        <textarea
            bind:value={caption}
            onblur={handleCaptionChange}
            onkeydown={(e) => {
                if (e.key === "Enter" && !e.shiftKey) {
                    e.preventDefault();
                    (e.target as HTMLTextAreaElement).blur();
                }
            }}
            placeholder="Ajouter une légende..."
            class="w-full p-2 border border-border-input rounded-input text-sm resize-none focus:border-border-input-hover focus:outline-hidden bg-background"
            rows="2"
        />
    </div>

    <!-- Remove Button -->
    <button
        class="self-start p-1 text-muted-foreground hover:text-destructive hover:scale-110 active:scale-95 transition-all duration-200 rounded-md hover:bg-destructive/10"
        onclick={() => onRemove(image.id)}
        title="Retirer du rapport"
    >
        <X size={18} />
    </button>
</div>

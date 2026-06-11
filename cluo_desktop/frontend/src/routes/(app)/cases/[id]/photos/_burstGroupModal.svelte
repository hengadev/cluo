<script lang="ts">
    import type { Image } from "./types";
    import { X, Plus, Check } from "@lucide/svelte";

    interface Props {
        images: Image[];
        reportedIds: Set<string>;
        onClose: () => void;
        onAdd: (image: Image) => void;
    }

    let { images, reportedIds, onClose, onAdd }: Props = $props();

    let selectedImageIds = $state<Set<string>>(new Set());

    function toggleImageSelection(id: string): void {
        if (selectedImageIds.has(id)) {
            selectedImageIds = new Set([...selectedImageIds].filter((x) => x !== id));
        } else {
            selectedImageIds = new Set([...selectedImageIds, id]);
        }
    }

    function addSelectedImages(): void {
        for (const id of selectedImageIds) {
            const image = images.find((img) => img.id === id);
            if (image && !reportedIds.has(id)) {
                onAdd(image);
            }
        }
        onClose();
    }

    function addAllImages(): void {
        for (const image of images) {
            if (!reportedIds.has(image.id)) {
                onAdd(image);
            }
        }
        onClose();
    }

    function handleBackdropClick(e: MouseEvent): void {
        if (e.target === e.currentTarget) {
            onClose();
        }
    }
</script>

<div
    class="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4"
    onclick={handleBackdropClick}
>
    <div class="bg-background rounded-card shadow-popover max-w-4xl w-full max-h-[80vh] flex flex-col">
        <!-- Header -->
        <div class="flex items-center justify-between p-4 border-b border-border-card">
            <div>
                <h2 class="text-lg font-semibold text-foreground">
                    Sélectionner les photos
                </h2>
                <p class="text-sm text-muted-foreground">
                    {images.length} photo{images.length > 1 ? "s" : ""} dans ce groupe
                </p>
            </div>
            <button
                class="w-8 h-8 rounded-full flex items-center justify-center text-muted-foreground hover:text-foreground hover:bg-muted transition-interactive"
                onclick={onClose}
            >
                <X size={20} />
            </button>
        </div>

        <!-- Image Grid -->
        <div class="flex-1 overflow-y-auto p-4">
            <div class="grid grid-cols-4 gap-4">
                {#each images as image (image.id)}
                    <div
                        class="relative group aspect-square rounded-card overflow-hidden border-2 cursor-pointer transition-interactive {selectedImageIds.has(
                            image.id,
                        )
                            ? 'border-primary ring-2 ring-primary ring-offset-2'
                            : 'border-border-card hover:border-border-input-hover'}"
                        onclick={() => toggleImageSelection(image.id)}
                    >
                        <img
                            src={image.url}
                            alt={image.filename}
                            class="w-full h-full object-cover"
                        />

                        <!-- Selection Overlay -->
                        <div
                            class="absolute inset-0 bg-black/0 group-hover:bg-black/10 transition-interactive"
                        ></div>

                        <!-- Checkbox -->
                        <div
                            class="absolute top-2 left-2 w-7 h-7 rounded-md border-2 flex items-center justify-center transition-interactive shadow-mini {selectedImageIds.has(
                                image.id,
                            )
                                ? 'bg-foreground border-foreground'
                                : 'bg-background border-dark'}"
                        >
                            {#if selectedImageIds.has(image.id)}
                                <Check size={16} class="text-background" />
                            {/if}
                        </div>

                        <!-- Already in report indicator -->
                        {#if reportedIds.has(image.id)}
                            <div
                                class="absolute top-2 right-2 w-6 h-6 bg-success rounded-full flex items-center justify-center"
                            >
                                <Check size={14} class="text-success-foreground" />
                            </div>
                        {/if}
                    </div>
                {/each}
            </div>
        </div>

        <!-- Footer -->
        <div class="flex items-center justify-between p-4 border-t border-border-card">
            <p class="text-sm text-muted-foreground">
                {selectedImageIds.size} photo{selectedImageIds.size !== 1
                    ? "s"
                    : ""} sélectionnée{selectedImageIds.size !== 1 ? "s" : ""}
            </p>
            <div class="flex gap-2">
                <button
                    class="px-4 py-2 rounded-full text-foreground hover:bg-muted transition-interactive font-medium"
                    onclick={onClose}
                >
                    Annuler
                </button>
                <button
                    class="px-4 py-2 rounded-full bg-muted text-foreground hover:bg-muted/80 transition-interactive font-medium"
                    onclick={addAllImages}
                >
                    Tout ajouter
                </button>
                <button
                    class="px-4 py-2 rounded-full bg-primary text-primary-foreground hover:bg-primary/90 transition-interactive font-medium disabled:opacity-50 disabled:cursor-not-allowed"
                    onclick={addSelectedImages}
                    disabled={selectedImageIds.size === 0}
                >
                    Ajouter la sélection
                </button>
            </div>
        </div>
    </div>
</div>

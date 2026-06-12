<script lang="ts">
    import LibraryPanel from "./_libraryPanel.svelte";
    import ReportPanel from "./_reportPanel.svelte";
    import FloatingToolbar, {
        type SortMode,
        type LayoutMode,
    } from "./_floatingToolbar.svelte";

    import { fetchCaseMedia, uploadMedia, updateMedia, deleteMedia } from "$lib/services/api";
    import type { Image, ReportImage, BurstGroup } from "./types";
    import type { MediaFile } from "$lib/types/entities";
    import { onMount } from "svelte";
    import { page } from "$app/stores";
    import { getToastContext } from "$lib/custom/global/toast/state.svelte";
    import { TOAST_LEVELS } from "$lib/custom/global/toast/type";
    import Spinner from "$lib/components/Spinner.svelte";

    const toastState = getToastContext();

    let allImages = $state<Image[]>([]);
    let loading = $state(true);
    let uploading = $state(false);

    function mediaToImage(m: MediaFile): Image {
        return {
            id: m.id,
            caseId: m.caseId,
            url: m.url,
            filename: m.fileName,
            filesize: m.fileSize,
            caption: m.caption,
            isPublished: m.isPublished,
            createdAt: m.createdAt,
        };
    }

    // Load images from the API on mount
    onMount(async () => {
        loading = true;
        try {
            const caseId = $page.params.id;
            const response = await fetchCaseMedia(caseId, 'image');
            allImages = response.media.map(mediaToImage);
        } catch (error) {
            console.error("Failed to fetch images:", error);
            toastState.add(TOAST_LEVELS.Error, "Erreur", "Impossible de charger les photos.");
            allImages = [];
        } finally {
            loading = false;
        }
    });

    let reportImages = $state<ReportImage[]>([]);
    let reportedIds = $derived(new Set(reportImages.map((img) => img.id)));

    // Toolbar state
    let selectMode = $state(false);
    let sortMode = $state<SortMode>("newest");
    let layoutMode = $state<LayoutMode>("split");
    let burstGroupsEnabled = $state(true);
    let selectedIds = $state<Set<string>>(new Set());
    let fileInput = $state<HTMLInputElement>();

    // Sorted and filtered images for library panel
    let displayImages = $derived.by(() => {
        let sorted = [...allImages];
        switch (sortMode) {
            case "newest":
                sorted.sort(
                    (a, b) =>
                        new Date(b.createdAt).getTime() -
                        new Date(a.createdAt).getTime(),
                );
                break;
            case "oldest":
                sorted.sort(
                    (a, b) =>
                        new Date(a.createdAt).getTime() -
                        new Date(b.createdAt).getTime(),
                );
                break;
            case "filename":
                sorted.sort((a, b) => a.filename.localeCompare(b.filename));
                break;
        }
        return sorted;
    });

    // Burst group detection
    let burstGroups = $derived.by<BurstGroup[]>(() => {
        if (!burstGroupsEnabled) return [];

        const BURST_TIME_WINDOW_MS = 2000;
        const MIN_GROUP_SIZE = 3;

        const sorted = [...allImages].sort(
            (a, b) =>
                new Date(a.createdAt).getTime() -
                new Date(b.createdAt).getTime(),
        );

        const groups: Image[][] = [];
        let currentGroup: Image[] = [sorted[0]];

        for (let i = 1; i < sorted.length; i++) {
            const prevTime = new Date(
                currentGroup[currentGroup.length - 1].createdAt,
            ).getTime();
            const currTime = new Date(sorted[i].createdAt).getTime();
            const diffMs = currTime - prevTime;

            if (diffMs <= BURST_TIME_WINDOW_MS) {
                currentGroup.push(sorted[i]);
            } else {
                groups.push(currentGroup);
                currentGroup = [sorted[i]];
            }
        }
        groups.push(currentGroup);

        return groups
            .filter((g) => g.length >= MIN_GROUP_SIZE)
            .map((groupImages, index) => ({
                id: `burst-${index}`,
                images: groupImages,
                startTimestamp: groupImages[0].createdAt,
                endTimestamp: groupImages[groupImages.length - 1].createdAt,
            }));
    });

    function addToReport(image: Image): void {
        if (reportedIds.has(image.id)) return;

        const newReportImage: ReportImage = {
            ...image,
            order: reportImages.length + 1,
            reportCaption: image.caption || "",
        };

        reportImages = [...reportImages, newReportImage];
    }

    function removeFromReport(id: string): void {
        reportImages = reportImages
            .filter((img) => img.id !== id)
            .map((img, index) => ({ ...img, order: index + 1 }));
    }

    function updateOrder(reorderedImages: ReportImage[]): void {
        reportImages = reorderedImages.map((img, index) => ({
            ...img,
            order: index + 1,
        }));
    }

    function updateCaption(id: string, caption: string): void {
        reportImages = reportImages.map((img) =>
            img.id === id ? { ...img, reportCaption: caption } : img,
        );
    }

    // Toolbar handlers
    function handleImport(): void {
        fileInput?.click();
    }

    async function handleFileSelect(event: Event): Promise<void> {
        const target = event.target as HTMLInputElement;
        const files = target.files;
        if (!files || files.length === 0) return;

        const caseId = $page.params.id;
        uploading = true;

        for (const file of Array.from(files)) {
            if (!file.type.startsWith("image/")) continue;

            try {
                const media = await uploadMedia(caseId, file);
                const image = mediaToImage(media);
                allImages = [...allImages, image];
            } catch (err) {
                console.error("Failed to upload:", file.name, err);
                toastState.add(
                    TOAST_LEVELS.Error,
                    "Erreur d'import",
                    `Impossible d'importer « ${file.name} ».`,
                );
            }
        }

        uploading = false;
        target.value = "";
    }

    async function handleTogglePublish(image: Image): Promise<void> {
        try {
            await updateMedia(image.id, { isPublished: !image.isPublished });
            allImages = allImages.map((img) =>
                img.id === image.id
                    ? { ...img, isPublished: !img.isPublished }
                    : img,
            );
        } catch (err) {
            toastState.add(
                TOAST_LEVELS.Error,
                "Erreur",
                "Impossible de modifier la publication.",
            );
        }
    }

    async function handleDelete(image: Image): Promise<void> {
        try {
            await deleteMedia(image.id);
            allImages = allImages.filter((img) => img.id !== image.id);
            removeFromReport(image.id);
        } catch (err) {
            toastState.add(
                TOAST_LEVELS.Error,
                "Erreur",
                "Impossible de supprimer la photo.",
            );
        }
    }

    function handleSelectModeToggle(): void {
        selectMode = !selectMode;
        if (!selectMode) {
            selectedIds = new Set();
        }
    }

    function handleBurstGroupToggle(): void {
        burstGroupsEnabled = !burstGroupsEnabled;
    }

    function handleLayoutModeChange(mode: LayoutMode): void {
        layoutMode = mode;
    }
</script>

<div class="p-8 flex flex-col flex-1 min-h-0 gap-6">
    <!-- Hidden File Input -->
    <input
        type="file"
        bind:this={fileInput}
        accept="image/*"
        multiple
        onchange={handleFileSelect}
        class="hidden"
    />

    <!-- Header -->
    <div class="flex items-center justify-between">
        <div>
            <h1 class="text-2xl font-bold text-foreground">Photos</h1>
            <p class="text-sm text-muted-foreground">
                {allImages.length} photo{allImages.length !== 1 ? "s" : ""} · {allImages.filter(i => i.isPublished).length} publiée{allImages.filter(i => i.isPublished).length !== 1 ? "s" : ""}
            </p>
        </div>
        {#if allImages.length > 0 || uploading}
            <button
                type="button"
                onclick={handleImport}
                disabled={uploading}
                class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center px-4 text-sm font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50"
            >
                {uploading ? "Import en cours…" : "Importer des photos"}
            </button>
        {/if}
    </div>

    <!-- Panel Layout -->
    <div class="flex-1 min-h-0 overflow-hidden">
        {#if loading}
            <div class="flex items-center justify-center h-full">
                <Spinner size="lg" />
            </div>
        {:else if allImages.length === 0}
            <div class="border border-dashed border-border rounded-lg bg-muted/20 flex flex-col items-center justify-center h-full gap-4">
                <p class="text-muted-foreground">Aucune photo pour ce dossier.</p>
                <button
                    type="button"
                    onclick={handleImport}
                    class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center px-4 text-sm font-semibold active:scale-[0.98] cursor-pointer"
                >
                    Importer des photos
                </button>
            </div>
        {:else if layoutMode === "library"}
            <LibraryPanel
                images={displayImages}
                burstGroups={burstGroups}
                {reportedIds}
                {selectMode}
                {selectedIds}
                onSelectionChange={(id) => {
                    if (selectedIds.has(id)) {
                        selectedIds = new Set(
                            [...selectedIds].filter((x) => x !== id),
                        );
                    } else {
                        selectedIds = new Set([...selectedIds, id]);
                    }
                }}
                onAdd={addToReport}
                onImport={handleImport}
                onTogglePublish={handleTogglePublish}
                onDelete={handleDelete}
            />
        {:else if layoutMode === "split"}
            <div class="flex items-center gap-6 h-full overflow-y-auto">
                <LibraryPanel
                    images={displayImages}
                    burstGroups={burstGroups}
                    {reportedIds}
                    {selectMode}
                    {selectedIds}
                    onSelectionChange={(id) => {
                        if (selectedIds.has(id)) {
                            selectedIds = new Set(
                                [...selectedIds].filter((x) => x !== id),
                            );
                        } else {
                            selectedIds = new Set([...selectedIds, id]);
                        }
                    }}
                    onAdd={addToReport}
                    onTogglePublish={handleTogglePublish}
                    onDelete={handleDelete}
                />
                <ReportPanel
                    images={reportImages}
                    onRemove={removeFromReport}
                    onReorder={updateOrder}
                    onCaptionChange={updateCaption}
                />
            </div>
        {:else}
            <ReportPanel
                images={reportImages}
                onRemove={removeFromReport}
                onReorder={updateOrder}
                onCaptionChange={updateCaption}
            />
        {/if}
    </div>

    <!-- Floating Toolbar -->
    <FloatingToolbar
        {selectMode}
        {sortMode}
        {layoutMode}
        hasBurstGroups={burstGroupsEnabled}
        onSelectModeToggle={handleSelectModeToggle}
        onImport={handleImport}
        onBurstGroupToggle={handleBurstGroupToggle}
        onSortModeChange={(mode) => (sortMode = mode)}
        onLayoutModeChange={handleLayoutModeChange}
    />
</div>

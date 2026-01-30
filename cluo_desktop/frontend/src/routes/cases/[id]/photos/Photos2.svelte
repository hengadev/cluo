<script lang="ts">
    import LibraryPanel from "./_libraryPanel.svelte";
    import ReportPanel from "./_reportPanel.svelte";
    import FloatingToolbar, {
        type SortMode,
        type LayoutMode,
    } from "./_floatingToolbar.svelte";

    import { isMockEnabled } from "$lib/config";
    import { images as mockImages } from "./mockData";
    import { fetchCaseImages } from "$lib/services/api";
    import type { Image, ReportImage, BurstGroup } from "./types";
    import { onMount } from "svelte";

    let allImages = $state<Image[]>(mockImages);
    let loading = $state(false);

    // Load images based on mock flag
    onMount(async () => {
        if (!isMockEnabled()) {
            loading = true;
            try {
                // TODO: Get caseId from route params when API is ready
                const apiImages = await fetchCaseImages("CASE-2024-0847");
                if (apiImages.length > 0) {
                    allImages = apiImages as Image[];
                }
            } catch (error) {
                console.error("Failed to fetch images:", error);
                allImages = [];
            } finally {
                loading = false;
            }
        }
    });
    let reportImages = $state<ReportImage[]>([]);
    let reportedIds = $derived(new Set(reportImages.map((img) => img.id)));

    // Toolbar state
    let selectMode = $state(false);
    let sortMode = $state<SortMode>("newest");
    let layoutMode = $state<LayoutMode>("split");
    // let layoutMode = $state<LayoutMode>("library");
    let burstGroupsEnabled = $state(true);
    let selectedIds = $state<Set<string>>(new Set());
    let fileInput = $state<HTMLInputElement>();

    // Sorted and filtered images for library panel
    let displayImages = $derived(() => {
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
    let burstGroups = $derived<BurstGroup[]>(() => {
        if (!burstGroupsEnabled) return [];

        const BURST_TIME_WINDOW_MS = 2000; // 2 seconds
        const MIN_GROUP_SIZE = 3;

        // Sort by timestamp
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

        // Filter to only actual bursts and convert to BurstGroup type
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

        const newImages: Image[] = [];

        for (const file of Array.from(files)) {
            if (!file.type.startsWith("image/")) continue;

            const url = URL.createObjectURL(file);
            const newImage: Image = {
                id: `img-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
                caseId: "CASE-2024-0847",
                url,
                filename: file.name,
                filesize: file.size,
                caption: "",
                isPublished: false,
                createdAt: new Date().toISOString(),
            };
            newImages.push(newImage);
        }

        allImages = [...allImages, ...newImages];

        // Reset the input so the same files can be selected again if needed
        target.value = "";
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

<div class="content p-6 pr-6 h-[calc(100vh-80px)] flex flex-col">
    <!-- Hidden File Input -->
    <input
        type="file"
        bind:this={fileInput}
        accept="image/*"
        multiple
        onchange={handleFileSelect}
        class="hidden"
    />

    <!-- Panel Layout -->
    <div class="flex-1 min-h-0 overflow-hidden">
        {#if loading}
            <div class="flex items-center justify-center h-full">
                <p class="text-muted-foreground">Chargement des photos...</p>
            </div>
        {:else if allImages.length === 0}
            <div class="flex items-center justify-center h-full">
                <p class="text-muted-foreground">
                    Aucune photo disponible. {isMockEnabled() ? '' : '(API non configurée)'}
                </p>
            </div>
        {:else if layoutMode === "library"}
            <!-- Library Only -->
            <LibraryPanel
                images={displayImages()}
                burstGroups={burstGroups()}
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
            />
        {:else if layoutMode === "split"}
            <!-- Split View -->
            <!-- <div class="grid grid-cols-2 gap-6 h-full overflow-y-auto"> -->
            <div class="flex items-center gap-6 h-full overflow-y-auto">
                <LibraryPanel
                    images={displayImages()}
                    burstGroups={burstGroups()}
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
                />
                <ReportPanel
                    images={reportImages}
                    onRemove={removeFromReport}
                    onReorder={updateOrder}
                    onCaptionChange={updateCaption}
                />
            </div>
        {:else}
            <!-- Report Only -->
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

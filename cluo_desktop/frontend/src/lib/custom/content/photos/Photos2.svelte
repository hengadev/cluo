<script lang="ts">
    import LibraryPanel from "./_libraryPanel.svelte";
    import ReportPanel from "./_reportPanel.svelte";
    import FloatingToolbar, { type ViewMode, type SortMode, type LayoutMode } from "./_floatingToolbar.svelte";

    import { images } from "./mockData";
    import type { Image, ReportImage } from "./types";

    let allImages = $state<Image[]>(images);
    let reportImages = $state<ReportImage[]>([]);
    let reportedIds = $derived(new Set(reportImages.map((img) => img.id)));

    // Toolbar state
    let selectMode = $state(false);
    let viewMode = $state<ViewMode>("grid-compact");
    let sortMode = $state<SortMode>("newest");
    let layoutMode = $state<LayoutMode>("split");
    let burstGroupsEnabled = $state(false);
    let selectedIds = $state<Set<string>>(new Set());
    let fileInput = $state<HTMLInputElement>();

    // Sorted and filtered images for library panel
    let displayImages = $derived(() => {
        let sorted = [...allImages];
        switch (sortMode) {
            case "newest":
                sorted.sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime());
                break;
            case "oldest":
                sorted.sort((a, b) => new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime());
                break;
            case "filename":
                sorted.sort((a, b) => a.filename.localeCompare(b.filename));
                break;
        }
        return sorted;
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
        // TODO: Implement burst grouping logic
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
    {#if layoutMode === "library"}
        <!-- Library Only -->
        <LibraryPanel
            images={displayImages()}
            {reportedIds}
            {viewMode}
            {selectMode}
            selectedIds={selectedIds}
            onSelectionChange={(id) => {
                if (selectedIds.has(id)) {
                    selectedIds = new Set([...selectedIds].filter(x => x !== id));
                } else {
                    selectedIds = new Set([...selectedIds, id]);
                }
            }}
            onAdd={addToReport}
        />
    {:else if layoutMode === "split"}
        <!-- Split View -->
        <div class="grid grid-cols-2 gap-6 h-full">
            <LibraryPanel
                images={displayImages()}
                {reportedIds}
                {viewMode}
                {selectMode}
                selectedIds={selectedIds}
                onSelectionChange={(id) => {
                    if (selectedIds.has(id)) {
                        selectedIds = new Set([...selectedIds].filter(x => x !== id));
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
        {viewMode}
        {sortMode}
        {layoutMode}
        hasBurstGroups={burstGroupsEnabled}
        onSelectModeToggle={handleSelectModeToggle}
        onImport={handleImport}
        onBurstGroupToggle={handleBurstGroupToggle}
        onViewModeChange={(mode) => viewMode = mode}
        onSortModeChange={(mode) => sortMode = mode}
        onLayoutModeChange={handleLayoutModeChange}
    />
</div>

<script lang="ts">
    import LibraryPanel from "./_libraryPanel.svelte";
    import ReportPanel from "./_reportPanel.svelte";
    import FloatingToolbar, { type ViewMode, type SortMode } from "./_floatingToolbar.svelte";

    import { images } from "./mockData";
    import type { Image, ReportImage } from "./types";

    let allImages = $state<Image[]>(images);
    let reportImages = $state<ReportImage[]>([]);
    let reportedIds = $derived(new Set(reportImages.map((img) => img.id)));

    // Toolbar state
    let selectMode = $state(false);
    let viewMode = $state<ViewMode>("grid-compact");
    let sortMode = $state<SortMode>("newest");
    let burstGroupsEnabled = $state(false);
    let selectedIds = $state<Set<string>>(new Set());

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
        // TODO: Implement camera import
        console.log("Import from camera");
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

    function handleShowSelected(): void {
        // TODO: Show selected photos panel/sidebar
        console.log("Show selected photos");
    }
</script>

<div class="content p-6 pb-32">
    <!-- Two-Panel Layout -->
    <div class="grid grid-cols-2 gap-6">
        <!-- Left: Project Image Library -->
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

        <!-- Right: Included in Report -->
        <ReportPanel
            images={reportImages}
            onRemove={removeFromReport}
            onReorder={updateOrder}
            onCaptionChange={updateCaption}
        />
    </div>

    <!-- Floating Toolbar -->
    <FloatingToolbar
        {selectMode}
        {viewMode}
        {sortMode}
        hasBurstGroups={burstGroupsEnabled}
        onSelectModeToggle={handleSelectModeToggle}
        onImport={handleImport}
        onBurstGroupToggle={handleBurstGroupToggle}
        onViewModeChange={(mode) => viewMode = mode}
        onSortModeChange={(mode) => sortMode = mode}
        onShowSelected={handleShowSelected}
    />
</div>

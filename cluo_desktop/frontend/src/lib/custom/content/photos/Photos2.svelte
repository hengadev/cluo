<script lang="ts">
    import { Button } from "bits-ui";
    import { Camera } from "@lucide/svelte";
    import LibraryPanel from "./_libraryPanel.svelte";
    import ReportPanel from "./_reportPanel.svelte";

    import { images } from "./mockData";
    import type { Image, ReportImage } from "./types";

    let allImages = $state<Image[]>(images);
    let reportImages = $state<ReportImage[]>([]);
    let reportedIds = $derived(new Set(reportImages.map((img) => img.id)));

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
</script>

<div class="content p-6">
    <!-- Header Row -->
    <header class="flex justify-end items-center mb-6">
        <Button.Root
            class="gap-2 items-center h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex px-4 text-[15px] font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98]"
        >
            <Camera size={18} />
            <span>Importer depuis la caméra</span>
        </Button.Root>
    </header>

    <!-- Two-Panel Layout -->
    <div class="grid grid-cols-2 gap-6">
        <!-- Left: Project Image Library -->
        <LibraryPanel
            images={allImages}
            {reportedIds}
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
</div>

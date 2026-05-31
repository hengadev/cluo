<script lang="ts">
    import { ChevronRight, Square, Loader2 } from "@lucide/svelte";
    import type { RecordingStatus } from "$lib/types/recording";

    interface Props {
        id?: number | string;
        title?: string;
        date?: string;
        startTime?: string;
        duration?: string;
        status?: RecordingStatus;
    }

    let {
        id = 0,
        title = "Titre de l'enregistrement",
        date = "01 janv. 2025",
        startTime = "00:00",
        duration = "00:00",
        status = "completed",
    }: Props = $props();

    const statusLabels: Record<string, string> = {
        uploading: "Téléchargement",
        transcribing: "Transcription",
        analyzing: "Analyse",
        completed: "Terminé",
        failed: "Échoué",
    };

    const isProcessing = $derived(
        status === "uploading" ||
        status === "transcribing" ||
        status === "analyzing"
    );

    const statusColor = $derived(() => {
        switch (status) {
            case "completed":
                return "bg-green-500";
            case "failed":
                return "bg-red-500";
            default:
                return "bg-yellow-500";
        }
    });
</script>

<a
    href="/recording/{id}"
    class="flex justify-between border-1 border-black-50 rounded-input px-3 py-4 hover:bg-dark-50 transition-colors cursor-pointer no-underline"
>
    <div>
        <p class="text-dark-700 font-medium text-sm">{title}</p>
        <div class="flex gap-4 items-center">
            <p class="text-dark-600 text-xxs">{date}</p>
            <Square class="bg-dark-100" size={4} />
            <p class="text-dark-300 text-xxs">{startTime}</p>
            {#if isProcessing}
                <div class="flex items-center gap-1">
                    <Loader2 size={10} class="animate-spin text-yellow-600" />
                    <p class="text-dark-500 text-xxs">{statusLabels[status] ?? status}</p>
                </div>
            {/if}
        </div>
    </div>
    <div class="flex items-center gap-4">
        <div class="flex items-center gap-2">
            {#if isProcessing}
                <div class="w-2 h-2 rounded-full {statusColor} animate-pulse"></div>
            {/if}
            <p
                class="flex justify-center items-center border-1 border-dark-100 rounded-3xl bg-dark-50 text-dark-600 py-1 px-2 text-xs"
            >
                {duration}
            </p>
        </div>
        <ChevronRight />
    </div>
</a>

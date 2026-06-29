<script lang="ts">
    import type { RecordingPurpose, RecordingStatus } from "$lib/types/recording";

    interface Props {
        id?: number | string;
        title?: string;
        date?: string;
        startTime?: string;
        duration?: number | string;
        status?: RecordingStatus;
        purpose?: RecordingPurpose;
    }

    let {
        id = 0,
        title = "Titre de l'enregistrement",
        date = "01 janv. 2025",
        startTime = "00:00",
        duration = 0,
        status = "completed",
        purpose = "general",
    }: Props = $props();

    function formatDuration(d: number | string): string {
        const secs = typeof d === "number" ? d : 0;
        if (secs <= 0) return "--:--";
        const m = Math.floor(secs / 60);
        const s = Math.floor(secs % 60);
        return `${m.toString().padStart(2, "0")}:${s.toString().padStart(2, "0")}`;
    }

    const statusLabels: Record<string, string> = {
        uploading: "Téléchargement",
        transcribing: "Transcription",
        analyzing: "Analyse",
        completed: "Terminé",
        failed: "Échoué",
    };

    const isProcessing = $derived(
        status === "uploading" || status === "transcribing" || status === "analyzing"
    );

    const statusTextColor = $derived(
        status === "failed" ? "text-destructive" : "text-tertiary"
    );
</script>

<a
    href="/recording/{id}"
    class="flex items-center gap-4 py-4 border-b border-dark-100 last:border-b-0 hover:opacity-60 transition-opacity cursor-pointer no-underline"
>
    <div class="flex-1 min-w-0">
        <div class="flex items-center gap-2">
            <p class="text-dark-900 font-semibold text-[15px] truncate">{title}</p>
            {#if purpose === "witness_interview"}
                <span class="flex-shrink-0 text-xs font-medium px-2 py-0.5 rounded-full border border-dark-200 text-dark-600 bg-dark-50">Audition témoin</span>
            {/if}
        </div>
        <div class="flex items-center gap-1.5 mt-0.5">
            <p class="text-dark-400 text-sm">{date}{startTime ? " · " + startTime : ""}</p>
            {#if isProcessing || status === "failed"}
                <span class="text-dark-300 text-sm">·</span>
                <span class="{statusTextColor} text-sm flex items-center gap-1">
                    {#if isProcessing}
                        <span class="w-1.5 h-1.5 rounded-full bg-current animate-pulse inline-block flex-shrink-0"></span>
                    {/if}
                    {statusLabels[status] ?? status}
                </span>
            {/if}
        </div>
    </div>
    <p class="text-dark-400 text-sm font-mono flex-shrink-0">{formatDuration(duration)}</p>
</a>

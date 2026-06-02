<script lang="ts">
    import { ChevronRight } from "@lucide/svelte";
    import type { RecordingStatus } from "$lib/types/recording";

    interface Props {
        id?: number | string;
        title?: string;
        date?: string;
        startTime?: string;
        duration?: number | string;
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
        status === "uploading" || status === "transcribing" || status === "analyzing"
    );

    const avatarBg = $derived(
        status === "completed" ? "bg-green-50" :
        status === "failed" ? "bg-red-50" :
        "bg-amber-50"
    );

    const barColor = $derived(
        status === "completed" ? "bg-green-400" :
        status === "failed" ? "bg-red-400" :
        "bg-amber-400"
    );

    const statusTextColor = $derived(
        status === "completed" ? "text-green-600" :
        status === "failed" ? "text-red-500" :
        "text-amber-600"
    );
</script>

<a
    href="/recording/{id}"
    class="flex items-center gap-3 border border-dark-100 rounded-card-sm px-3 py-3 hover:bg-dark-50 transition-colors cursor-pointer no-underline"
>
    <!-- Waveform avatar -->
    <div class="flex-shrink-0 w-10 h-10 rounded-card-sm {avatarBg} flex items-center justify-center gap-[3px]">
        <div class="w-[3px] rounded-full {barColor} h-2 {isProcessing ? 'animate-pulse' : ''}"></div>
        <div class="w-[3px] rounded-full {barColor} h-4 {isProcessing ? 'animate-pulse' : ''}"></div>
        <div class="w-[3px] rounded-full {barColor} h-6 {isProcessing ? 'animate-pulse' : ''}"></div>
        <div class="w-[3px] rounded-full {barColor} h-3 {isProcessing ? 'animate-pulse' : ''}"></div>
        <div class="w-[3px] rounded-full {barColor} h-5 {isProcessing ? 'animate-pulse' : ''}"></div>
    </div>

    <!-- Content -->
    <div class="flex-1 min-w-0">
        <p class="text-dark-700 font-medium text-sm truncate">{title}</p>
        <div class="flex gap-1.5 items-center mt-0.5">
            <p class="text-dark-500 text-xxs">{date}</p>
            <span class="text-dark-300 text-xxs">·</span>
            <p class="text-dark-400 text-xxs">{startTime}</p>
        </div>
    </div>

    <!-- Duration + status -->
    <div class="flex flex-col items-end gap-0.5 flex-shrink-0">
        <p class="text-dark-700 font-mono text-xs font-medium">{duration}</p>
        <span class="text-xxs {statusTextColor} flex items-center gap-1">
            {#if isProcessing}
                <span class="w-1.5 h-1.5 rounded-full bg-current animate-pulse inline-block"></span>
            {/if}
            {statusLabels[status] ?? status}
        </span>
    </div>

    <ChevronRight size={16} class="text-dark-300 flex-shrink-0" />
</a>

<script lang="ts">
    import { ArrowLeft, Play, Trash2, Download } from "@lucide/svelte";

    // Data is passed from +page.server.ts load function
    let { data } = $props();
    const recording = data.recording;
</script>

<div class="min-h-screen flex flex-col gap-6 pb-24">
    <!-- Header with back button -->
    <div class="flex items-center gap-3">
        <a
            href="/"
            class="flex items-center justify-center w-10 h-10 rounded-full hover:bg-dark-50 transition-colors"
        >
            <ArrowLeft class="text-dark-700" />
        </a>
        <h1 class="text-dark-900 font-extrabold text-xl">Recording Details</h1>
    </div>

    <!-- Recording Info Card -->
    <div class="flex flex-col gap-4 p-4 border-1 border-dark-100 rounded-2xl">
        <div class="flex justify-between items-start">
            <div>
                <h2 class="text-dark-800 font-bold text-lg">
                    {recording.title}
                </h2>
                <div class="flex gap-2 items-center mt-1">
                    <p class="text-dark-600 text-sm">{recording.date}</p>
                    <span class="text-dark-300">•</span>
                    <p class="text-dark-400 text-sm">{recording.startTime}</p>
                </div>
            </div>
            <div class="flex items-center gap-2">
                <p
                    class="flex justify-center items-center border-1 border-dark-100 rounded-3xl bg-dark-50 text-dark-600 py-1 px-3 text-sm font-medium"
                >
                    {recording.duration}
                </p>
            </div>
        </div>

        <!-- Audio Player Placeholder -->
        <div
            class="flex items-center gap-4 p-4 bg-dark-50 rounded-xl border-1 border-dark-100"
        >
            <button
                class="flex items-center justify-center w-12 h-12 bg-dark-700 rounded-full hover:bg-dark-600 transition-colors"
            >
                <Play class="text-foreground" fill="currentColor" size={20} />
            </button>
            <div class="flex-1 h-2 bg-dark-200 rounded-full">
                <!-- Progress bar would go here -->
                <div class="h-full w-0 bg-dark-700 rounded-full"></div>
            </div>
        </div>
    </div>

    <!-- Tags -->
    {#if recording.tags && recording.tags.length > 0}
        <div class="flex flex-col gap-2">
            <p class="text-dark-700 font-bold text-base">Tags</p>
            <div class="flex gap-2 flex-wrap">
                {#each recording.tags as tag}
                    <span
                        class="px-3 py-1 bg-accent text-accent-foreground rounded-full text-sm"
                    >
                        {tag}
                    </span>
                {/each}
            </div>
        </div>
    {/if}

    <!-- Transcript Section -->
    <div class="flex flex-col gap-3">
        <p class="text-dark-700 font-bold text-base">Transcript</p>
        <div
            class="p-4 border-1 border-dark-100 rounded-2xl bg-background-alt min-h-40"
        >
            <p class="text-dark-700 text-sm leading-relaxed">
                {recording.transcript}
            </p>
        </div>
    </div>

    <!-- Action Buttons -->
    <div class="flex gap-3 mt-4">
        <button
            class="flex-1 flex items-center justify-center gap-2 px-4 py-3 bg-dark-700 hover:bg-dark-600 text-foreground rounded-xl transition-colors"
        >
            <Download size={18} />
            <span class="font-medium">Download</span>
        </button>
        <button
            class="flex-1 flex items-center justify-center gap-2 px-4 py-3 bg-destructive hover:bg-destructive/90 text-white rounded-xl transition-colors"
        >
            <Trash2 size={18} />
            <span class="font-medium">Delete</span>
        </button>
    </div>
</div>

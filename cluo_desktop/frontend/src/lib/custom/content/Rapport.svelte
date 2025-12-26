<script lang="ts">
    // this part should be so that I can decide on the size of each split
    // with a maximum width for the one with AI.
    import {
        MessageSquareMore,
        Lightbulb,
        PencilLine,
        AudioLines,
    } from "@lucide/svelte";

    type AIButton = {
        icon: typeof import("@lucide/svelte").Icon;
        title: string;
    };

    const aiButtons: AIButton[] = [
        { icon: MessageSquareMore, title: "Chat" },
        { icon: Lightbulb, title: "Ideas" },
        { icon: PencilLine, title: "Review" },
        { icon: AudioLines, title: "Audio" },
    ];
    // TODO: make sure that the row 1 just its size and the row takes the remaining space.
</script>

<div class="grid grid-cols-2 h-full">
    <div class="bg-red-300">rapport the part with the text editor</div>
    <div class="grid grid-rows-[auto_1fr]">
        <div
            class="flex justify-center gap-6 border-b-1 border-dark-50 text-center"
        >
            <!-- AI buttons toolbar will go here -->
        </div>
    </div>
</div>

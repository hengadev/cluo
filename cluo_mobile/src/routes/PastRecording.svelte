<script lang="ts">
    import { ChevronRight, Square } from "@lucide/svelte";

    interface Props {
        id?: number | string;
        title?: string;
        date?: string;
        startTime?: string;
        duration?: string;
    }

    let {
        id = 0,
        title = "The title of the recording",
        date = "01 Jan, 2025",
        startTime = "00:00",
        duration = "00:00",
    }: Props = $props();
</script>

<a
    href="/recording/{id}"
    class="flex justify-between border-1 border-black-50 rounded-input px-3 py-4 hover:bg-dark-50 transition-colors cursor-pointer no-underline"
>
    <div>
        <p class="text-dark-700 font-medium text-sm">{title}</p>
        <div class="flex gap-4">
            <p class="text-dark-600 text-xxs">{date}</p>
            <Square class="bg-dark-100" size={4} />
            <p class="text-dark-300 text-xxs">{startTime}</p>
        </div>
    </div>
    <div class="flex items-center gap-4">
        <p
            class="flex justify-center items-center border-1 border-dark-100 rounded-3xl bg-dark-50 text-dark-600 py-1 px-2 text-xs"
        >
            {duration}
        </p>
        <ChevronRight />
    </div>
</a>

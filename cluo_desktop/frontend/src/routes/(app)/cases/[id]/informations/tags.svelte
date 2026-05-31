<script lang="ts">
    import { FileText, Mic, Image, UserRoundSearch } from "@lucide/svelte";

    // TODO: I need to implement this later
    // interface Props {}
    // let {}: Props = $props();

    interface Tag {
        icon: typeof import("@lucide/svelte").Icon;
        count: number;
        title: string;
    }
    const tags: Tag[] = [
        { icon: FileText, count: 15, title: "Documents" },
        { icon: Mic, count: 12, title: "Audios" },
        { icon: Image, count: 3, title: "Image" },
        { icon: UserRoundSearch, count: 4, title: "Personnes" },
    ];
</script>

<div class="grid grid-cols-4 gap-32 justify-between items-center">
    {#each tags as tag}
        {@const Icon = tag.icon}
        <div
            class="flex gap-8 items-center px-6 py-3 rounded-input border-1 border-dark-50"
        >
            <div class="p-3 bg-dark-50 rounded-input">
                <Icon size={24} />
            </div>
            <div class="grid">
                <p class="text-dark-900 font-extrabold text-2xl">{tag.count}</p>
                <p class="text-dark-700 text-xl">{tag.title}</p>
            </div>
        </div>
    {/each}
</div>

<script lang="ts">
    import { Tabs, Button } from "bits-ui";
    import {
        Upload,
        Grid,
        LayoutList,
        ArrowUpNarrowWide,
    } from "@lucide/svelte";
    import AllPhotos from "./allPhotos.svelte";
</script>

<div class="content p-6">
    <Tabs.Root value="all" class="rounded-card w-full p-3">
        <div class="flex justify-between items-center">
            <Tabs.List
                class="rounded-9px bg-dark-10 shadow-mini-inset dark:bg-background grid w-full grid-cols-2 gap-1 p-1 text-sm font-semibold leading-[0.01em] dark:border dark:border-neutral-600/30 max-w-[400px]"
            >
                <Tabs.Trigger
                    value="all"
                    class="data-[state=active]:shadow-mini dark:data-[state=active]:bg-muted h-8 rounded-[7px] bg-transparent py-2 data-[state=active]:bg-white"
                    >Toutes les photos</Tabs.Trigger
                >
                <Tabs.Trigger
                    value="selection"
                    class="data-[state=active]:shadow-mini dark:data-[state=active]:bg-muted h-8 rounded-[7px] bg-transparent py-2 data-[state=active]:bg-white"
                    >Photos selectionnees</Tabs.Trigger
                >
            </Tabs.List>
            <Button.Root
                class="gap-2 items-center h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex px-4 text-[15px] font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98]"
            >
                <Upload />
                <p>Ajouter une photo</p>
            </Button.Root>
        </div>
        <Tabs.Content value="all" class="select-none pt-3">
            <AllPhotos />
        </Tabs.Content>
        <Tabs.Content value="selection" class="select-none pt-3">
            <div>here is the content with the selection</div>
        </Tabs.Content>
    </Tabs.Root>
</div>

<script lang="ts">
    import { Slack, User } from "@lucide/svelte";
    import { Button, Tooltip } from "bits-ui";
    import { items, type SidebarItem } from "$lib/constructor/sidebar";
    let selected: string = $state(items[0].title);
</script>

<div
    class="grid-area h-full p-2 pt-8 flex flex-col gap-12 items-center bg-[#fafafa] border-2 border-[#e5e7eb]"
>
    <div
        class="border-dark border-2 mt-24 size-16 rounded-input flex items-center justify-center"
    >
        <Slack size={36} strokeWidth={1.5} />
    </div>
    <div class="flex flex-col justify-between h-full">
        <div class="flex flex-col items-center gap-4">
            {#each items as item}
                {@render button(item)}
            {/each}
        </div>
        <Button.Root
            class="rounded-full flex items-center justify-center border-1 border-[#e5e7eb] mx-auto size-12 bg-muted"
        >
            <User size={32} />
        </Button.Root>
    </div>
</div>

{#snippet button(item: SidebarItem)}
    <Tooltip.Provider>
        <Tooltip.Root delayDuration={100}>
            <Tooltip.Trigger
                class="align-center border-border-input rounded-10px bg-background-alt ring-offset-background active:scale-[0.98] active:transition:all 
		focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex size-12 items-center justify-center focus-visible:ring-2 focus-visible:ring-offset-2 {item.title ===
                selected
                    ? 'bg-dark text-white'
                    : 'bg-transparent text-black hover:bg-[#f0f0f0]'}"
                onclick={() => {
                    item.fn;
                    selected = item.title;
                }}
            >
                {@const Icon = item.icon}
                <Button.Root class="cursor-pointer">
                    <Icon size={32} strokeWidth={1.5} />
                </Button.Root>
            </Tooltip.Trigger>
            <Tooltip.Content sideOffset={8} side="right">
                <div
                    class="rounded-input text-[1rem] align-center bg-dark text-white font-medium gap-3 border-dark-10 shadow-popover outline-hidden z-0 flex items-center justify-center border p-2"
                >
                    {item.title}
                </div>
            </Tooltip.Content>
        </Tooltip.Root>
    </Tooltip.Provider>
{/snippet}

<style>
    .grid-area {
        grid-area: sidebar;
    }
</style>

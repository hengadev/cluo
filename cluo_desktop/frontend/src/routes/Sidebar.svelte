<script lang="ts">
    import { Button, Tooltip } from "bits-ui";
    import { items, type SidebarItem } from "$lib/constructor/sidebar";
    let selected: string = $state(items[0].title);
</script>

<div
    class="grid-area h-full px-2 py-2 relative flex flex-col gap-4 bg-[#fafafa] border-2 border-[#e5e7eb]"
>
    {#each items as item}
        {@render button(item)}
    {/each}
</div>

{#snippet button(item: SidebarItem)}
    <Tooltip.Provider>
        <Tooltip.Root delayDuration={100}>
            <Tooltip.Trigger
                class="border-border-input rounded-10px p-2 bg-background-alt ring-offset-background active:scale-[0.98] active:transition:all 
		focus-visible:ring-dark cursor-pointer focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex size-10 items-center justify-center focus-visible:ring-2 focus-visible:ring-offset-2 {item.title ===
                selected
                    ? 'bg-dark text-white'
                    : 'bg-transparent text-black'}"
            >
                {@const Icon = item.icon}

                <Button.Root
                    onclick={() => {
                        item.fn;
                        selected = item.title;
                    }}
                >
                    <Icon size={32} />
                </Button.Root>
            </Tooltip.Trigger>
            <Tooltip.Content sideOffset={8} side="right">
                <div
                    class="rounded-input text-[1rem] align-center bg-dark text-white font-semibold gap-3 border-dark-10 shadow-popover outline-hidden z-0 flex items-center justify-center border p-2"
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

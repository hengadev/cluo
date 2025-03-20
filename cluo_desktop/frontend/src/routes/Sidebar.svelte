<script lang="ts">
    import { Button } from "bits-ui";
    import { items, type SidebarItem } from "$lib/constructor/sidebar";
    import { PanelLeftClose, PanelLeftOpen } from "@lucide/svelte";
    let collapsed: boolean = $state(false);
    let selected = $state<SidebarItem>();
    // TODO: add dialog thing in this
    // TODO: add transiton on width change
    // I think I need to use a ref for this
</script>

<div class="grid-area h-full px-4 py-2 relative flex flex-col gap-2">
    <Button.Root
        class="absolute right-[1rem] bottom-[1rem]"
        onclick={() => (collapsed = !collapsed)}
    >
        {#if collapsed}
            <PanelLeftOpen />
        {:else}
            <PanelLeftClose />
        {/if}
    </Button.Root>
    {#each items as item}
        {@render button(item)}
    {/each}
</div>

{#snippet button(item: SidebarItem)}
    <Button.Root
        onclick={() => {
            item.fn;
            selected = item;
        }}
        class="rounded-input text-[1rem] flex shadow-mini align-center bg-dark hover:bg-dark/95 text-background font-semibold active:scale-[0.98] active:transition:all p-4 gap-3"
    >
        {@const Icon = item.icon}
        <Icon />
        {#if !collapsed}
            <p>{item.title}</p>
        {/if}
    </Button.Root>
{/snippet}

<style>
    .grid-area {
        grid-area: sidebar;
    }
</style>

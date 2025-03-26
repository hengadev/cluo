<script lang="ts">
    import { Folder } from "@lucide/svelte";
    import Search from "$lib/custom/header/Search.svelte";
    import { items, type HeaderItem } from "$lib/constructor/header";
    // TODO: pour la partie client ou type d'enquete, il faut un select avec tous les anciens clients + un bouton plus pour ajouter un nouvel element
    import { Button, Tooltip } from "bits-ui";
</script>

<div class="header bg-[#fafafa] border-1 border-[#e5e7eb]">
    <div class="grid">
        <p class="text-base font-semibold">Website performance issue</p>
        <div class="left">
            <div class="current-case">
                <div class="bg-yellow-700 p-2 rounded-input">
                    <Folder size={16} color="white" />
                </div>
                <p>#CS-1234</p>
            </div>
            <p>&bull;</p>
            <p>Jean DUPONT</p>
        </div>
    </div>
    <Search />
    <div class="flex align-center gap-2">
        <div class="buttons">
            {#each items as item}
                {@const DialogOrPopover = item.uiComponent}
                <DialogOrPopover>
                    {@render headerItem(item)}
                </DialogOrPopover>
            {/each}
        </div>
    </div>
</div>

{#snippet headerItem(item: HeaderItem)}
    {@const Icon = item.icon}
    <Tooltip.Provider>
        <Tooltip.Root delayDuration={100}>
            <Tooltip.Trigger
                class="border-border-input border-1 rounded-10px p-2 bg-background-alt ring-offset-background active:scale-[0.98] active:transition:all 
		focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex size-14 items-center justify-center focus-visible:ring-2 focus-visible:ring-offset-2 
                    hover:bg-[#f0f0f0]'} {item.bg} text-{item.fg}"
            >
                <Button.Root class="cursor-pointer">
                    <Icon size={32} strokeWidth={1.5} />
                </Button.Root>
            </Tooltip.Trigger>
            <Tooltip.Content sideOffset={8} side="bottom">
                <div
                    class="rounded-input text-[1rem] align-center bg-dark text-white font-medium border-dark-10 shadow-popover outline-hidden z-0 flex items-center justify-center border p-2"
                >
                    {item.title}
                </div>
            </Tooltip.Content>
        </Tooltip.Root>
    </Tooltip.Provider>
{/snippet}

<style>
    .header {
        grid-area: header;
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 0.5rem 2rem;
        gap: 2rem;
    }
    .left,
    .buttons {
        flex: 1;
    }
    .left {
        display: flex;
        gap: 0.5rem;
        align-items: center;
        font-weight: 500;
    }
    .current-case {
        display: flex;
        gap: 0.5rem;
        align-items: center;
    }
    .buttons {
        display: flex;
        justify-content: right;
        gap: 0.5rem;
        margin-left: auto;
        text-align: right;
    }
</style>

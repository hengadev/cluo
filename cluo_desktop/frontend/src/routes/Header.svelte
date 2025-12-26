<script lang="ts">
    import { Folder } from "@lucide/svelte";
    import Search from "$lib/custom/header/Search.svelte";
    import { items, type HeaderItem } from "$lib/constructor/header";
    import ThemeToggle from "$lib/components/ThemeToggle.svelte";
    // TODO: pour la partie client ou type d'enquete, il faut un select avec tous les anciens clients + un bouton plus pour ajouter un nouvel element
    import { Button, Tooltip } from "bits-ui";
</script>

<div class="header border-1 border-dark-50">
    <div class="grid">
        <div class="left">
            <div class="current-case">
                <div class="p-2 rounded-input bg-foreground">
                    <Folder size={16} class="text-background" />
                </div>
                <p>Cabinet DUPONT</p>
            </div>
            <p>&bull;</p>
            <p>Affaire Ti Sonson</p>
        </div>
    </div>
    <Search />
    <div class="flex align-center gap-2">
        <div class="buttons">
            <ThemeToggle />
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
                class="rounded-10px p-3 bg-dark-50 ring-offset-background active:scale-[0.98] active:transition:all 
		focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center focus-visible:ring-2 focus-visible:ring-offset-2 
                    hover:bg-dark-100/50 {item.bg} text-{item.fg}"
            >
                <Button.Root class="cursor-pointer">
                    <Icon size={24} strokeWidth={1.75} />
                </Button.Root>
            </Tooltip.Trigger>
            <Tooltip.Content sideOffset={8} side="bottom">
                <div
                    class="rounded-input text-[1rem] align-center bg-foreground text-background font-medium border-dark-10 shadow-popover outline-hidden z-0 flex items-center justify-center border p-2"
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

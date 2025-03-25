<script lang="ts">
    import { Slack, Folder, MessageSquareMore, Settings } from "@lucide/svelte";
    // import { Input } from "$lib/components/ui/input";
    // import { Button } from "$lib/components/ui/button";
    import NewCase from "$lib/custom/NewCase.svelte";
    // TODO: pour la partie client ou type d'enquete, il faut un select avec tous les anciens clients + un bouton plus pour ajouter un nouvel element
    import { Button, Tooltip } from "bits-ui";
</script>

{#snippet headerOption(Icon: any, clickAction: any)}
    <button onclick={clickAction}>
        <Icon />
    </button>
{/snippet}

<div class="header">
    <div class="left">
        <Slack size={36} />
        <div class="current-case">
            <Folder size={24} />
            <p>HENRY Gary</p>
        </div>
    </div>
    <input class="search" />
    <!-- <Input /> -->
    <div class="buttons">
        <NewCase />
        {@render headerOption(MessageSquareMore, function (e: MouseEvent) {
            console.log("here is the event:", e);
        })}
        <button>Button</button>
        <Settings />
    </div>
</div>

{#snippet buttons(item: HeaderItem)}
    {@const Icon = item.icon}
    <Tooltip.Provider>
        <Tooltip.Root delayDuration={100}>
            <Tooltip.Trigger
                class="border-border-input border-1 rounded-10px p-2 bg-background-alt ring-offset-background active:scale-[0.98] active:transition:all 
		focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex size-16 items-center justify-center focus-visible:ring-2 focus-visible:ring-offset-2 
                    hover:bg-[#f0f0f0]'} {item.bg} text-{item.fg}"
            >
                <Button.Root onclick={item.fn} class="cursor-pointer">
                    <Icon size={32} />
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
        gap: 4rem;
    }
    .current-case {
        display: flex;
        gap: 1rem;
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

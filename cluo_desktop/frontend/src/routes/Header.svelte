<script lang="ts">
    import { Slack, Folder, MessageSquareMore, Settings } from "@lucide/svelte";
    // import { Input } from "$lib/components/ui/input";
    // import { Button } from "$lib/components/ui/button";
    import NewCase from "$lib/custom/NewCase.svelte";
    // TODO: pour la partie client ou type d'enquete, il faut un select avec tous les anciens clients + un bouton plus pour ajouter un nouvel element
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

<style>
    .header {
        grid-area: header;
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 1rem 2rem;
        gap: 2rem;
    }
    .search {
        width: 100%;
        max-width: 800px;
        flex: 0 0 auto;
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

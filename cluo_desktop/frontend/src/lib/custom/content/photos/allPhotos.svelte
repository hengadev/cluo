<script lang="ts">
    import { type LayoutState, LAYOUT_STATES } from "./types";

    import { type CardType } from "./types";
    import CardGrid from "./_cardGrid.svelte";
    import CardList from "./_cardList.svelte";

    let viewMode = $state<LayoutState>(LAYOUT_STATES.Grid);
    // let viewMode = $state<LayoutState>(LAYOUT_STATES.List);

    const cards: CardType[] = [
        { title: "Client contract", tags: ["contract", "signed"] },
        { title: "Evidence photo", tags: ["contract", "signed"] },
        { title: "Medical report", tags: ["contract", "signed"] },
    ];
</script>

<div class="grid">
    {#if viewMode === LAYOUT_STATES.Grid}
        {#each cards as card}
            <CardGrid {card} />
        {/each}
    {:else if viewMode === LAYOUT_STATES.List}
        {#each cards as card}
            <CardList {card} />
        {/each}
    {/if}
</div>

<script lang="ts">
    import { type LayoutState, type Image, LAYOUT_STATES } from "./types";

    import CardGrid from "./_cardGrid.svelte";
    import CardList from "./_cardList.svelte";

    interface Props {
        layoutValue: LayoutState;
        images: Image[];
    }

    let { layoutValue, images }: Props = $props();
</script>

<div class="pt-8 grid grid-cols-6 gap-4">
    {#each images as image}
        {#if layoutValue === LAYOUT_STATES.Grid}
            <CardGrid card={image} />
        {:else if layoutValue === LAYOUT_STATES.List}
            <CardList card={image} />
        {/if}
    {/each}
</div>

<script lang="ts">
    import { type CardType } from "./types";
    type Props = {
        card: CardType;
    };
    let { card }: Props = $props();
</script>

<div class="border-border-input border-1">
    <div class="image"></div>
    <div class="p-2">
        <p class="title text-base font-semibold">
            this is the card list{card.title}
        </p>
        <div class="tags p-2 flex gap-2">
            {#each card.tags as tag}
                <p class="tag text-xs rounded-xl">{tag}</p>
            {/each}
        </div>
    </div>
</div>

<style>
    .image {
        background-color: #fafafa;
        height: 200px;
        width: 200px;
    }
    .tag {
        background-color: #dedede;
        padding: 0.5rem 1rem;
        max-width: 100px;
    }
</style>

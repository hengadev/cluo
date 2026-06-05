<script lang="ts">
    import { type CardType } from "./types";
    import { Check } from "@lucide/svelte";
    type Props = {
        card: CardType;
    };
    let { card }: Props = $props();

    import { Button } from "bits-ui";
    let selected = $state(false);
</script>

<Button.Root
    class="border-border-input border-1 max-w-120 rounded-input relative {selected
        ? `border border-red-500`
        : `border-0`} p-4"
    onclick={() => (selected = !selected)}
>
    <div>here goes the caption if there is any</div>
    <div class="aspect-square bg-dark-50 w-full overflow-hidden"></div>
    <Button.Root
        class="absolute top-4 right-4 w-6 aspect-square bg-background"
        onclick={() => (selected = !selected)}
    >
        {#if selected}
            <Check />
        {/if}
    </Button.Root>
    <div>here is some other text</div>
</Button.Root>

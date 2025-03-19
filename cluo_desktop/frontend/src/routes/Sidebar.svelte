<script lang="ts">
    // import { Button } from "$lib/components/ui/button";
    import { elements } from "$lib/constructor/sidebar";
    import { PanelLeftClose, PanelLeftOpen } from "@lucide/svelte";
    let collapsed: boolean = $state(true);
    // let active: string = $state("Informations");
    let active: string = $state(elements[0].title);
    // TODO: add dialog thing in this
</script>

{#snippet sidebarOption(Icon: any, title: string)}
    <button
        class="btn"
        class:active={title === active}
        onclick={() => (active = title)}
    >
        <Icon />
        {#if collapsed}
            <p class="title">{title}</p>
        {/if}
    </button>
{/snippet}

<div class="sidebar">
    <button class="collapser" onclick={() => (collapsed = !collapsed)}>
        {#if collapsed}
            <PanelLeftClose />
        {:else}
            <PanelLeftOpen />
        {/if}
    </button>
    {#each elements as element}
        {@render sidebarOption(element.image, element.title)}
    {/each}
</div>

<style>
    .sidebar {
        grid-area: sidebar;
        border: 2px solid red;
        height: 100%;
        padding: 1rem 0.5rem;
        position: relative;
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }
    .btn {
        /* border: 1px solid pink; */
        display: flex;
        align-items: center;
        gap: 1rem;
        padding: 0.5rem 1rem;
    }
    .btn:is(:hover, :focus) {
        background-color: #f5f4f4;
    }
    .active {
        background-color: #f5f4f4;
    }
    .collapser {
        position: absolute;
        right: 1rem;
        bottom: 1rem;
    }
</style>

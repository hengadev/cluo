<script lang="ts">
    import { onMount } from "svelte";
    import { Home, User, ChevronRight, ChevronLeft } from "@lucide/svelte";
    import { Button, Tooltip } from "bits-ui";
    import ProfilePopover from "$lib/custom/sidebar/ProfilePopover.svelte";
    import { items, type SidebarItem } from "$lib/constructor/sidebar";
    import { currentCase } from "$lib/stores/case";
    import { goto } from "$app/navigation";
    import { page } from "$app/stores";
    import NoCaseDialog from "$lib/custom/global/NoCaseDialog.svelte";
    import NewCase from "$lib/custom/header/NewCase.svelte";

    const size: number = 24;

    let noCaseOpen: boolean = $state(false);
    let newCaseOpen: boolean = $state(false);

    onMount(() => {
        if (!$currentCase.id) {
            noCaseOpen = true;
        }
    });

    let isExpanded: boolean = $state(false);

    const regularItems = items.filter(i => !i.disabled);
    const disabledItems = items.filter(i => i.disabled);

    // Get the current route path for highlighting
    function getRouteForItem(item: SidebarItem): string {
        // Routes without :id are used as-is
        if (!item.path.includes(':id')) return item.path;
        const caseId = $currentCase.id;
        if (!caseId) return '';
        return item.path.replace(':id', caseId);
    }

    function isActive(item: SidebarItem): boolean {
        const routePath = getRouteForItem(item);
        if (!routePath) return false;
        return $page.url.pathname === routePath;
    }

    function handleItemClick(item: SidebarItem) {
        const routePath = getRouteForItem(item);
        if (!routePath) {
            noCaseOpen = true;
            return;
        }
        goto(routePath);
    }
</script>

<div
    class="grid-area-sidebar h-full p-1 pt-2 flex flex-col gap-10 bg-background-alt border-1 border-dark-50 relative transition-all duration-300 animate-fade-in"
    style="animation-delay: 200ms;"
    style:width={isExpanded ? '200px' : 'auto'}
    style:align-items={isExpanded ? 'stretch' : 'center'}
>
    <div class="flex {isExpanded ? 'justify-between' : 'justify-center'} items-center w-full">
        <Button.Root
            class="p-2 !mt-1 rounded-input flex items-center cursor-pointer transition-all duration-300 bg-transparent text-foreground hover:bg-foreground/10 {isExpanded
                ? 'justify-start gap-3 px-4'
                : 'justify-center'}"
            onclick={() => goto("/")}
        >
            <Home {size} strokeWidth={1.5} />
            {#if isExpanded}
                <span class="text-sm font-medium">Home</span>
            {/if}
        </Button.Root>

        <!-- Chevron toggle button -->
        <button
            onclick={() => (isExpanded = !isExpanded)}
            class="p-2 rounded-input text-dark-500 hover:bg-foreground/10 hover:scale-110 active:scale-95 transition-all duration-200 {isExpanded ? '' : '!mt-1'}"
            title={isExpanded ? "Collapse sidebar" : "Expand sidebar"}
            type="button"
        >
            {#if isExpanded}
                <ChevronLeft size={24} strokeWidth={1.5} />
            {:else}
                <ChevronRight size={24} strokeWidth={1.5} />
            {/if}
        </button>
    </div>
    <div class="flex flex-col justify-between h-full">
        <div
            class="flex flex-col gap-2"
            style:align-items={isExpanded ? 'stretch' : 'center'}
        >
            {#each regularItems as item}
                {@render button(item)}
            {/each}
        </div>
        <div class="flex flex-col gap-2 mb-2" style:align-items={isExpanded ? 'stretch' : 'center'}>
            {#each disabledItems as item}
                {@render button(item)}
            {/each}
        </div>
        <ProfilePopover>
            <Button.Root
                class="rounded-10px flex items-center border-1 border-border-input bg-background cursor-pointer transition-all duration-300 {isExpanded
                    ? 'justify-start gap-3 px-4 py-3 w-full'
                    : 'justify-center mx-auto size-12'}"
            >
                <User {size} />
                {#if isExpanded}
                    <span class="text-sm font-medium">Profile</span>
                {/if}
            </Button.Root>
        </ProfilePopover>
    </div>
</div>

<NoCaseDialog bind:open={noCaseOpen} onCreateCase={() => (newCaseOpen = true)} />
<NewCase bind:open={newCaseOpen} />

{#snippet button(item: SidebarItem)}
    {@const Icon = item.icon}
    {@const active = isActive(item)}
    {@const disabled = item.disabled ?? false}
    {#if isExpanded}
        <button
            class="align-center border-border-input rounded-10px bg-background-alt ring-offset-background active:scale-[0.98] active:transition:all
		focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden flex items-center gap-3 px-4 py-3 focus-visible:ring-2 focus-visible:ring-offset-2 {active
                ? 'bg-foreground text-background'
                : 'bg-transparent text-foreground hover:bg-foreground/10'} {disabled ? 'opacity-35 cursor-not-allowed' : ''}"
            onclick={disabled ? undefined : () => handleItemClick(item)}
            disabled={disabled}
        >
            <Icon size={24} strokeWidth={1.75} />
            <span class="text-sm font-medium whitespace-nowrap"
                >{item.title}</span
            >
        </button>
    {:else}
        <Tooltip.Provider>
            <Tooltip.Root delayDuration={300}>
                <Tooltip.Trigger
                    class="align-center border-border-input rounded-10px bg-background-alt ring-offset-background
		focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex size-12 items-center justify-center focus-visible:ring-2 focus-visible:ring-offset-2 {active
                        ? 'bg-foreground text-background'
                        : 'bg-transparent text-foreground hover:bg-foreground/10'} {disabled ? 'opacity-35 cursor-not-allowed' : 'active:scale-[0.98] active:transition:all'}"
                    onclick={disabled ? undefined : () => handleItemClick(item)}
                    disabled={disabled}
                >
                    <Button.Root class={disabled ? 'cursor-not-allowed pointer-events-none' : 'cursor-pointer'}>
                        <Icon size={24} strokeWidth={1.75} />
                    </Button.Root>
                </Tooltip.Trigger>
                <Tooltip.Content sideOffset={8} side="right">
                    <div
                        class="rounded-input text-[1rem] align-center bg-foreground text-background font-medium gap-3 border-dark-10 shadow-popover outline-hidden z-0 flex items-center justify-center border p-2"
                    >
                        {item.title}
                    </div>
                </Tooltip.Content>
            </Tooltip.Root>
        </Tooltip.Provider>
    {/if}
{/snippet}

<style>
    .grid-area-sidebar {
        grid-area: sidebar;
    }
</style>

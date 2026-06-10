<script lang="ts">
    import { Home, User, ChevronRight, ChevronLeft } from "@lucide/svelte";
    import { Button, Tooltip } from "bits-ui";
    import ProfilePopover from "$lib/custom/sidebar/ProfilePopover.svelte";
    import { groups, type SidebarItem } from "$lib/constructor/sidebar";
    import { currentCase } from "$lib/stores/case";
    import { goto } from "$app/navigation";
    import { page } from "$app/stores";

    let isExpanded: boolean = $state(false);

    let noCaseOpen = $derived(!$currentCase.id);

    function getRouteForItem(item: SidebarItem): string {
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
        if (!routePath) return;
        goto(routePath);
    }
</script>

<div
    class="grid-area-sidebar h-full p-1 pt-2 flex flex-col gap-10 bg-background-alt border border-border-input relative transition-all duration-300 animate-fade-in overflow-visible"
    style="animation-delay: 200ms;"
    style:width={isExpanded ? '200px' : 'auto'}
    style:align-items={isExpanded ? 'stretch' : 'center'}
>
    <!-- Home row with inline collapse toggle -->
    <div class="flex {isExpanded ? 'justify-between' : 'justify-center'} items-center w-full">
        <Button.Root
            class="p-2 !mt-1 rounded-input flex items-center cursor-pointer transition-all duration-300 bg-transparent text-foreground hover:bg-foreground/10 {isExpanded
                ? 'justify-start gap-3 px-4'
                : 'justify-center'}"
            onclick={() => goto("/")}
        >
            <Home size={24} strokeWidth={1.5} />
            {#if isExpanded}
                <span class="text-sm font-medium">Home</span>
            {/if}
        </Button.Root>
        <button
            onclick={() => (isExpanded = !isExpanded)}
            class="p-1.5 rounded-full text-muted-foreground hover:text-foreground hover:bg-foreground/10 active:scale-95 transition-all duration-200 cursor-pointer {isExpanded ? 'mr-2' : ''}"
            title={isExpanded ? "Réduire" : "Agrandir"}
            type="button"
        >
            {#if isExpanded}
                <ChevronLeft size={18} strokeWidth={2} />
            {:else}
                <ChevronRight size={18} strokeWidth={2} />
            {/if}
        </button>
    </div>

    <div class="flex flex-col justify-between h-full">
        <div class="flex flex-col gap-6" style:align-items={isExpanded ? 'stretch' : 'center'}>
            {#each groups as group, i}
                {#if i > 0}
                    <div class="flex flex-col {isExpanded ? 'gap-2 px-2' : 'gap-1.5 items-center'}">
                        {#if group.label}
                            <span class="font-semibold uppercase tracking-widest text-foreground/35 select-none {isExpanded ? 'text-[10px] px-2' : 'text-[8px]'}">
                                {isExpanded ? group.label : group.label.slice(0, 3)}
                            </span>
                        {/if}
                        <hr class="border-foreground/20 {isExpanded ? 'w-full' : 'w-6'}" />
                    </div>
                {/if}
                <div class="flex flex-col gap-2" style:align-items={isExpanded ? 'stretch' : 'center'}>
                    {#each group.items as item}
                        {@render button(item, noCaseOpen || !!item.disabled)}
                    {/each}
                </div>
            {/each}
        </div>

        <div class="flex flex-col gap-2" style:align-items={isExpanded ? 'stretch' : 'center'}>
            <ProfilePopover>
                <Button.Root
                    class="rounded-card-sm flex items-center border-1 border-border-input bg-background cursor-pointer transition-all duration-300 {isExpanded
                        ? 'justify-start gap-3 px-4 py-3 w-full'
                        : 'justify-center mx-auto size-12'}"
                >
                    <User size={24} />
                    {#if isExpanded}
                        <span class="text-sm font-medium">Profile</span>
                    {/if}
                </Button.Root>
            </ProfilePopover>
        </div>
    </div>
</div>

{#snippet button(item: SidebarItem, disabled: boolean)}
    {@const Icon = item.icon}
    {@const active = !disabled && isActive(item)}
    {#if isExpanded}
        <button
            class="align-center border-border-input rounded-card-sm bg-background-alt ring-offset-background active:scale-[0.98] active:transition:all
	focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden flex items-center gap-3 px-4 py-3 focus-visible:ring-2 focus-visible:ring-offset-2 {active
                ? 'bg-surface-hover text-foreground hover:bg-surface-active'
                : 'bg-transparent text-foreground hover:bg-surface'} {disabled ? 'opacity-35 cursor-not-allowed' : ''}"
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
                    class="align-center border-border-input rounded-card-sm bg-background-alt ring-offset-background
	focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex size-12 items-center justify-center focus-visible:ring-2 focus-visible:ring-offset-2 {active
                        ? 'bg-surface-hover text-foreground hover:bg-surface-active'
                        : 'bg-transparent text-foreground hover:bg-surface'} {disabled ? 'opacity-35 cursor-not-allowed' : 'active:scale-[0.98] active:transition:all'}"
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

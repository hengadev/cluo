<script lang="ts">
    import { Home, User, ChevronRight, ChevronLeft } from "@lucide/svelte";
    import { Button, Tooltip } from "bits-ui";
    import ProfilePopover from "$lib/custom/sidebar/ProfilePopover.svelte";
    import { items, type SidebarItem } from "$lib/constructor/sidebar";
    import { currentCase } from "$lib/stores/case";
    import { goto } from "$app/navigation";
    import { page } from "$app/stores";

    const size: number = 24;

    import { getToastContext } from "$lib/custom/global/toast/state.svelte";
    import { TOAST_LEVELS } from "$lib/custom/global/toast/type";
    const toastState = getToastContext();

    let isExpanded: boolean = $state(false);

    // Get the current route path for highlighting
    function getRouteForItem(item: SidebarItem): string {
        // Replace :id with actual case ID if available
        const caseId = $currentCase.id;
        if (!caseId) {
            // For case-specific routes without a case, return empty
            if (item.path.includes(':id')) return '';
            return item.path;
        }
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
            // No case ID selected - show toast
            toastState.add(
                TOAST_LEVELS.Warning,
                "Aucun dossier sélectionné",
                "Veuillez d'abord sélectionner un dossier"
            );
            return;
        }

        goto(routePath);
    }
</script>

<div
    class="grid-area-sidebar h-full p-1 pt-2 flex flex-col gap-10 bg-background-alt border-1 border-dark-50 relative transition-all duration-300"
    style:width={isExpanded ? '200px' : 'auto'}
    style:align-items={isExpanded ? 'stretch' : 'center'}
>
    <!-- Chevron toggle button -->
    <button
        onclick={() => (isExpanded = !isExpanded)}
        class="absolute bottom-16 text-dark-500 right-[-1.25rem] border-dark-50 bg-dark-50 p-2 rounded-3xl hover:bg-muted transition-colors"
        title={isExpanded ? "Collapse sidebar" : "Expand sidebar"}
        type="button"
    >
        {#if isExpanded}
            <ChevronLeft size={24} strokeWidth={1.5} />
        {:else}
            <ChevronRight size={24} strokeWidth={1.5} />
        {/if}
    </button>

    <div class="grid gap-4">
        <Button.Root
            class="bg-background border-1 border-border-input p-2 !mt-1 rounded-input flex items-center cursor-pointer transition-all duration-300 {isExpanded
                ? 'justify-start gap-3 px-4'
                : 'justify-center'}"
            onclick={() => goto("/")}
        >
            <Home {size} strokeWidth={1.5} />
            {#if isExpanded}
                <span class="text-sm font-medium">Home</span>
            {/if}
        </Button.Root>
    </div>
    <div class="flex flex-col justify-between h-full">
        <div
            class="flex flex-col gap-2"
            style:align-items={isExpanded ? 'stretch' : 'center'}
        >
            {#each items as item}
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

{#snippet button(item: SidebarItem)}
    {@const Icon = item.icon}
    {@const active = isActive(item)}
    {#if isExpanded}
        <button
            class="align-center border-border-input rounded-10px bg-background-alt ring-offset-background active:scale-[0.98] active:transition:all
		focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden flex items-center gap-3 px-4 py-3 focus-visible:ring-2 focus-visible:ring-offset-2 {active
                ? 'bg-foreground text-background'
                : 'bg-transparent text-foreground hover:bg-muted'}"
            onclick={() => handleItemClick(item)}
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
                    class="align-center border-border-input rounded-10px bg-background-alt ring-offset-background active:scale-[0.98] active:transition:all
		focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex size-12 items-center justify-center focus-visible:ring-2 focus-visible:ring-offset-2 {active
                        ? 'bg-foreground text-background'
                        : 'bg-transparent text-foreground hover:bg-muted'}"
                    onclick={() => handleItemClick(item)}
                >
                    <Button.Root class="cursor-pointer">
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

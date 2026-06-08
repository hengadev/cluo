<script lang="ts">
    import { navItems, utilityItems, type UtilityItem } from "$lib/constructor/header";
    import ThemeToggle from "$lib/components/ThemeToggle.svelte";
    import { Button, Tooltip } from "bits-ui";
    import { page } from "$app/stores";
</script>

<div class="header border-1 border-dark-50 animate-fade-in" style="animation-delay: 100ms;">
    <nav class="nav-items">
        {#each navItems as item}
            {@const Icon = item.icon}
            {@const isActive = $page.url.pathname === item.href || $page.url.pathname.startsWith(item.href + "/")}
            <a
                href={item.href}
                class="nav-button {isActive ? 'active' : ''}"
            >
                <Icon size={18} strokeWidth={1.75} />
                <span>{item.label}</span>
            </a>
        {/each}
    </nav>
    <div class="flex align-center gap-2">
        <div class="buttons">
            <ThemeToggle />
            {#each utilityItems as item}
                {@const DialogOrPopover = item.uiComponent}
                <DialogOrPopover>
                    {@render utilityItem(item)}
                </DialogOrPopover>
            {/each}
        </div>
    </div>
</div>

{#snippet utilityItem(item: UtilityItem)}
    {@const Icon = item.icon}
    <Tooltip.Provider>
        <Tooltip.Root delayDuration={100}>
            <Tooltip.Trigger
                class="rounded-10px p-3 bg-dark-50 ring-offset-background active:scale-[0.98] active:transition:all
		focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center focus-visible:ring-2 focus-visible:ring-offset-2
                    hover:bg-dark-100/50 hover:scale-105 transition-all duration-200 {item.bg} text-{item.fg}"
            >
                <Button.Root class="cursor-pointer">
                    <Icon size={24} strokeWidth={1.75} />
                </Button.Root>
            </Tooltip.Trigger>
            <Tooltip.Content sideOffset={8} side="bottom">
                <div
                    class="rounded-input text-[1rem] align-center bg-foreground text-background font-medium border-dark-10 shadow-popover outline-hidden z-0 flex items-center justify-center border p-2"
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
    .nav-items {
        display: flex;
        gap: 0.25rem;
        align-items: center;
    }
    .nav-button {
        display: inline-flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.5rem 1rem;
        border-radius: 10px;
        font-size: 0.875rem;
        font-weight: 500;
        color: var(--foreground-alt);
        text-decoration: none;
        transition: all 150ms ease;
    }
    .nav-button:hover {
        background: var(--dark-50);
        color: var(--foreground);
    }
    .nav-button.active {
        background: var(--dark-50);
        color: var(--foreground);
        font-weight: 600;
    }
    .buttons {
        display: flex;
        justify-content: right;
        gap: 0.5rem;
        margin-left: auto;
        text-align: right;
    }
</style>

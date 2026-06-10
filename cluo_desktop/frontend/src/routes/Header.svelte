<script lang="ts">
    import { navItems, utilityItems, type UtilityItem } from "$lib/constructor/header";
    import ThemeToggle from "$lib/components/ThemeToggle.svelte";
    import { Button, Tooltip } from "bits-ui";
    import { page } from "$app/stores";
    import { currentCase, recentCases } from "$lib/stores/case";
    import { caseStatusBadge } from "$lib/utils/badgeVariants";
    import type { CaseStatus } from "$lib/types/entities";

    const STATUS_LABELS: Record<CaseStatus, string> = {
        in_progress: "En cours",
        ready: "Prêt",
        released: "Clôturé",
    };

    const STATUS_DOT: Record<CaseStatus, string> = {
        in_progress: "bg-accent",
        ready: "bg-success",
        released: "bg-muted-foreground",
    };

    $: currentCaseInfo = $currentCase.id
        ? ($recentCases.find(c => c.id === $currentCase.id) ?? null)
        : null;
</script>

<div class="header border-1 border-border-input animate-fade-in" style="animation-delay: 100ms;">
    <nav class="nav-items">
        {#each navItems as item}
            {@const Icon = item.icon}
            {@const isActive = $page.url.pathname === item.href || $page.url.pathname.startsWith(item.href + "/")}
            <a
                href={item.href}
                class="nav-button {isActive ? 'active' : ''}"
            >
                <Icon size={16} strokeWidth={1.75} />
                <span>{item.label}</span>
            </a>
        {/each}
    </nav>

    <div class="case-banner">
        {#if currentCaseInfo}
            <div class="case-chip">
                <span class="status-dot {STATUS_DOT[currentCaseInfo.status]}"></span>
                <span class="case-title">{currentCaseInfo.title}</span>
                <span class="status-label {caseStatusBadge(currentCaseInfo.status)}">
                    {STATUS_LABELS[currentCaseInfo.status]}
                </span>
            </div>
        {/if}
    </div>

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

{#snippet utilityItem(item: UtilityItem)}
    {@const Icon = item.icon}
    <Tooltip.Provider>
        <Tooltip.Root delayDuration={100}>
            <Tooltip.Trigger
                class="rounded-card-sm p-3 bg-surface ring-offset-background active:scale-[0.98] active:transition:all
		focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center focus-visible:ring-2 focus-visible:ring-offset-2
                    hover:bg-surface-hover hover:scale-105 transition-all duration-200 {item.bg} text-{item.fg}"
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
        display: grid;
        grid-template-columns: 1fr auto 1fr;
        align-items: center;
        padding: 0.5rem 2rem;
        gap: 1rem;
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
        border-radius: 0;
        border-bottom: 2px solid transparent;
        font-size: 0.875rem;
        font-weight: 400;
        color: var(--foreground-alt);
        text-decoration: none;
        background: transparent;
        transition: all 150ms ease;
    }
    .nav-button:hover {
        background: transparent;
        color: var(--foreground);
    }
    .nav-button.active {
        border-bottom-color: var(--accent);
        background: transparent;
        color: var(--foreground);
        font-weight: 500;
    }
    .case-banner {
        display: flex;
        justify-content: center;
        align-items: center;
    }
    .case-chip {
        display: inline-flex;
        align-items: center;
        gap: 0.625rem;
        padding: 0.375rem 0.875rem;
        border-radius: 999px;
        background: var(--surface);
        border: 1px solid var(--border-input);
        max-width: 380px;
    }
    .status-dot {
        flex-shrink: 0;
        width: 7px;
        height: 7px;
        border-radius: 50%;
    }
    .case-title {
        font-size: 0.8125rem;
        font-weight: 600;
        color: var(--foreground);
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        max-width: 240px;
    }
    .status-label {
        flex-shrink: 0;
        font-size: 0.6875rem;
        font-weight: 500;
        padding: 0.125rem 0.5rem;
        border-radius: 999px;
        white-space: nowrap;
    }
    .buttons {
        display: flex;
        justify-content: flex-end;
        gap: 0.5rem;
    }
</style>

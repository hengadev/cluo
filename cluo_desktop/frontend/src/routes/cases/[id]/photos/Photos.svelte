<script lang="ts">
    import type { Component } from "svelte";
    import { Tabs, Button, Select } from "bits-ui";
    import {
        Check,
        Upload,
        Grid,
        LayoutList,
        ChevronDown,
        ClockArrowDown,
        ClockArrowUp,
        CircleCheckBig,
        Circle,
    } from "@lucide/svelte";
    import AllPhotos from "./allPhotos.svelte";

    import {
        type LayoutState,
        type SortState,
        LAYOUT_STATES,
        SORT_STATES,
    } from "./types";

    import { isMockEnabled } from "$lib/config";
    import { images as mockImages } from "./mockData";
    import { fetchCaseImages } from "$lib/services/api";
    import { onMount } from "svelte";
    import { page } from "$app/stores";

    // NOTE: That thing should be props to be honest
    let images = $state(mockImages);
    let loading = $state(false);

    // Load images based on mock flag
    onMount(async () => {
        if (!isMockEnabled()) {
            loading = true;
            try {
                const caseId = $page.params.id;
                const apiImages = await fetchCaseImages(caseId);
                if (apiImages.length > 0) {
                    images = apiImages as typeof images;
                }
            } catch (error) {
                console.error("Failed to fetch images:", error);
                images = [];
            } finally {
                loading = false;
            }
        }
    });

    interface Layout {
        value: LayoutState;
        label: string;
        icon: Component;
    }

    const layouts: Layout[] = [
        { value: LAYOUT_STATES.Grid, label: "Grille", icon: Grid },
        { value: LAYOUT_STATES.List, label: "Liste", icon: LayoutList },
    ];

    let layoutValue = $state<LayoutState>(LAYOUT_STATES.Grid);
    const layoutLabel = $derived(
        layoutValue
            ? layouts.find((theme) => theme.value === layoutValue)?.label
            : "Grille",
    );
    const LayoutIcon: Component | undefined = $derived(
        layouts.find((theme) => theme.value === layoutValue)?.icon,
    );

    interface Sort {
        value: SortState;
        label: string;
        icon: Component;
    }

    const sorts: Sort[] = [
        {
            value: SORT_STATES.NewestFirst,
            label: "Plus récents",
            icon: ClockArrowDown,
        },
        {
            value: SORT_STATES.OldestFirst,
            label: "Plus anciens",
            icon: ClockArrowUp,
        },
        {
            value: SORT_STATES.SelectedFirst,
            label: "Sélectionnés en premier",
            icon: CircleCheckBig,
        },
        {
            value: SORT_STATES.NonSelectedFirst,
            label: "Non sélectionnés en premier",
            icon: Circle,
        },
    ];

    let sortValue = $state<string>(SORT_STATES.NewestFirst);
    const sortLabel = $derived(
        sortValue
            ? sorts.find((theme) => theme.value === sortValue)?.label
            : "Grille",
    );
    const SortIcon: Component | undefined = $derived(
        sorts.find((theme) => theme.value === sortValue)?.icon,
    );
</script>

<div class="content p-6">
    <Tabs.Root value="all" class="rounded-card w-full p-3">
        <div class="flex justify-between items-center">
            <Tabs.List
                class="rounded-9px bg-dark-10 shadow-mini-inset dark:bg-background grid w-full grid-cols-2 gap-1 p-1 text-sm font-semibold leading-[0.01em] dark:border dark:border-neutral-600/30 max-w-[400px]"
            >
                <Tabs.Trigger
                    value="all"
                    class="data-[state=active]:shadow-mini dark:data-[state=active]:bg-muted h-8 rounded-[7px] bg-transparent py-2 data-[state=active]:bg-white"
                    >Toutes les photos</Tabs.Trigger
                >
                <Tabs.Trigger
                    value="selection"
                    class="data-[state=active]:shadow-mini dark:data-[state=active]:bg-muted h-8 rounded-[7px] bg-transparent py-2 data-[state=active]:bg-white"
                    >Photos selectionnees</Tabs.Trigger
                >
            </Tabs.List>
            <Button.Root
                class="gap-2 items-center h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex px-4 text-[15px] font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98]"
            >
                <Upload />
                <p>Ajouter une photo</p>
            </Button.Root>
            <div>
                {@render select_bits_ui(
                    layoutLabel,
                    LayoutIcon,
                    layouts,
                    (value) => (layoutValue = value),
                )}
                {@render select_bits_ui(
                    sortLabel,
                    SortIcon,
                    sorts,
                    (value) => (sortValue = value),
                )}
            </div>
        </div>
        <Tabs.Content value="all" class="select-none pt-3">
            {#if loading}
                <p class="text-muted-foreground">Chargement des photos...</p>
            {:else if images.length === 0}
                <p class="text-muted-foreground">
                    Aucune photo disponible. {isMockEnabled() ? '' : '(API non configurée)'}
                </p>
            {:else}
                <AllPhotos {layoutValue} {images} />
            {/if}
        </Tabs.Content>
        <Tabs.Content value="selection" class="select-none pt-3">
            <div>here is the content with the selection</div>
        </Tabs.Content>
    </Tabs.Root>
</div>

{#snippet select_bits_ui(
    labelState: string | undefined,
    IconState: Component | undefined,
    items: any[],
    onSelect: (v: LayoutState) => void,
)}
    <Select.Root
        type="single"
        onValueChange={(v) => onSelect(v as LayoutState)}
        {items}
        allowDeselect={true}
    >
        <Select.Trigger
            class="rounded-10px p-3 ring-offset-background active:scale-[0.98] active:transition:all focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center gap-2 focus-visible:ring-2 focus-visible:ring-offset-2 hover:bg-dark-100/50 text-dark-900 bg-dark-50"
            title="Grille"
            aria-label="Choisis une disposition"
        >
            <IconState />
            <p class="capitalize">{labelState}</p>
            <ChevronDown size={16} class="text-dark-600" strokeWidth={1.75} />
        </Select.Trigger>
        <Select.Portal>
            <Select.Content
                class="focus-override border-muted bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 outline-hidden z-50 h-96 max-h-[var(--bits-select-content-available-height)] w-[var(--bits-select-anchor-width)] min-w-[var(--bits-select-anchor-width)] select-none rounded-xl border px-1 py-3 data-[side=bottom]:translate-y-1 data-[side=left]:-translate-x-1 data-[side=right]:translate-x-1 data-[side=top]:-translate-y-1"
                sideOffset={10}
            >
                <Select.Viewport class="p-1">
                    {#each items as item, i (i + item.value)}
                        <Select.Item
                            class="rounded-button data-highlighted:bg-muted outline-hidden data-disabled:opacity-50 flex gap-2 h-10 w-full select-none items-center py-3 pl-5 pr-1.5 text-sm capitalize"
                            value={item.value}
                            label={item.label}
                            disabled={item.disabled}
                        >
                            {#snippet children({ selected })}
                                {item.label}
                                {#if selected}
                                    <div class="ml-auto">
                                        <Check aria-label="check" />
                                    </div>
                                {/if}
                            {/snippet}
                        </Select.Item>
                    {/each}
                </Select.Viewport>
            </Select.Content>
        </Select.Portal>
    </Select.Root>
{/snippet}

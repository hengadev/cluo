<script lang="ts">
    import { Combobox } from "bits-ui";
    import {
        ChevronDown,
        ChevronsUp,
        ChevronsDown,
        Check,
        X,
    } from "@lucide/svelte";
    import { cn } from "$lib/utils";
    import { createEventDispatcher } from "svelte";

    interface Props {
        options: Array<{ value: string; label: string; disabled?: boolean }>;
        id?: string;
        name?: string;
        placeholder?: string;
        value?: string;
        disabled?: boolean;
        required?: boolean;
        class?: string;
        size?: "sm" | "md" | "lg";
    }

    let {
        options,
        id,
        name,
        placeholder = "Select an option",
        value = $bindable(""),
        disabled = false,
        required = false,
        class: className = "",
        size = "md",
        ...restProps
    }: Props = $props();

    const dispatch = createEventDispatcher();

    let searchValue = $state("");
    let isOpen = $state(false);

    const selectedLabel = $derived(
        value
            ? options.find((opt) => opt.value === value)?.label
            : placeholder,
    );

    const filteredOptions = $derived(
        searchValue === ""
            ? options
            : options.filter((opt) =>
                  opt.label.toLowerCase().includes(searchValue.toLowerCase()),
              ),
    );

    const baseClasses =
        "flex w-full bg-dark-50 placeholder-dark-300 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-dark focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 rounded-input";

    const sizeClasses = {
        sm: "h-9 px-3 text-sm",
        md: "h-10 px-3 py-2 text-sm",
        lg: "h-12 px-3 py-2 text-base",
    };

    const inputClass = cn(baseClasses, sizeClasses[size], className);

    function handleValueChange(v: string | undefined) {
        if (v !== undefined) {
            value = v;
            dispatch("change", { value: v });
        }
    }

    function handleOpenChange(o: boolean) {
        isOpen = o;
        if (!o) {
            searchValue = "";
        }
    }

    function clearValue(e: MouseEvent) {
        e.stopPropagation();
        value = "";
        dispatch("change", { value: "" });
    }
</script>

<Combobox.Root
    type="single"
    {name}
    {disabled}
    {required}
    onValueChange={handleValueChange}
    onOpenChange={handleOpenChange}
>
    <div class="relative">
        <Combobox.Input
            {id}
            oninput={(e) => (searchValue = e.currentTarget.value)}
            class={inputClass}
            placeholder={isOpen ? "Search..." : selectedLabel}
            aria-label={placeholder}
            {...restProps}
        />
        {#if value && !disabled}
            <button
                type="button"
                onclick={clearValue}
                class="absolute right-10 top-1/2 -translate-y-1/2 rounded p-1 hover:bg-dark-10"
                aria-label="Clear selection"
            >
                <X class="h-4 w-4" />
            </button>
        {/if}
        <Combobox.Trigger
            class="absolute end-3 top-1/2 size-6 -translate-y-1/2"
        >
            <ChevronDown class="text-muted-foreground size-6" />
        </Combobox.Trigger>
    </div>
    <Combobox.Portal>
        <Combobox.Content
            class="focus-override border-muted bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 outline-hidden z-50 h-96 max-h-[var(--bits-combobox-content-available-height)] w-[var(--bits-combobox-anchor-width)] min-w-[var(--bits-combobox-anchor-width)] select-none rounded-xl border px-1 py-3 data-[side=bottom]:translate-y-1 data-[side=left]:-translate-x-1 data-[side=right]:translate-x-1 data-[side=top]:-translate-y-1"
            sideOffset={10}
        >
            <Combobox.ScrollUpButton
                class="flex w-full items-center justify-center py-1"
            >
                <ChevronsUp class="size-3" />
            </Combobox.ScrollUpButton>
            <Combobox.Viewport class="p-1">
                {#each filteredOptions as option, i (i + option.value)}
                    <Combobox.Item
                        class="rounded-button data-highlighted:bg-muted outline-hidden flex h-10 w-full select-none items-center py-3 pl-5 pr-1.5 text-sm capitalize"
                        value={option.value}
                        label={option.label}
                        disabled={option.disabled}
                    >
                        {#snippet children({ selected })}
                            {option.label}
                            {#if selected}
                                <div class="ml-auto">
                                    <Check />
                                </div>
                            {/if}
                        {/snippet}
                    </Combobox.Item>
                {:else}
                    <span class="block px-5 py-2 text-sm text-muted-foreground">
                        No results found, try again.
                    </span>
                {/each}
            </Combobox.Viewport>
            <Combobox.ScrollDownButton
                class="flex w-full items-center justify-center py-1"
            >
                <ChevronsDown class="size-3" />
            </Combobox.ScrollDownButton>
        </Combobox.Content>
    </Combobox.Portal>
</Combobox.Root>

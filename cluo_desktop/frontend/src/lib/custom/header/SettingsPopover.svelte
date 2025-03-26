<script lang="ts">
    import { Popover, Separator, Select, Label, Switch } from "bits-ui";
    import {
        Languages,
        Check,
        ChevronsDown,
        ChevronsUp,
        ChevronsUpDown,
    } from "@lucide/svelte";

    const languages = [
        { value: "english", label: "English", disabled: false },
        { value: "french", label: "French", disabled: false },
    ];

    let value = $state<string>("");
    const selectedLabel = $derived(
        value
            ? languages.find((language) => language.value === value)?.label
            : "Select a language",
    );

    type Props = { children: import("svelte").Snippet };
    let { children }: Props = $props();
</script>

<Popover.Root>
    <Popover.Trigger>
        {@render children()}
    </Popover.Trigger>
    <Popover.Portal>
        <Popover.Content
            class="border-dark-10 bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 z-30 w-full max-w-[360px] rounded-[12px] border p-4"
            sideOffset={8}
            align="end"
        >
            <div class="flex items-center">
                <div class="flex flex-col">
                    <h4
                        class="text-[17px] font-semibold leading-5 tracking-[-0.01em]"
                    >
                        Parametres
                    </h4>
                    <p class="text-muted-foreground text-sm font-medium">
                        Customise ton experience utilisateur
                    </p>
                </div>
            </div>
            <Separator.Root
                class="bg-dark-10 -mx-4 !mb-6 !mt-[17px] block h-px"
            />
            <div class="flex flex-col items-center pb-2 gap-2">
                {@render switchComponent("Mode sombre")}
                {@render switchComponent("Notifications")}
                {@render switchComponent("Compact view")}
                <Separator.Root
                    class="w-full bg-dark-10 -mx-4 !mb-6 !mt-[17px] block h-px"
                />
                {@render selector()}
                {@render selector()}
            </div>
        </Popover.Content>
    </Popover.Portal>
</Popover.Root>

{#snippet switchComponent(label: string)}
    <div class="flex items-center justify-between space-x-3 w-full">
        <Label.Root for="dnd" class="text-sm font-medium text-dark"
            >{label}</Label.Root
        >
        <Switch.Root
            id="dnd"
            name="hello"
            class="focus-visible:ring-foreground focus-visible:ring-offset-background data-[state=checked]:bg-foreground data-[state=unchecked]:bg-dark-10 data-[state=unchecked]:shadow-mini-inset dark:data-[state=checked]:bg-foreground focus-visible:outline-hidden peer inline-flex h-[24px] min-h-[24px] w-[48px] shrink-0 cursor-pointer items-center rounded-full px-[3px] transition-colors focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
        >
            <Switch.Thumb
                class="bg-background data-[state=unchecked]:shadow-mini dark:border-background/30 dark:bg-foreground dark:shadow-popover pointer-events-none block size-[18px] shrink-0 rounded-full transition-transform data-[state=checked]:translate-x-6 data-[state=unchecked]:translate-x-0 dark:border dark:data-[state=unchecked]:border"
            />
        </Switch.Root>
    </div>
{/snippet}

{#snippet selector()}
    <Select.Root
        type="single"
        onValueChange={(v) => (value = v)}
        items={languages}
    >
        <Select.Trigger
            class="h-input rounded-9px border-border-input bg-background data-placeholder:text-foreground-alt/50 inline-flex justify-between w-[296px] select-none items-center border px-[11px] text-sm transition-colors"
            aria-label="Select a language"
        >
            <div class="flex gap-2">
                <Languages class="text-muted-foreground mr-[9px] size-6" />
                {selectedLabel}
            </div>
            <ChevronsUpDown class="text-muted-foreground ml-auto size-6" />
        </Select.Trigger>
        <Select.Portal>
            <Select.Content
                class="focus-override border-muted bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 outline-hidden z-50 h-96 max-h-[var(--bits-select-content-available-height)] w-[var(--bits-select-anchor-width)] min-w-[var(--bits-select-anchor-width)] select-none rounded-xl border px-1 py-3 data-[side=bottom]:translate-y-1 data-[side=left]:-translate-x-1 data-[side=right]:translate-x-1 data-[side=top]:-translate-y-1"
                sideOffset={10}
            >
                <Select.ScrollUpButton
                    class="flex w-full items-center justify-center"
                >
                    <ChevronsUp class="size-3" />
                </Select.ScrollUpButton>
                <Select.Viewport class="p-1">
                    {#each languages as language, i (i + language.value)}
                        <Select.Item
                            class="rounded-button data-highlighted:bg-muted outline-hidden data-disabled:opacity-50 flex align-center justify-between h-10 w-full select-none items-center py-3 pl-5 pr-1.5 text-sm  capitalize"
                            value={language.value}
                            label={language.label}
                            disabled={language.disabled}
                        >
                            {#snippet children({ selected })}
                                {language.label}
                                {#if selected}
                                    <div class="ml-auto">
                                        <Check />
                                    </div>
                                {/if}
                            {/snippet}
                        </Select.Item>
                    {/each}
                </Select.Viewport>
                <Select.ScrollDownButton
                    class="flex w-full items-center justify-center"
                >
                    <ChevronsDown class="size-3" />
                </Select.ScrollDownButton>
            </Select.Content>
        </Select.Portal>
    </Select.Root>
{/snippet}

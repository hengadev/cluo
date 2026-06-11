<!-- Hallmark · component: popover/settings · genre: modern-minimal · theme: project-system
     states: default · hover · focus · active
     contrast: pass (foreground/muted-foreground on background) -->
<script lang="ts">
    import { Popover, Select, Label, Switch } from "bits-ui";
    import {
        Languages,
        Check,
        ChevronsDown,
        ChevronsUp,
        ChevronsUpDown,
        RefreshCw,
    } from "@lucide/svelte";
    import { updateDialogOpen } from '$lib/stores/update';

    const languages = [
        { value: "french", label: "Français", disabled: false },
        { value: "english", label: "English", disabled: false },
    ];

    let language = $state<string>("");
    const selectedLabel = $derived(
        language
            ? languages.find(l => l.value === language)?.label
            : "Sélectionner une langue"
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
            class="border-border-card bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 z-30 w-[320px] rounded-card-sm border p-5"
            sideOffset={8}
            align="end"
        >
            <div class="mb-5">
                <h4 class="text-sm font-semibold tracking-tight text-foreground">Paramètres</h4>
                <p class="text-xs text-muted-foreground mt-0.5">Personnalisez votre expérience</p>
            </div>

            <div class="flex flex-col gap-5">
                <div class="flex flex-col gap-3">
                    {@render sectionLabel("Préférences")}
                    {@render switchRow("Mode sombre", "switch-dark-mode")}
                    {@render switchRow("Notifications", "switch-notifications")}
                    {@render switchRow("Vue compacte", "switch-compact")}
                </div>

                <div class="flex flex-col gap-3">
                    {@render sectionLabel("Langue")}
                    {@render selector()}
                </div>

                <div class="flex flex-col gap-3">
                    {@render sectionLabel("Application")}
                    <button
                        onclick={() => updateDialogOpen.set(true)}
                        class="flex items-center gap-2.5 w-full text-sm font-medium text-muted-foreground hover:text-foreground transition-interactive duration-150 cursor-pointer py-0.5"
                    >
                        <RefreshCw class="size-4 shrink-0" />
                        Rechercher des mises à jour
                    </button>
                </div>
            </div>
        </Popover.Content>
    </Popover.Portal>
</Popover.Root>

{#snippet sectionLabel(label: string)}
    <div class="flex items-center gap-3">
        <span class="text-xs font-medium text-muted-foreground uppercase tracking-wider whitespace-nowrap">{label}</span>
        <div class="h-px flex-1 bg-border-input"></div>
    </div>
{/snippet}

{#snippet switchRow(label: string, id: string)}
    <div class="flex items-center justify-between">
        <Label.Root for={id} class="text-sm font-medium text-foreground cursor-pointer">{label}</Label.Root>
        <Switch.Root
            {id}
            name={id}
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
        onValueChange={(v) => (language = v)}
        items={languages}
    >
        <Select.Trigger
            class="h-input rounded-input border-border-input bg-background data-placeholder:text-muted-foreground/50 inline-flex justify-between w-full select-none items-center border px-3 text-sm transition-interactive duration-150 hover:border-border-input-hover"
            aria-label="Sélectionner une langue"
        >
            <div class="flex items-center gap-2">
                <Languages class="text-muted-foreground size-4 shrink-0" />
                {selectedLabel}
            </div>
            <ChevronsUpDown class="text-muted-foreground size-4 shrink-0" />
        </Select.Trigger>
        <Select.Portal>
            <Select.Content
                class="focus-override border-border-card bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 outline-hidden z-50 max-h-[var(--bits-select-content-available-height)] w-[var(--bits-select-anchor-width)] min-w-[var(--bits-select-anchor-width)] select-none rounded-card-sm border px-1 py-2"
                sideOffset={6}
            >
                <Select.ScrollUpButton class="flex w-full items-center justify-center py-1">
                    <ChevronsUp class="size-3" />
                </Select.ScrollUpButton>
                <Select.Viewport class="p-1">
                    {#each languages as lang, i (i + lang.value)}
                        <Select.Item
                            class="rounded-button data-highlighted:bg-muted outline-hidden data-disabled:opacity-50 flex items-center justify-between h-9 w-full select-none py-2 pl-4 pr-2 text-sm"
                            value={lang.value}
                            label={lang.label}
                            disabled={lang.disabled}
                        >
                            {#snippet children({ selected })}
                                {lang.label}
                                {#if selected}
                                    <Check class="size-4 text-foreground" />
                                {/if}
                            {/snippet}
                        </Select.Item>
                    {/each}
                </Select.Viewport>
                <Select.ScrollDownButton class="flex w-full items-center justify-center py-1">
                    <ChevronsDown class="size-3" />
                </Select.ScrollDownButton>
            </Select.Content>
        </Select.Portal>
    </Select.Root>
{/snippet}

<script lang="ts">
    import { Home, User } from "@lucide/svelte";
    import { Button, Tooltip } from "bits-ui";
    import ProfilePopover from "$lib/custom/sidebar/ProfilePopover.svelte";
    import { items, type SidebarItem } from "$lib/constructor/sidebar";
    let selected: string = $state(items[0].title);
    const size: number = 32;

    import { getToastContext } from "$lib/custom/global/toast/state.svelte";
    import { TOAST_LEVELS } from "$lib/custom/global/toast/type";
    const toastState = getToastContext();
    const mockToast = {
        level: TOAST_LEVELS.Info,
        title: "Un titre pour le toast",
        message:
            "Le message de ce toast n'est la qu'a titre indicatif en realite.",
    };
</script>

<div
    class="grid-area-sidebar h-full p-1 pt-2 flex flex-col gap-10 items-center bg-[#fafafa] border-2 border-[#e5e7eb]"
>
    <div class="grid gap-4">
        <Button.Root
            class="bg-white border-dark border-2 p-2 !mt-1 rounded-input flex items-center justify-center cursor-pointer"
            onclick={() =>
                toastState.add(
                    mockToast.level,
                    mockToast.title,
                    mockToast.message,
                )}
        >
            <Home {size} strokeWidth={1.5} />
        </Button.Root>
    </div>
    <div class="flex flex-col justify-between h-full">
        <div class="flex flex-col items-center gap-4">
            {#each items as item}
                {@render button(item)}
            {/each}
        </div>
        <ProfilePopover>
            <Button.Root
                class="rounded-10px flex items-center justify-center border-1 border-[#e5e7eb] mx-auto size-12 bg-white cursor-pointer"
            >
                <User {size} />
            </Button.Root>
        </ProfilePopover>
    </div>
</div>

{#snippet button(item: SidebarItem)}
    <Tooltip.Provider>
        <Tooltip.Root delayDuration={100}>
            <Tooltip.Trigger
                class="align-center border-border-input rounded-10px bg-background-alt ring-offset-background active:scale-[0.98] active:transition:all 
		focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex size-12 items-center justify-center focus-visible:ring-2 focus-visible:ring-offset-2 {item.title ===
                selected
                    ? 'bg-dark text-white'
                    : 'bg-transparent text-black hover:bg-[#f0f0f0]'}"
                onclick={() => {
                    item.fn;
                    selected = item.title;
                }}
            >
                {@const Icon = item.icon}
                <Button.Root class="cursor-pointer">
                    <Icon size={32} strokeWidth={1.5} />
                </Button.Root>
            </Tooltip.Trigger>
            <Tooltip.Content sideOffset={8} side="right">
                <div
                    class="rounded-input text-[1rem] align-center bg-dark text-white font-medium gap-3 border-dark-10 shadow-popover outline-hidden z-0 flex items-center justify-center border p-2"
                >
                    {item.title}
                </div>
            </Tooltip.Content>
        </Tooltip.Root>
    </Tooltip.Provider>
{/snippet}

<style>
    .grid-area-sidebar {
        grid-area: sidebar;
    }
</style>

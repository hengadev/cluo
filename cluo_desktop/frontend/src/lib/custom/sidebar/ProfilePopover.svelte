<script lang="ts">
    import { Button, Popover, Separator } from "bits-ui";
    import { User, LogOut, BadgeCheck, Bell, CreditCard, RefreshCw } from "@lucide/svelte";
    import ConfirmDialog from "$lib/custom/global/ConfirmDialog.svelte";
    import UpdateDialog from "$lib/custom/global/UpdateDialog.svelte";

    type Props = { children: import("svelte").Snippet };
    let { children }: Props = $props();

    type ButtonItem = {
        icon: typeof import("@lucide/svelte").Icon;
        title: string;
        onclick?: () => void;
    };
    let buttons: ButtonItem[] = [
        { icon: BadgeCheck, title: "Account" },
        { icon: CreditCard, title: "Billing" },
        { icon: Bell, title: "Notifications" },
    ];

    let updateDialogOpen = $state(false);
</script>

<Popover.Root>
    <Popover.Trigger>
        {@render children()}
    </Popover.Trigger>
    <Popover.Portal>
        <Popover.Content
            class=" border-dark-10 bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 z-30 w-full max-w-[480px] rounded-[12px] border p-1"
            sideOffset={12}
            side="right"
            align="end"
        >
            <div class="flex items-center gap-4 p-2">
                <div
                    class="rounded-10px flex items-center justify-center border-1 border-border-input mx-auto size-8 bg-background cursor-pointer"
                >
                    <User size={24} />
                </div>
                <div>
                    <p class="font-semibold text-base">John</p>
                    <p class="text-sm">johndoe@example.com</p>
                </div>
            </div>
            {@render separator()}
            <div class="max-h-[400px] h-full flex flex-col items-center">
                {#each buttons as btn}
                    {@render button(btn)}
                {/each}
            </div>
            {@render separator()}
            <Button.Root
                class="p-2 w-full rounded-input hover:bg-muted cursor-pointer"
                onclick={() => (updateDialogOpen = true)}
            >
                <div class="text-base flex gap-2">
                    <RefreshCw size={16} />
                    <p>Check for Updates</p>
                </div>
            </Button.Root>
            {@render separator()}
            <ConfirmDialog>
                {@render button({ icon: LogOut, title: "Se deconnecter" })}
            </ConfirmDialog>
        </Popover.Content>
    </Popover.Portal>
</Popover.Root>

<UpdateDialog bind:open={updateDialogOpen} />

{#snippet separator()}
    <Separator.Root class="bg-dark-10 !my-1 -mx-4 block h-px" />
{/snippet}
{#snippet button(btn: ButtonItem)}
    {@const Icon = btn.icon}
    <Button.Root class="p-2 w-full rounded-input hover:bg-muted cursor-pointer">
        <div class="text-base flex gap-2">
            <Icon size={16} />
            <p>{btn.title}</p>
        </div>
    </Button.Root>
{/snippet}

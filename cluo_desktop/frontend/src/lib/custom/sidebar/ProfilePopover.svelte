<script lang="ts">
    import { Button, Popover, Separator } from "bits-ui";
    import { User, LogOut, BadgeCheck, Bell, CreditCard, RefreshCw } from "@lucide/svelte";
    import ConfirmDialog from "$lib/custom/global/ConfirmDialog.svelte";
    import UpdateDialog from "$lib/custom/global/UpdateDialog.svelte";
    import { auth } from "$lib/stores/auth";

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

    async function handleLogout() {
        await auth.logout();
    }
</script>

<Popover.Root>
    <Popover.Trigger>
        {@render children()}
    </Popover.Trigger>
    <Popover.Portal>
        <Popover.Content
            class="border-border-input bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 z-30 w-full max-w-[280px] rounded-card-sm border overflow-hidden"
            sideOffset={12}
            side="right"
            align="end"
        >
            <div class="flex items-center gap-3 px-4 py-3">
                <div class="flex items-center justify-center size-8 rounded-input border border-border-input bg-muted shrink-0">
                    <User size={15} class="text-foreground-alt" />
                </div>
                <div class="min-w-0">
                    <p class="text-sm font-semibold truncate">{$auth.user?.email ?? 'Admin'}</p>
                    <p class="text-xs text-foreground-alt capitalize">{$auth.user?.role ?? 'admin'}</p>
                </div>
            </div>
            {@render separator()}
            <div class="py-1">
                {#each buttons as btn}
                    {@render button(btn)}
                {/each}
            </div>
            {@render separator()}
            <div class="py-1">
                <Button.Root
                    class="flex items-center gap-2.5 w-full px-4 py-2 text-sm text-foreground-alt hover:text-foreground hover:bg-muted cursor-pointer transition-colors"
                    onclick={() => (updateDialogOpen = true)}
                >
                    <RefreshCw size={14} />
                    Vérifier les mises à jour
                </Button.Root>
            </div>
            {@render separator()}
            <div class="py-1">
                <ConfirmDialog
                    onConfirm={handleLogout}
                    title="Se déconnecter"
                    description="Êtes-vous sûr de vouloir vous déconnecter ?"
                >
                    {@render button({ icon: LogOut, title: "Se déconnecter" })}
                </ConfirmDialog>
            </div>
        </Popover.Content>
    </Popover.Portal>
</Popover.Root>

<UpdateDialog bind:open={updateDialogOpen} />

{#snippet separator()}
    <Separator.Root class="bg-border-input block h-px" />
{/snippet}
{#snippet button(btn: ButtonItem)}
    {@const Icon = btn.icon}
    <Button.Root class="flex items-center gap-2.5 w-full px-4 py-2 text-sm text-foreground-alt hover:text-foreground hover:bg-muted cursor-pointer transition-colors">
        <Icon size={14} />
        {btn.title}
    </Button.Root>
{/snippet}
